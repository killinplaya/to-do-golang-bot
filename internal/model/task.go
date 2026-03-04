package model

import (
	"errors"
	"time"
)

type Status string

const (
	StatusActive Status = "active"
	StatusDone   Status = "done"
)

type Task struct {
	ID        int64
	UserID    int64
	Text      string
	Status    Status
	CreatedAt time.Time
	DoneAt    *time.Time
}

type DialogState string

const (
	StateIdle        DialogState = "idle"
	StateWaitingTask DialogState = "waiting_task_text"
)

type State struct {
	UserID int64
	Name   DialogState
}

var ErrTaskNotFound = errors.New("task not found")
