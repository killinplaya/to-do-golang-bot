package todo

import (
	"context"
	"errors"
	"sort"
	"strings"
	"tg-bot-go/internal/model"
)

type service struct {
	repo           TaskRepository
	maxActiveTasks int
}

func New(repo TaskRepository) Service {
	return &service{
		repo:           repo,
		maxActiveTasks: DefaultMaxActiveTasks,
	}
}

func (s *service) AddTask(ctx context.Context, userID int64, text string) (model.Task, error) {
	normalized := strings.TrimSpace(text)
	if normalized == "" {
		return model.Task{}, ErrEmptyText
	}
	if len([]rune(normalized)) < MinTaskTextLen {
		return model.Task{}, ErrTextTooShort
	}

	activeStatus := model.StatusActive
	activeTasks, err := s.repo.List(ctx, userID, &activeStatus)
	if err != nil {
		return model.Task{}, err
	}
	if s.maxActiveTasks > 0 && len(activeTasks) >= s.maxActiveTasks {
		return model.Task{}, ErrActiveLimitReached
	}

	return s.repo.Create(ctx, userID, normalized)
}

func (s *service) ListTasks(ctx context.Context, userID int64, status *model.Status) ([]model.Task, error) {
	tasks, err := s.repo.List(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (s *service) MarkDone(ctx context.Context, userID, taskID int64) error {
	if err := s.repo.SetDone(ctx, userID, taskID); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *service) DeleteTask(ctx context.Context, userID, taskID int64) error {
	if err := s.repo.Delete(ctx, userID, taskID); err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (s *service) ClearDone(ctx context.Context, userID int64) (int, error) {
	return s.repo.ClearDone(ctx, userID)
}
