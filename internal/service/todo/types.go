package todo

import (
	"context"
	"errors"
	"tg-bot-go/internal/model"
)

var (
	ErrEmptyText          = errors.New("task text is empty")
	ErrTextTooShort       = errors.New("task text is too short")
	ErrNotFound           = errors.New("task not found")
	ErrActiveLimitReached = errors.New("active task limit reached")
)

const (
	MinTaskTextLen        = 2
	DefaultMaxActiveTasks = 100
)

type TaskRepository interface {
	Create(ctx context.Context, userID int64, text string) (model.Task, error)
	List(ctx context.Context, userID int64, status *model.Status) ([]model.Task, error)
	SetDone(ctx context.Context, userID, taskID int64) error
	Delete(ctx context.Context, userID, taskID int64) error
	ClearDone(ctx context.Context, userID int64) (int, error)
}

type Service interface {
	AddTask(ctx context.Context, userID int64, text string) (model.Task, error)
	ListTasks(ctx context.Context, userID int64, status *model.Status) ([]model.Task, error)
	MarkDone(ctx context.Context, userID, taskID int64) error
	DeleteTask(ctx context.Context, userID, taskID int64) error
	ClearDone(ctx context.Context, userID int64) (int, error)
}
