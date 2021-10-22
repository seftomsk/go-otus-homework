package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage"
)

var ErrDateBusy = errors.New("date busy")

var ErrNotFound = errors.New("not found")

var ErrStartDateMoreThanEndDate = errors.New("start date more than end date")

var ErrInvalidInitialization = errors.New("start date more than end date")

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Event interface {
	ID() string
	UserID() string
	Title() string
	Datetime() time.Time
	Duration() time.Duration
	Description() storage.EventDescription
	RemindBefore() time.Duration
}

type Storage interface {
	GetByID(ctx context.Context, id string) (Event, error)
	GetByDatetime(ctx context.Context, date time.Time) (Event, error)
	Create(ctx context.Context, event Event) error
	Update(ctx context.Context, id string, event Event) error
	Delete(ctx context.Context, id string) error
	GetByDateRange(ctx context.Context, start, end time.Time) (map[string]Event, error)
}

type EventDTO struct {
	UserID       string
	Title        string
	Description  storage.EventDescription
	Datetime     time.Time
	Duration     time.Duration
	RemindBefore time.Duration
}

type UpdateEventDTO struct {
	ID           string
	UserID       string
	Title        string
	Description  storage.EventDescription
	Datetime     time.Time
	Duration     time.Duration
	RemindBefore time.Duration
}

type App struct {
	logger  Logger
	storage Storage
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) eventsByDateRange(ctx context.Context, start, end time.Time) (map[string]Event, error) {
	events, err := a.storage.GetByDateRange(ctx, start, end)
	if err != nil {
		if errors.Is(err, storage.ErrStartDateMoreThanEndDate) {
			return nil, fmt.Errorf("eventsByDateRange: %w", ErrStartDateMoreThanEndDate)
		}
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("eventsByDateRange: %w", ErrNotFound)
		}

		return nil, fmt.Errorf("eventsByDateRange: %w", err)
	}
	return events, nil
}

func (a *App) CreateEvent(ctx context.Context, dto EventDTO) (Event, error) {
	if a.storage == nil {
		return nil, ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	_, err := a.storage.GetByDatetime(ctx, dto.Datetime)

	if err == nil {
		return nil, fmt.Errorf("createEvent: %w", ErrDateBusy)
	}

	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, fmt.Errorf("createEvent: %w", err)
	}

	event := storage.NewEvent(
		uuid.NewString(),
		dto.UserID,
		dto.Title,
		dto.Description,
		dto.Datetime,
		dto.Duration,
		dto.RemindBefore)

	err = a.storage.Create(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("createEvent: %w", err)
	}
	return event, nil
}

func (a *App) UpdateEvent(ctx context.Context, dto UpdateEventDTO) (Event, error) {
	if a.storage == nil {
		return nil, ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	e, err := a.storage.GetByID(ctx, dto.ID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("updateEvent: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("updateEvent: %w", err)
	}

	if e.Datetime() != dto.Datetime {
		_, err = a.storage.GetByDatetime(ctx, dto.Datetime)
		if err == nil {
			return nil, fmt.Errorf("updateEvent: %w", ErrDateBusy)
		}

		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("updateEvent: %w", err)
		}
	}

	event := storage.NewEvent(
		dto.ID,
		dto.UserID,
		dto.Title,
		dto.Description,
		dto.Datetime,
		dto.Duration,
		dto.RemindBefore)

	err = a.storage.Update(ctx, dto.ID, event)

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, fmt.Errorf("updateEvent: %w", ErrNotFound)
		}
		return nil, fmt.Errorf("updateEvent: %w", err)
	}
	return event, nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	if a.storage == nil {
		return ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	err := a.storage.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return fmt.Errorf("deleteEvent: %w", ErrNotFound)
		}
		return fmt.Errorf("deleteEvent: %w", err)
	}
	return nil
}

func (a *App) EventsForDay(ctx context.Context, date time.Time) (map[string]Event, error) {
	if a.storage == nil {
		return nil, ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	start := date.Truncate(time.Hour * 24)

	return a.eventsByDateRange(ctx, start, start)
}

func (a *App) EventsForWeek(ctx context.Context, date time.Time) (map[string]Event, error) {
	if a.storage == nil {
		return nil, ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	start := date.Truncate(time.Hour * 24)
	if start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, -int(start.Weekday()))
	}
	end := start.AddDate(0, 0, 6)

	return a.eventsByDateRange(ctx, start, end)
}

func (a *App) EventsForMonth(ctx context.Context, date time.Time) (map[string]Event, error) {
	if a.storage == nil {
		return nil, ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	start := date.Truncate(time.Hour * 24)
	if start.Day() != 1 {
		start = start.AddDate(0, 0, -start.Day()+1)
	}
	end := start.AddDate(0, 1, -1)

	return a.eventsByDateRange(ctx, start, end)
}
