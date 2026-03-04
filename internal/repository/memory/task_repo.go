package memory

import (
	"context"
	"sync"
	"sync/atomic"
	"tg-bot-go/internal/model"
	"time"
)

type TaskRepository struct {
	mu          sync.RWMutex
	tasksByUser map[int64]map[int64]model.Task
	nextID      atomic.Int64
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasksByUser: make(map[int64]map[int64]model.Task),
	}
}

func (r *TaskRepository) Create(_ context.Context, userID int64, text string) (model.Task, error) {
	taskID := r.nextID.Add(1)
	task := model.Task{
		ID:        taskID,
		UserID:    userID,
		Text:      text,
		Status:    model.StatusActive,
		CreatedAt: time.Now(),
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasksByUser[userID]; !ok {
		r.tasksByUser[userID] = make(map[int64]model.Task)
	}
	r.tasksByUser[userID][taskID] = task

	return task, nil
}

func (r *TaskRepository) List(_ context.Context, userID int64, status *model.Status) ([]model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userTasks, ok := r.tasksByUser[userID]
	if !ok {
		return []model.Task{}, nil
	}

	result := make([]model.Task, 0, len(userTasks))
	for _, task := range userTasks {
		if status != nil && task.Status != *status {
			continue
		}
		result = append(result, task)
	}

	return result, nil
}

func (r *TaskRepository) SetDone(_ context.Context, userID, taskID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userTasks, ok := r.tasksByUser[userID]
	if !ok {
		return model.ErrTaskNotFound
	}

	task, ok := userTasks[taskID]
	if !ok {
		return model.ErrTaskNotFound
	}

	if task.Status == model.StatusDone {
		return nil
	}

	now := time.Now()
	task.Status = model.StatusDone
	task.DoneAt = &now
	userTasks[taskID] = task

	return nil
}

func (r *TaskRepository) Delete(_ context.Context, userID, taskID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userTasks, ok := r.tasksByUser[userID]
	if !ok {
		return model.ErrTaskNotFound
	}

	if _, ok := userTasks[taskID]; !ok {
		return model.ErrTaskNotFound
	}

	delete(userTasks, taskID)
	if len(userTasks) == 0 {
		delete(r.tasksByUser, userID)
	}

	return nil
}

func (r *TaskRepository) ClearDone(_ context.Context, userID int64) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userTasks, ok := r.tasksByUser[userID]
	if !ok {
		return 0, nil
	}

	deleted := 0
	for taskID, task := range userTasks {
		if task.Status == model.StatusDone {
			delete(userTasks, taskID)
			deleted++
		}
	}

	if len(userTasks) == 0 {
		delete(r.tasksByUser, userID)
	}

	return deleted, nil
}
