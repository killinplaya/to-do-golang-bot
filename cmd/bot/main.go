package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"tg-bot-go/internal/app"
	"tg-bot-go/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := newLogger(cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}

	logger.Info("starting bot")
	if err := application.Run(ctx); err != nil {
		logger.Error("bot stopped with error", "error", err)
		os.Exit(1)
	}
}

func newLogger(level slog.Level) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	return slog.New(handler)
}
