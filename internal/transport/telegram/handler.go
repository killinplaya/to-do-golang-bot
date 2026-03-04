package telegram

import (
	"context"
	"log/slog"
	"tg-bot-go/internal/model"
	"tg-bot-go/internal/service/todo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StateStorage interface {
	Get(userID int64) model.State
	Set(userID int64, state model.DialogState)
	Clear(userID int64)
}

type Handler struct {
	bot         *tgbotapi.BotAPI
	todoService todo.Service
	states      StateStorage
	logger      *slog.Logger
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	todoService todo.Service,
	states StateStorage,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		bot:         bot,
		todoService: todoService,
		states:      states,
		logger:      logger,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 30

	updates := h.bot.GetUpdatesChan(updateCfg)
	defer h.bot.StopReceivingUpdates()

	h.logger.Info("telegram update loop started")

	for {
		select {
		case <-ctx.Done():
			h.logger.Info("telegram update loop stopped by context")
			return nil
		case update, ok := <-updates:
			if !ok {
				h.logger.Info("telegram update channel is closed")
				return nil
			}

			if err := h.routeUpdate(ctx, update); err != nil {
				h.logger.Error("failed to process update", "error", err)
			}
		}
	}
}
