package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage"
)

type InMemory struct {
	mu     sync.Mutex
	events map[string]app.Event
}

func New() *InMemory {
	return &InMemory{
		events: make(map[string]app.Event),
	}
}

func (m *InMemory) GetByDatetime(ctx context.Context, date time.Time) (app.Event, error) {
	if m.events == nil {
		return nil, storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, event := range m.events {
		if event.Datetime() == date {
			return event, nil
		}
	}
	return nil, storage.ErrNotFound
}

func (m *InMemory) GetByID(ctx context.Context, id string) (app.Event, error) {
	if m.events == nil {
		return nil, storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if e, ok := m.events[id]; ok {
		return e, nil
	}

	return nil, storage.ErrNotFound
}

func (m *InMemory) Create(ctx context.Context, event app.Event) error {
	if m.events == nil {
		return storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.events[event.ID()] = event
	return nil
}

func (m *InMemory) Update(ctx context.Context, id string, event app.Event) error {
	if m.events == nil {
		return storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.events[id]; !ok {
		return fmt.Errorf("update: %w", storage.ErrNotFound)
	}

	m.events[id] = event
	return nil
}

func (m *InMemory) Delete(ctx context.Context, id string) error {
	if m.events == nil {
		return storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.events[id]; !ok {
		return fmt.Errorf("delete: %w", storage.ErrNotFound)
	}

	delete(m.events, id)
	return nil
}

func (m *InMemory) GetByDateRange(ctx context.Context, start, end time.Time) (map[string]app.Event, error) {
	if m.events == nil {
		return nil, storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	events := make(map[string]app.Event)

	if start.After(end) {
		return nil, fmt.Errorf("getByDateRange: %w", storage.ErrStartDateMoreThanEndDate)
	}

	start = start.AddDate(0, 0, -1)
	end = end.AddDate(0, 0, 1)

	for _, event := range m.events {
		if event.Datetime().After(start) && event.Datetime().Before(end) {
			events[event.ID()] = event
		}
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("getByDateRange: %w", storage.ErrNotFound)
	}
	return events, nil
}
