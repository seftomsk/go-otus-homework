package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	// Register some standard stuff.
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage"
)

type InPostgres struct {
	db *sql.DB
}

func New(user, pwd, host, port, dbname string) (*InPostgres, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, pwd, host, port, dbname)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("new %w", err)
	}

	return &InPostgres{
		db: db,
	}, nil
}

func (s *InPostgres) getBy(ctx context.Context, field string, value interface{}) (app.Event, error) {
	var query string
	switch field {
	case "id":
		query = `SELECT id, user_id, title, description, datetime, duration, remind_before
		FROM events WHERE id = $1`
	case "datetime":
		query = `SELECT id, user_id, title, description, datetime, duration, remind_before
		FROM events WHERE datetime = $1`
	default:
		return nil, fmt.Errorf("getBy - invalid field value")
	}

	row := s.db.QueryRowContext(ctx, query, value)
	var id string
	var userID string
	var title string
	var tempDescription sql.NullString
	var datetime time.Time
	var duration time.Duration
	var remindBefore time.Duration
	var tempDuration int
	var tempRemindBefore int
	err := row.Scan(&id, &userID, &title, &tempDescription, &datetime, &tempDuration, &tempRemindBefore)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getBy: %w", storage.ErrNotFound)
		}
		return nil, fmt.Errorf("getBy: %w", err)
	}

	duration = time.Duration(tempDuration) * time.Millisecond
	remindBefore = time.Duration(tempRemindBefore) * time.Millisecond

	description := storage.EventDescription{}
	if tempDescription.Valid {
		description = storage.EventDescription{Data: tempDescription.String, Valid: true}
	}

	return storage.NewEvent(id, userID, title, description, datetime, duration, remindBefore), nil
}

func (s *InPostgres) Connect(ctx context.Context) error {
	if s.db == nil {
		return storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	if err := goose.Up(s.db, "migrations"); err != nil {
		return fmt.Errorf("connect - up migration: %w", err)
	}
	return nil
}

func (s *InPostgres) Close() error {
	if s.db == nil {
		return nil
	}

	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	return nil
}

func (s *InPostgres) GetByDatetime(ctx context.Context, date time.Time) (app.Event, error) {
	if s.db == nil {
		return nil, storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	res, err := s.getBy(ctx, "datetime", date)
	if err != nil {
		return nil, fmt.Errorf("getByDatetime: %w", err)
	}
	return res, nil
}

func (s *InPostgres) GetByID(ctx context.Context, id string) (app.Event, error) {
	if s.db == nil {
		return nil, storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	res, err := s.getBy(ctx, "id", id)
	if err != nil {
		return nil, fmt.Errorf("getById: %w", err)
	}
	return res, nil
}

func (s *InPostgres) Create(ctx context.Context, event app.Event) error {
	if s.db == nil {
		return storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	query := `INSERT INTO events(id, user_id, title, description, datetime, duration, remind_before)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	var err error
	if event.Description().Valid {
		_, err = s.db.ExecContext(
			ctx,
			query,
			event.ID(),
			event.UserID(),
			event.Title(),
			event.Description().Data,
			event.Datetime(),
			event.Duration().Milliseconds(),
			event.RemindBefore().Milliseconds())
	} else {
		_, err = s.db.ExecContext(
			ctx,
			query,
			event.ID(),
			event.UserID(),
			event.Title(),
			nil,
			event.Datetime(),
			event.Duration().Milliseconds(),
			event.RemindBefore().Milliseconds())
	}
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	return nil
}

func (s *InPostgres) Update(ctx context.Context, id string, event app.Event) error {
	if s.db == nil {
		return storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	query := `UPDATE events SET (user_id, title, description, datetime, duration, remind_before)
		= ($1, $2, $3, $4, $5, $6) WHERE id = $7`
	var err error
	var row sql.Result
	if event.Description().Valid {
		row, err = s.db.ExecContext(
			ctx,
			query,
			event.UserID(),
			event.Title(),
			event.Description().Data,
			event.Datetime(),
			event.Duration().Milliseconds(),
			event.RemindBefore().Milliseconds(),
			id)
	} else {
		row, err = s.db.ExecContext(
			ctx,
			query,
			event.UserID(),
			event.Title(),
			nil,
			event.Datetime(),
			event.Duration().Milliseconds(),
			event.RemindBefore().Milliseconds(),
			id)
	}
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	n, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("update: %w", storage.ErrNotFound)
	}

	return nil
}

func (s *InPostgres) Delete(ctx context.Context, id string) error {
	if s.db == nil {
		return storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	query := `DELETE FROM events WHERE id = $1`
	row, err := s.db.ExecContext(
		ctx,
		query,
		id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	n, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	if n != 1 {
		return fmt.Errorf("delete: %w", storage.ErrNotFound)
	}
	return nil
}

func (s *InPostgres) GetByDateRange(ctx context.Context, start, end time.Time) (map[string]app.Event, error) {
	if s.db == nil {
		return nil, storage.ErrInvalidInitialization
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	events := make(map[string]app.Event)

	if start.After(end) {
		return nil, fmt.Errorf("getByDateRange: %w", storage.ErrStartDateMoreThanEndDate)
	}

	query := `SELECT id, user_id, title, description, datetime, duration, remind_before
		FROM events WHERE date_trunc('day', datetime) BETWEEN $1 AND $2`
	rows, err := s.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("getByDateRange: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	var id string
	var userID string
	var title string
	var tempDescription sql.NullString
	var datetime time.Time
	var duration time.Duration
	var remindBefore time.Duration
	var tempDuration int
	var tempRemindBefore int

	for rows.Next() {
		err = rows.Scan(&id, &userID, &title, &tempDescription, &datetime, &tempDuration, &tempRemindBefore)
		if err != nil {
			return nil, fmt.Errorf("getByDateRange: %w", err)
		}

		duration = time.Duration(tempDuration) * time.Millisecond
		remindBefore = time.Duration(tempRemindBefore) * time.Millisecond

		description := storage.EventDescription{}
		if tempDescription.Valid {
			description = storage.EventDescription{Data: tempDescription.String, Valid: true}
		}

		events[id] = storage.NewEvent(id, userID, title, description, datetime, duration, remindBefore)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("getByDateRange: %w", err)
	}

	if len(events) == 0 {
		return nil, storage.ErrNotFound
	}

	return events, nil
}
