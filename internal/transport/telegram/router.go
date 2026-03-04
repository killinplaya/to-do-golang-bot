package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tg-bot-go/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	commandStart = "start"
	commandHelp  = "help"
)

func (h *Handler) routeUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.Message != nil {
		return h.handleMessage(ctx, update.Message)
	}
	if update.CallbackQuery != nil {
		return h.handleCallback(ctx, update.CallbackQuery)
	}

	return nil
}

func (h *Handler) handleMessage(ctx context.Context, msg *tgbotapi.Message) error {
	if msg == nil {
		return nil
	}
	if msg.From == nil {
		return nil
	}

	userID := msg.From.ID
	chatID := msg.Chat.ID

	if msg.IsCommand() {
		return h.handleCommand(chatID, userID, msg.Command())
	}

	switch msg.Text {
	case buttonAddTask:
		h.states.Set(userID, model.StateWaitingTask)
		return h.sendMessage(chatID, RenderNeedTaskText(), MainKeyboard())
	case buttonListTasks:
		return h.sendTaskList(ctx, chatID, userID, model.StatusActive)
	case buttonDoneTasks:
		return h.sendTaskList(ctx, chatID, userID, model.StatusDone)
	case buttonClearDone:
		removed, err := h.todoService.ClearDone(ctx, userID)
		if err != nil {
			return h.sendMessage(chatID, RenderError(err), MainKeyboard())
		}
		return h.sendMessage(chatID, RenderClearedDone(removed), MainKeyboard())
	}

	state := h.states.Get(userID)
	if state.Name == model.StateWaitingTask {
		task, err := h.todoService.AddTask(ctx, userID, msg.Text)
		if err != nil {
			return h.sendMessage(chatID, RenderError(err), MainKeyboard())
		}

		h.states.Clear(userID)
		return h.sendMessage(chatID, RenderTaskAdded(task), MainKeyboard())
	}

	return h.sendMessage(chatID, RenderUnknownInput(), MainKeyboard())
}

func (h *Handler) handleCommand(chatID, userID int64, command string) error {
	switch strings.ToLower(command) {
	case commandStart:
		h.states.Clear(userID)
		return h.sendMessage(chatID, RenderWelcome(), MainKeyboard())
	case commandHelp:
		return h.sendMessage(chatID, RenderHelp(), MainKeyboard())
	default:
		return h.sendMessage(chatID, "Неизвестная команда. Используй /help.", MainKeyboard())
	}
}

func (h *Handler) sendTaskList(
	ctx context.Context,
	chatID int64,
	userID int64,
	status model.Status,
) error {
	tasks, err := h.todoService.ListTasks(ctx, userID, &status)
	if err != nil {
		return h.sendMessage(chatID, RenderError(err), MainKeyboard())
	}

	msg := tgbotapi.NewMessage(chatID, RenderTaskList(tasks, status))
	if len(tasks) > 0 {
		inline := TaskInlineKeyboard(tasks, status)
		msg.ReplyMarkup = inline
	}

	_, err = h.bot.Send(msg)
	return err
}

func (h *Handler) handleCallback(ctx context.Context, query *tgbotapi.CallbackQuery) error {
	if query == nil {
		return nil
	}
	if query.From == nil {
		return nil
	}

	action, taskID, err := parseCallbackData(query.Data)
	if err != nil {
		_, _ = h.bot.Request(tgbotapi.NewCallback(query.ID, "Некорректные данные callback"))
		return nil
	}

	userID := query.From.ID
	switch action {
	case callbackActionDone:
		err = h.todoService.MarkDone(ctx, userID, taskID)
	case callbackActionDelete:
		err = h.todoService.DeleteTask(ctx, userID, taskID)
	default:
		_, _ = h.bot.Request(tgbotapi.NewCallback(query.ID, "Неизвестное действие"))
		return nil
	}

	if err != nil {
		_, _ = h.bot.Request(tgbotapi.NewCallback(query.ID, RenderError(err)))
		return nil
	}

	_, _ = h.bot.Request(tgbotapi.NewCallback(query.ID, RenderActionResult(action)))

	if query.Message != nil {
		if err := h.refreshCallbackMessage(ctx, query.Message, userID); err != nil {
			h.logger.Warn("failed to refresh callback message", "error", err)
		}
	}

	return nil
}

func (h *Handler) refreshCallbackMessage(
	ctx context.Context,
	msg *tgbotapi.Message,
	userID int64,
) error {
	status := model.StatusActive
	if IsDoneListMessage(msg.Text) {
		status = model.StatusDone
	}

	tasks, err := h.todoService.ListTasks(ctx, userID, &status)
	if err != nil {
		return err
	}

	editCfg := tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, RenderTaskList(tasks, status))
	inline := TaskInlineKeyboard(tasks, status)
	editCfg.ReplyMarkup = &inline

	_, err = h.bot.Send(editCfg)
	return err
}

func (h *Handler) sendMessage(chatID int64, text string, markup interface{}) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if markup != nil {
		msg.ReplyMarkup = markup
	}

	_, err := h.bot.Send(msg)
	return err
}

func parseCallbackData(data string) (action string, taskID int64, err error) {
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return "", 0, errors.New("invalid callback data format")
	}

	parsedID, parseErr := strconv.ParseInt(parts[1], 10, 64)
	if parseErr != nil {
		return "", 0, fmt.Errorf("parse task id: %w", parseErr)
	}
	if parsedID <= 0 {
		return "", 0, errors.New("invalid task id")
	}

	return parts[0], parsedID, nil
}
