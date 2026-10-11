package main

import (
	"log/slog"
	"os"

	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/server/internal/config"
)

func main() {
	cfg := config.MustLoad()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	// TODO: init App + GC
}
