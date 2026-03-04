package telegram

import (
	"errors"
	"fmt"
	"strings"
	"tg-bot-go/internal/model"
	"tg-bot-go/internal/service/todo"
	"time"
)

const (
	activeListHeader = "📋 Активные задачи"
	doneListHeader   = "✅ Выполненные задачи"
)

func RenderWelcome() string {
	return strings.Join([]string{
		"Привет! Я ToDo-бот.",
		"Помогаю хранить задачи прямо в Telegram.",
		"",
		"Используй кнопки внизу или команду /help.",
	}, "\n")
}

func RenderHelp() string {
	return strings.Join([]string{
		"Справка:",
		"/start — приветствие и меню",
		"/help — показать эту справку",
		"",
		"Кнопки:",
		"➕ Add task — добавить задачу",
		"📋 List tasks — показать активные задачи",
		"✅ Done tasks — показать выполненные задачи",
		"🧹 Clear done — удалить все выполненные задачи",
	}, "\n")
}

func RenderTaskAdded(task model.Task) string {
	return fmt.Sprintf("Добавлено ✅\n#%d: %s", task.ID, task.Text)
}

func RenderClearedDone(removed int) string {
	if removed == 0 {
		return "Выполненных задач для очистки нет."
	}
	return fmt.Sprintf("Готово: удалено %d выполненных задач.", removed)
}

func RenderTaskList(tasks []model.Task, status model.Status) string {
	if status == model.StatusDone {
		return renderDoneTasks(tasks)
	}
	return renderActiveTasks(tasks)
}

func RenderUnknownInput() string {
	return "Не понял сообщение. Используй кнопки меню или /help."
}

func RenderNeedTaskText() string {
	return "Ок, напиши текст задачи."
}

func RenderActionResult(action string) string {
	switch action {
	case callbackActionDone:
		return "Задача отмечена как выполненная ✅"
	case callbackActionDelete:
		return "Задача удалена 🗑"
	default:
		return "Готово"
	}
}

func RenderError(err error) string {
	switch {
	case errors.Is(err, todo.ErrEmptyText):
		return "Текст задачи пустой. Введи нормальный текст."
	case errors.Is(err, todo.ErrTextTooShort):
		return fmt.Sprintf("Текст слишком короткий. Минимум %d символа.", todo.MinTaskTextLen)
	case errors.Is(err, todo.ErrNotFound):
		return "Задача не найдена."
	case errors.Is(err, todo.ErrActiveLimitReached):
		return fmt.Sprintf("Достигнут лимит активных задач (%d).", todo.DefaultMaxActiveTasks)
	default:
		return "Произошла ошибка. Попробуй ещё раз."
	}
}

func IsDoneListMessage(text string) bool {
	return strings.HasPrefix(text, doneListHeader)
}

func renderActiveTasks(tasks []model.Task) string {
	if len(tasks) == 0 {
		return activeListHeader + "\n\nПока пусто."
	}

	lines := make([]string, 0, len(tasks)+2)
	lines = append(lines, activeListHeader, "")
	for i, task := range tasks {
		lines = append(lines, fmt.Sprintf("%d. #%d — %s", i+1, task.ID, task.Text))
	}

	return strings.Join(lines, "\n")
}

func renderDoneTasks(tasks []model.Task) string {
	if len(tasks) == 0 {
		return doneListHeader + "\n\nПока пусто."
	}

	lines := make([]string, 0, len(tasks)+2)
	lines = append(lines, doneListHeader, "")
	for i, task := range tasks {
		doneAt := "-"
		if task.DoneAt != nil {
			doneAt = task.DoneAt.Format(time.DateTime)
		}
		lines = append(lines, fmt.Sprintf("%d. #%d — %s (done: %s)", i+1, task.ID, task.Text, doneAt))
	}

	return strings.Join(lines, "\n")
}
