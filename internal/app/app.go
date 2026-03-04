package app

import (
	"context"
	"log/slog"
	"tg-bot-go/internal/config"
	"tg-bot-go/internal/repository/memory"
	"tg-bot-go/internal/service/todo"
	"tg-bot-go/internal/transport/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type App struct {
	handler *telegram.Handler
}

func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = cfg.LogLevel <= slog.LevelDebug

	taskRepo := memory.NewTaskRepository()
	stateRepo := memory.NewStateRepository()
	todoService := todo.New(taskRepo)

	handler := telegram.NewHandler(bot, todoService, stateRepo, logger)

	return &App{
		handler: handler,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.handler.Start(ctx)
}
