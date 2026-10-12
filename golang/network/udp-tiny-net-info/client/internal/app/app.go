package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/client/internal/config"
	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/client/pkg/assembly"
	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/client/pkg/protocol"
)

type App struct {
	cfg       *config.Config
	log       *slog.Logger
	conn      *net.UDPConn
	serverUDP *net.UDPAddr
	assembler *assembly.Assembler
	reqSeq    atomic.Uint64
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	serverAddr, err := net.ResolveUDPAddr("udp", cfg.ServerAddr)
	if err != nil {
		return nil, fmt.Errorf("resolve server addr: %w", err)
	}

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, fmt.Errorf("listen client udp: %w", err)
	}

	return &App{
		cfg:       cfg,
		log:       log,
		conn:      conn,
		serverUDP: serverAddr,
		assembler: assembly.NewAssembler(cfg.AssemblerTTL),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.conn.Close()
	defer a.assembler.Stop()

	a.log.Info("client started",
		slog.String("local_addr", a.conn.LocalAddr().String()),
		slog.String("target_server", a.serverUDP.String()),
	)

	fmt.Println("Available commands: GET <protocol>, LIST, COUNT, SEARCH <query>, EXIT")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fmt.Print("\n> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if strings.EqualFold(input, "EXIT") {
			return nil
		}

		tokens := strings.Fields(input)
		req := &protocol.Request{
			Command: tokens[0],
			Args:    tokens[1:],
		}

		resp, err := a.SendWithRetry(ctx, req)
		if err != nil {
			a.log.Error("request failed", slog.Any("error", err))
			continue
		}

		fmt.Printf("[%s] %s\n", resp.Status, resp.Body)
	}

	return scanner.Err()
}

func (a *App) SendWithRetry(ctx context.Context, req *protocol.Request) (*protocol.Response, error) {
	reqID := fmt.Sprintf("%d", a.reqSeq.Add(1))
	reqPayload := req.Marshal()
	chunks := protocol.SplitIntoChunks(reqID, reqPayload, a.cfg.MaxChunkSize)

	for attempt := 1; attempt <= a.cfg.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		a.log.Debug("sending request",
			slog.String("req_id", reqID),
			slog.Int("attempt", attempt),
			slog.String("cmd", req.Command),
			slog.Int("chunks", len(chunks)),
		)

		for _, chunk := range chunks {
			if _, err := a.conn.WriteToUDP(chunk.Marshal(), a.serverUDP); err != nil {
				return nil, fmt.Errorf("write udp: %w", err)
			}
		}

		resp, err := a.awaitResponse(reqID)
		if err == nil {
			return resp, nil
		}

		if errors.Is(err, os.ErrDeadlineExceeded) {
			a.log.Warn("request timed out, retrying...",
				slog.String("req_id", reqID),
				slog.Int("attempt", attempt),
			)
			continue
		}

		return nil, err
	}

	return nil, fmt.Errorf("server unreachable after %d attempts", a.cfg.MaxRetries)
}

func (a *App) awaitResponse(targetReqID string) (*protocol.Response, error) {
	deadline := time.Now().Add(a.cfg.Timeout)
	if err := a.conn.SetReadDeadline(deadline); err != nil {
		return nil, fmt.Errorf("set deadline: %w", err)
	}

	buf := make([]byte, 2048)

	for {
		n, from, err := a.conn.ReadFromUDP(buf)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				return nil, os.ErrDeadlineExceeded
			}
			return nil, err
		}

		if !from.IP.Equal(a.serverUDP.IP) || from.Port != a.serverUDP.Port {
			continue
		}

		rawPacket := make([]byte, n)
		copy(rawPacket, buf[:n])

		chunk, err := protocol.UnmarshalChunk(rawPacket)
		if err != nil {
			a.log.Warn("malformed server packet", slog.Any("error", err))
			continue
		}

		if chunk.MessageID != targetReqID {
			continue
		}

		var fullPayload []byte
		if chunk.Total > 1 {
			a.log.Debug("received chunk",
				slog.String("req_id", chunk.MessageID),
				slog.Int("part", chunk.Part),
				slog.Int("total", chunk.Total),
				slog.Int("size", len(chunk.Data)),
			)
			sessionKey := fmt.Sprintf("%s#%s", from.String(), chunk.MessageID)
			payload, complete := a.assembler.AddChunk(sessionKey, chunk.Part, chunk.Total, chunk.Data)
			if !complete {
				continue
			}
			a.log.Debug("all chunks assembled",
				slog.String("req_id", chunk.MessageID),
				slog.Int("total_chunks", chunk.Total),
				slog.Int("total_bytes", len(payload)),
			)
			fullPayload = payload
		} else {
			a.log.Debug("received single chunk response",
				slog.String("req_id", chunk.MessageID),
				slog.Int("size", len(chunk.Data)),
			)
			fullPayload = chunk.Data
		}

		return protocol.UnmarshalResponse(fullPayload)
	}
}
