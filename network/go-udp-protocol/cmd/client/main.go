package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/sudo-odner/tmp/network/go-udp-protocol/pkg/protocol"
	"github.com/sudo-odner/tmp/network/go-udp-protocol/pkg/protocol/payload"
)

type Config struct {
	serverAddr     string
	timeout        time.Duration
	assembleTTL    time.Duration
	maxRetry       int
	maxPayloadSize int
	reqSeq         atomic.Uint64
	logDebug       bool
}

func main() {
	// Config
	var cfg Config
	flag.StringVar(&cfg.serverAddr, "server-addr", "127.0.0.1:5008", "Адресс сервера (default: '127.0.0.1:5008')")
	flag.DurationVar(&cfg.timeout, "timeout", 2*time.Second, "Время ожидания (default: 2s)")
	flag.DurationVar(&cfg.assembleTTL, "assemble-ttl", 6*time.Second, "Время жизни пакетов на сборку пакетов (default: 6s)")
	flag.IntVar(&cfg.maxRetry, "max-retry", 3, "Колличество retry на запрос (default: 3)")
	flag.IntVar(&cfg.maxPayloadSize, "max-payload-size", 100, "Ограничение payload по byte на датаграмму (default: 100)")
	flag.BoolVar(&cfg.logDebug, "debug-log", true, "Вывести debug логи")
	flag.Parse()

	// Logger
	var opts slog.HandlerOptions
	if cfg.logDebug {
		opts.Level = slog.LevelDebug
	} else {
		opts.Level = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts))

	if cfg.maxRetry < 0 {
		cfg.maxRetry = 0
	}

	logger.Info("Start with config",
		slog.String("server-addr", cfg.serverAddr),
		slog.Duration("timeout", cfg.timeout),
		slog.Duration("assemble-ttl", cfg.assembleTTL),
		slog.Int("max-retry", cfg.maxPayloadSize),
		slog.Int("max-payload-size", cfg.maxPayloadSize),
		slog.Bool("debug-log", cfg.logDebug),
	)

	application := New(&cfg, logger)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		slog.Error("application stopped with error", slog.Any("err", err))
		os.Exit(1)
	}
}

type App struct {
	cfg *Config
	log *slog.Logger

	packetAssember *protocol.PacketAssembler

	conn          *net.UDPConn
	serverUDPAddr *net.UDPAddr
}

func New(cfg *Config, log *slog.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,

		packetAssember: protocol.NewPackageAssember(cfg.assembleTTL, cfg.maxPayloadSize),
	}
}

func (a *App) Run(ctx context.Context) error {
	serverAddr, err := net.ResolveUDPAddr("udp", a.cfg.serverAddr)
	if err != nil {
		a.log.Error("failed resolve udp", slog.Any("error", err))
		return fmt.Errorf("failed resolve udp")
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		a.log.Error("failed resolve udp", slog.Any("error", err))
		return fmt.Errorf("failed resolve udp")
	}
	a.conn = conn
	a.serverUDPAddr = serverAddr
	defer a.conn.Close()
	defer a.packetAssember.Stop()

	a.log.Info("client starting",
		slog.String("client_addr", a.conn.LocalAddr().String()),
		slog.String("server_addr", serverAddr.String()),
	)
	go func() {
		<-ctx.Done()
		_ = a.conn.Close()
	}()

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
		req := &payload.Request{
			Command: payload.Commands(tokens[0]),
			Args:    tokens[1:],
		}

		resp, err := a.send(ctx, req)
		if err != nil {
			a.log.Error("request failed", slog.Any("error", err))
			continue
		}

		fmt.Printf("[%s] %s\n", string(resp.Status), string(resp.Body))
	}

	return scanner.Err()
}

func (a *App) send(ctx context.Context, req *payload.Request) (*payload.Response, error) {
	reqID := fmt.Sprintf("%d", a.cfg.reqSeq.Add(1))
	payload := req.Marshal()
	packets := protocol.SplitDataOnPackage(reqID, payload, a.cfg.maxPayloadSize)

	for attempt := 1; attempt < a.cfg.maxRetry+1; attempt += 1 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		a.log.Debug("sending request",
			slog.String("req_id", reqID),
			slog.Int("attempt", attempt),
			slog.String("cmd", string(req.Command)),
			slog.Int("packats", len(packets)),
		)

		for _, packet := range packets {
			if _, err := a.conn.Write(packet.Marshal()); err != nil {
				return nil, fmt.Errorf("write udp: %w", err)
			}
		}

		resp, err := a.awaitRespones(reqID)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				a.log.Warn("request timed out, retrying...",
					slog.String("req_id", reqID),
					slog.Int("attempt", attempt),
				)
				continue
			}
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("server unreachable after %d attempts", a.cfg.maxRetry)
}

func (a *App) awaitRespones(msgID string) (*payload.Response, error) {
	deadline := time.Now().Add(a.cfg.timeout)
	if err := a.conn.SetReadDeadline(deadline); err != nil {
		return nil, fmt.Errorf("set deadline: %w", err)
	}

	buf := make([]byte, 2048)
	for {
		n, from, err := a.conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return nil, os.ErrDeadlineExceeded
			}
			return nil, err
		}

		if !from.IP.Equal(a.serverUDPAddr.IP) || from.Port != a.serverUDPAddr.Port {
			continue
		}

		raw := make([]byte, n)
		copy(raw, buf[:n])

		a.log.Debug("udp read datagamm", slog.String("packet", string(raw)))
		respPacket, err := protocol.Unmarshal(raw)
		if err != nil {
			return nil, fmt.Errorf("failed unmarshal packet: %w", err)
		}

		if respPacket.MessageID != msgID {
			continue
		}

		a.log.Debug("packet unmarshal",
			slog.String("msg_id", respPacket.MessageID),
			slog.Int("part", respPacket.Part),
			slog.Int("total", respPacket.Total),
			slog.String("payload", string(respPacket.Payload)),
		)

		respPacketsPayload, ok := a.packetAssember.AddPacakge(respPacket.MessageID, respPacket.Part, respPacket.Total, respPacket.Payload)
		if !ok {
			continue
		}

		respPayload, err := payload.UnmarshalResponse(respPacketsPayload)
		if err != nil {
			return nil, fmt.Errorf("failed unmarshal payload: %w", err)
		}
		return respPayload, nil
	}
}
