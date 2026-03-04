package telegram

import (
	"fmt"
	"tg-bot-go/internal/model"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	buttonAddTask   = "➕ Add task"
	buttonListTasks = "📋 List tasks"
	buttonDoneTasks = "✅ Done tasks"
	buttonClearDone = "🧹 Clear done"

	callbackActionDone   = "done"
	callbackActionDelete = "delete"
)

func MainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttonAddTask),
			tgbotapi.NewKeyboardButton(buttonListTasks),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(buttonDoneTasks),
			tgbotapi.NewKeyboardButton(buttonClearDone),
		),
	)
	keyboard.ResizeKeyboard = true

	return keyboard
}

func TaskInlineKeyboard(tasks []model.Task, status model.Status) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(tasks))

	for _, task := range tasks {
		deleteButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("🗑 delete #%d", task.ID),
			fmt.Sprintf("%s:%d", callbackActionDelete, task.ID),
		)

		if status == model.StatusDone {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(deleteButton))
			continue
		}

		doneButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("✅ done #%d", task.ID),
			fmt.Sprintf("%s:%d", callbackActionDone, task.ID),
		)

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(doneButton, deleteButton))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
