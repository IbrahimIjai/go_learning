package models

import (
	"context"
	"errors"
	"sync"
)

// Sentinel errors. Callers check these with errors.Is, letting the controller
// layer map them to HTTP status codes without the store knowing about HTTP.
var (
	ErrNotFound = errors.New("store: code not found")
	ErrConflict = errors.New("store: code already exists")
)

// URLStore is everything the rest of the app is allowed to know about
// persistence.
type URLStore interface {
	Save(ctx context.Context, code, target string) error
	GetByCode(ctx context.Context, code string) (string, error)
}

// Memory is a thread-safe in-memory URLStore. The mutex guards the map against
// concurrent access from per-request goroutines.
type Memory struct {
	mu   sync.RWMutex
	urls map[string]string
}

var _ URLStore = (*Memory)(nil)

func NewMemory() *Memory {
	return &Memory{urls: make(map[string]string)}
}

func (m *Memory) Save(ctx context.Context, code, target string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.urls[code]; exists {
		return ErrConflict
	}
	m.urls[code] = target
	return nil
}

func (m *Memory) GetByCode(ctx context.Context, code string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	target, ok := m.urls[code]
	if !ok {
		return "", ErrNotFound
	}
	return target, nil
}
