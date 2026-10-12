package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"

	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/server/internal/config"
	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/server/internal/repository/cache"
	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/server/internal/repository/storage"
	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/server/pkg/assembly"
	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/server/pkg/protocol"
)

// Datagram in UTF-8 string format
//
// If short message less 64 byte:
// Request:		ReqID|CMD|ARG1|ARG2...
// Response:	ReqID|OK|Result or ReqID|ERROR|Code
//
// Else, send parts:
// ReqID|Part|Total|Data

type App struct {
	cfg      *config.Config
	log      *slog.Logger
	repo     *storage.Storage
	cache    *cache.LRUCache
	assembly *assembly.Assembler
	conn     *net.UDPConn
}

func New(cfg *config.Config, log *slog.Logger) *App {
	return &App{
		cfg:      cfg,
		log:      log,
		repo:     storage.New(),
		cache:    cache.New(cfg.Cache.Capacity),
		assembly: assembly.NewAssembler(cfg.Server.Timeout),
	}
}

func (a *App) Run(ctx context.Context) error {
	// Set up UDP Server config
	addr, err := net.ResolveUDPAddr("udp", a.cfg.Server.Address)
	if err != nil {
		return fmt.Errorf("falied resolve UDP addr: %w", err)
	}

	// Open connection and Listen UPD
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("falied listen udp: %w", err)
	}
	a.conn = conn
	defer a.conn.Close()
	defer a.assembly.Stop()

	a.log.Info("server listening", slog.String("addr", addr.String()))
	go func() {
		<-ctx.Done()
		_ = a.conn.Close()
	}()

	buf := make([]byte, 2048)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			a.log.Error("read datagram error", slog.Any("error", err))
			continue
		}

		rawPacket := make([]byte, n)
		copy(rawPacket, buf[:n])
		a.handlePacket(clientAddr, rawPacket)
	}
}

func (a *App) handlePacket(clientAddr *net.UDPAddr, rawPacket []byte) {
	chunk, err := protocol.UnmarshalChunk(rawPacket)
	if err != nil {
		a.log.Warn("malformed packet", slog.String("from", clientAddr.String()), slog.Any("err", err))
		return
	}

	sessionKey := fmt.Sprintf("%s#%s", clientAddr.String(), chunk.MessageID)
	var fullPayload []byte

	if chunk.Total > 1 {
		payload, complete := a.assembly.AddChunk(sessionKey, chunk.Part, chunk.Total, chunk.Data)
		if !complete {
			return
		}
		fullPayload = payload
	} else {
		fullPayload = chunk.Data
	}

	if cachedResp, ok := a.cache.Get(sessionKey); ok {
		a.log.Info("duplicate request served from cache", slog.String("req_id", chunk.MessageID))
		a.sendResponse(clientAddr, chunk.MessageID, []byte(cachedResp))
		return
	}

	resp := a.executeCommand(fullPayload)
	respBytes := resp.Marshal()

	a.cache.Put(sessionKey, string(respBytes))

	// miss datagram
	if a.cfg.LossRate > 0 && rand.Float64() < a.cfg.LossRate {
		a.log.Warn("simulate response drop", slog.String("req_id", chunk.MessageID), slog.Float64("loss_rate", a.cfg.LossRate))
		return
	}

	a.sendResponse(clientAddr, chunk.MessageID, respBytes)
}

func (a *App) executeCommand(payload []byte) *protocol.Response {
	req, err := protocol.UnmarshalRequest(payload)
	if err != nil {
		return &protocol.Response{Status: protocol.StatusError, Body: "BAD_REQUEST"}
	}

	switch strings.ToUpper(req.Command) {
	case "GET":
		if len(req.Args) == 0 {
			return &protocol.Response{Status: protocol.StatusError, Body: "MISSING_ARG"}
		}
		info, ok := a.repo.ByKey(req.Args[0])
		if !ok {
			return &protocol.Response{Status: protocol.StatusError, Body: "NOT_FOUND"}
		}
		return &protocol.Response{Status: protocol.StatusOK, Body: info}

	case "LIST":
		list := a.repo.List()
		return &protocol.Response{Status: protocol.StatusOK, Body: strings.Join(list, ";")}

	case "SEARCH":
		if len(req.Args) == 0 {
			return &protocol.Response{Status: protocol.StatusError, Body: "MISSING_ARG"}
		}
		res := a.repo.Search(req.Args[0])
		union := make([]string, len(res))
		for i, v := range res {
			union[i] = fmt.Sprintf("%s: %s", v[0], v[1])
		}
		return &protocol.Response{Status: protocol.StatusOK, Body: strings.Join(union, ";")}

	case "COUNT":
		count := a.repo.Count()
		return &protocol.Response{Status: protocol.StatusOK, Body: strconv.Itoa(count)}

	default:
		return &protocol.Response{Status: protocol.StatusError, Body: "UNKNOWN_COMMAND"}
	}
}

func (a *App) sendResponse(clientAddr *net.UDPAddr, msgID string, payload []byte) {
	chunks := protocol.SplitIntoChunks(msgID, payload, a.cfg.MaxChunkSize)
	a.log.Debug("sending response",
		slog.String("client", clientAddr.String()),
		slog.String("msg_id", msgID),
		slog.Int("payload_bytes", len(payload)),
		slog.Int("chunk_size", a.cfg.MaxChunkSize),
		slog.Int("total_chunks", len(chunks)),
	)
	for _, ch := range chunks {
		if _, err := a.conn.WriteToUDP(ch.Marshal(), clientAddr); err != nil {
			a.log.Error("failed to write datagram", slog.String("client", clientAddr.String()), slog.Any("err", err))
		}
	}
}
