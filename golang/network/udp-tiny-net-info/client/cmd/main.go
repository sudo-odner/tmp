package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/client/internal/app"
	"github.com/sudo-odner/tmp/golang/network/udp-tiny-net-info/client/internal/config"
)

func main() {
	cfg := config.MustLoad()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("falied init application", slog.Any("error", err))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		slog.Error("application stopped with error", slog.Any("err", err))
		os.Exit(1)
	}
}
