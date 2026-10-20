package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sudo-odner/tmp/network/go-udp-protocol/cmd/server/storage"
	"github.com/sudo-odner/tmp/network/go-udp-protocol/pkg/cachelru"
	"github.com/sudo-odner/tmp/network/go-udp-protocol/pkg/protocol"
	"github.com/sudo-odner/tmp/network/go-udp-protocol/pkg/protocol/payload"
)

type Config struct {
	host           string
	port           int
	timeout        time.Duration
	assembleTTL    time.Duration
	lossRate       float64
	maxPayloadSize int
	cacheCapacity  int
}

func main() {
	// Loggerk
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Config
	var cfg Config
	flag.StringVar(&cfg.host, "host", "127.0.0.1", "Айпи сервера (default: '127.0.0.1')")
	flag.IntVar(&cfg.port, "port", 5003, "Порт сервера (defalut: 5008)")
	flag.DurationVar(&cfg.timeout, "timeout", 2*time.Second, "Время ожидания (default: 2s)")
	flag.DurationVar(&cfg.timeout, "assemble-ttl", 6*time.Second, "Время жизни пакетов на сборку пакетов (default: 6s)")
	flag.Float64Var(&cfg.lossRate, "loss-rate", 0.3, "Процент потерь пакетов в коэффициенте от 0-1 (default: 0.3)")
	flag.IntVar(&cfg.maxPayloadSize, "max-payload-size", 100, "Ограничение payload по byte на датаграмму (default: 100)")
	flag.IntVar(&cfg.cacheCapacity, "cache-capacity", 100, "Ограничение LRU cache (default: 100)")
	flag.Parse()

	if 0.0 > cfg.lossRate && cfg.lossRate > 1.0 {
		logger.Error("loss-rate - should be upper zero and less one")
		os.Exit(1)
	}

	logger.Info("Start with config",
		slog.String("host", cfg.host),
		slog.Int("port", cfg.port),
		slog.Duration("timeout", cfg.timeout),
		slog.Float64("loss-rate", cfg.lossRate),
		slog.Int("max-payload-size", cfg.maxPayloadSize),
		slog.Int("cache-capacity", cfg.cacheCapacity),
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
	cfg     *Config
	log     *slog.Logger
	storage *storage.Storage

	packetAssembler *protocol.PacketAssembler
	cache           *cachelru.CacheLRU

	conn *net.UDPConn
}

func New(cfg *Config, log *slog.Logger) *App {
	return &App{
		cfg:     cfg,
		log:     log,
		storage: storage.New(),

		packetAssembler: protocol.NewPackageAssember(cfg.timeout, cfg.maxPayloadSize),
		cache:           cachelru.New(cfg.cacheCapacity),
	}
}

func (a *App) Run(ctx context.Context) error {
	ip := net.ParseIP(a.cfg.host)
	if ip == nil {
		a.log.Error("Incurrect IP-host format", "host", ip)
	}
	fmt.Print(ip, a.cfg.port)
	addr := &net.UDPAddr{
		IP:   ip,
		Port: a.cfg.port,
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		a.log.Error("failed listed udp", slog.Any("error", err))
		return fmt.Errorf("falied listed udp")
	}
	a.conn = conn
	defer a.conn.Close()
	defer a.packetAssembler.Stop()

	a.log.Info("server lissening", slog.String("addr", addr.String()))
	go func() {
		<-ctx.Done()
		_ = a.conn.Close()
	}()

	buffer := make([]byte, 2048)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			a.log.Error("failed read datagramm", slog.Any("error", err))
			continue
		}

		raw := make([]byte, n)
		copy(raw, buffer[:n])
		a.log.Info("read from udp", slog.String("datagramm", string(raw)))
		a.handlerPacket(clientAddr, raw)
	}
}

func (a *App) handlerPacket(clientAddr *net.UDPAddr, raw []byte) {
	reqPacket, err := protocol.Unmarshal(raw)
	if err != nil {
		a.log.Error("falied unmarshal packet", slog.Any("error", err))
		return
	}

	// // Miss datagramm in
	// if rand.Float64() < a.cfg.lossRate {
	// 	a.log.Warn("simulate request drop",
	// 		slog.String("msg_id", packetRequest.MessageID),
	// 		slog.Float64("loss_rate", a.cfg.lossRate),
	// 	)
	// 	return
	// }

	sessionKey := fmt.Sprintf("%s#%s", clientAddr.String(), reqPacket.MessageID)

	var reqPayload []byte
	if reqPacket.Total > 1 {
		completePayload, complete := a.packetAssembler.AddPacakge(
			sessionKey,
			reqPacket.Part,
			reqPacket.Total,
			reqPacket.Payload,
		)
		if !complete {
			return
		}
		reqPayload = completePayload
	} else {
		reqPayload = reqPacket.Payload
	}

	var respPayload []byte
	if cache, ok := a.cache.Get(sessionKey); ok {
		// LRU Cache: Fast path
		a.log.Debug("fast path: read from LRU cache", slog.String("sessionKey", sessionKey))
		respPayload = cache
	} else {
		// Starage & Command: Slow path
		a.log.Debug("slow path: read from storage", slog.String("sessionKey", sessionKey))
		respPayload = a.execute(reqPayload).Marshal()
		a.cache.Put(sessionKey, respPayload)
	}

	respPackets := protocol.SplitDataOnPackage(
		reqPacket.MessageID,
		respPayload,
		a.cfg.maxPayloadSize,
	)

	for _, respPacket := range respPackets {
		// Miss datagramm out
		if rand.Float64() < a.cfg.lossRate {
			a.log.Warn("simulate response drop",
				slog.String("msg_id", reqPacket.MessageID),
				slog.Float64("loss_rate", a.cfg.lossRate),
			)
			return
		}
		if _, err := a.conn.WriteToUDP(respPacket.Marshal(), clientAddr); err != nil {
			a.log.Error("failed to write datagram",
				slog.String("client", clientAddr.String()),
				slog.Any("err", err),
			)
		} else {
			a.log.Debug("sending response",
				slog.String("client", clientAddr.String()),
				slog.String("msg_id", respPacket.MessageID),
				slog.Int("payload_bytes", len(respPacket.Payload)),
				slog.Int("chunk_size", a.cfg.maxPayloadSize),
				slog.Int("total_chunks", len(respPackets)),
			)
		}
	}
}

func (a *App) execute(raw []byte) *payload.Response {
	req, err := payload.UnmarshalRequset(raw)
	if err != nil {
		return &payload.Response{Status: payload.StatusERROR, Body: "BAD_REQUEST"}
	}

	switch strings.ToUpper(string(req.Command)) {
	case "GET":
		if len(req.Args) == 0 {
			return &payload.Response{Status: payload.StatusERROR, Body: "MISSING_ARG"}
		}
		info, ok := a.storage.ByKey(req.Args[0])
		if !ok {
			return &payload.Response{Status: payload.StatusERROR, Body: "NOT_FOUND"}
		}
		return &payload.Response{Status: payload.StatusOK, Body: info}

	case "LIST":
		list := a.storage.List()
		return &payload.Response{Status: payload.StatusOK, Body: strings.Join(list, ";")}

	case "SEARCH":
		if len(req.Args) == 0 {
			return &payload.Response{Status: payload.StatusERROR, Body: "MISSING_ARG"}
		}
		res := a.storage.Search(req.Args[0])
		union := make([]string, len(res))
		for i, v := range res {
			union[i] = fmt.Sprintf("%s: %s", v[0], v[1])
		}
		return &payload.Response{Status: payload.StatusOK, Body: strings.Join(union, ";")}

	case "COUNT":
		count := a.storage.Count()
		return &payload.Response{Status: payload.StatusOK, Body: strconv.Itoa(count)}

	default:
		return &payload.Response{Status: payload.StatusERROR, Body: "UNKNOWN_COMMAND"}
	}
}
