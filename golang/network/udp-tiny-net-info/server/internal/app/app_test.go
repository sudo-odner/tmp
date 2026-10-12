package app

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/server/internal/config"
	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/server/pkg/assembly"
	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/server/pkg/protocol"
)

func TestApp_Integration(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{
		MaxChunkSize: 32, // force chunking for testing
		LossRate:     0,
		Server: config.Server{
			Address: "127.0.0.1:5508",
			Timeout: time.Second,
		},
		Cache: config.Cache{
			Capacity: 10,
		},
	}

	app := New(cfg, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = app.Run(ctx)
	}()

	// Wait for server to bind
	time.Sleep(50 * time.Millisecond)

	serverAddr, err := net.ResolveUDPAddr("udp", cfg.Server.Address)
	if err != nil {
		t.Fatalf("resolve server addr: %v", err)
	}

	clientConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		t.Fatalf("listen client udp: %v", err)
	}
	defer clientConn.Close()

	asm := assembly.NewAssembler(time.Second)
	defer asm.Stop()

	sendAndReceiveMsg := func(customID, cmd string, args ...string) *protocol.Response {
		req := &protocol.Request{
			Command: cmd,
			Args:    args,
		}
		reqPayload := req.Marshal()
		msgID := customID
		if msgID == "" {
			msgID = fmt.Sprintf("test-%d", time.Now().UnixNano())
		}
		chunk := protocol.Chunk{
			MessageID: msgID,
			Part:      1,
			Total:     1,
			Data:      reqPayload,
		}
		if _, err := clientConn.WriteToUDP(chunk.Marshal(), serverAddr); err != nil {
			t.Fatalf("write to udp error: %v", err)
		}

		_ = clientConn.SetReadDeadline(time.Now().Add(time.Second))
		buf := make([]byte, 2048)
		for {
			n, from, err := clientConn.ReadFromUDP(buf)
			if err != nil {
				t.Fatalf("read from udp error: %v", err)
			}

			respChunk, err := protocol.UnmarshalChunk(buf[:n])
			if err != nil {
				t.Fatalf("unmarshal chunk error: %v", err)
			}

			if respChunk.MessageID != msgID {
				continue
			}

			var payload []byte
			if respChunk.Total > 1 {
				sessionKey := fmt.Sprintf("%s#%s", from.String(), respChunk.MessageID)
				full, complete := asm.AddChunk(sessionKey, respChunk.Part, respChunk.Total, respChunk.Data)
				if !complete {
					continue
				}
				payload = full
			} else {
				payload = respChunk.Data
			}

			resp, err := protocol.UnmarshalResponse(payload)
			if err != nil {
				t.Fatalf("unmarshal response error: %v", err)
			}
			return resp
		}
	}

	sendAndReceive := func(cmd string, args ...string) *protocol.Response {
		return sendAndReceiveMsg("", cmd, args...)
	}

	t.Run("GET HTTP", func(t *testing.T) {
		resp := sendAndReceive("GET", "HTTP")
		if resp.Status != protocol.StatusOK {
			t.Fatalf("expected StatusOK, got %s: %s", resp.Status, resp.Body)
		}
	})

	t.Run("GET nonexistent", func(t *testing.T) {
		resp := sendAndReceive("GET", "UNKNOWN_PROTO")
		if resp.Status != protocol.StatusError || resp.Body != "NOT_FOUND" {
			t.Fatalf("expected NOT_FOUND, got %s: %s", resp.Status, resp.Body)
		}
	})

	t.Run("COUNT", func(t *testing.T) {
		resp := sendAndReceive("COUNT")
		if resp.Status != protocol.StatusOK {
			t.Fatalf("expected StatusOK, got %s: %s", resp.Status, resp.Body)
		}
	})

	t.Run("SEARCH web", func(t *testing.T) {
		resp := sendAndReceive("SEARCH", "web")
		if resp.Status != protocol.StatusOK {
			t.Fatalf("expected StatusOK, got %s: %s", resp.Status, resp.Body)
		}
	})

	t.Run("Duplicate request caching", func(t *testing.T) {
		resp1 := sendAndReceiveMsg("dup-123", "GET", "SSH")
		resp2 := sendAndReceiveMsg("dup-123", "GET", "SSH")

		if resp1.Body != resp2.Body {
			t.Fatalf("expected identical cached response, got %s vs %s", resp1.Body, resp2.Body)
		}
	})
}
