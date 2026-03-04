package memory

import (
	"sync"
	"tg-bot-go/internal/model"
)

type StateRepository struct {
	mu     sync.RWMutex
	states map[int64]model.State
}

func NewStateRepository() *StateRepository {
	return &StateRepository{
		states: make(map[int64]model.State),
	}
}

func (r *StateRepository) Get(userID int64) model.State {
	r.mu.RLock()
	defer r.mu.RUnlock()

	state, ok := r.states[userID]
	if !ok {
		return model.State{
			UserID: userID,
			Name:   model.StateIdle,
		}
	}

	return state
}

func (r *StateRepository) Set(userID int64, state model.DialogState) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.states[userID] = model.State{
		UserID: userID,
		Name:   state,
	}
}

func (r *StateRepository) Clear(userID int64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.states, userID)
}
