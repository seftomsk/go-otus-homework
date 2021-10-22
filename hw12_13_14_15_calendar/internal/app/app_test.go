package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testLogger struct{}

func (l testLogger) Info(_ string) {}

func (l testLogger) Error(_ string) {}

func (l testLogger) Warn(_ string) {}

func (l testLogger) Debug(_ string) {}

type AppSuite struct {
	suite.Suite
	app *app.App
	ctx context.Context
}

func (e *AppSuite) SetupTest() {
	e.app = app.New(testLogger{}, memory.New())
	e.ctx = context.Background()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}

func (e *AppSuite) TestInvalidInitializationGetErr() {
	a := app.App{}
	createDto := app.EventDTO{}
	dto := app.UpdateEventDTO{}
	now := time.Now().UTC()

	_, err := a.CreateEvent(e.ctx, createDto)
	require.ErrorIs(e.T(), err, app.ErrInvalidInitialization)

	_, err = a.UpdateEvent(e.ctx, dto)
	require.ErrorIs(e.T(), err, app.ErrInvalidInitialization)

	err = a.DeleteEvent(e.ctx, "")
	require.ErrorIs(e.T(), err, app.ErrInvalidInitialization)

	_, err = a.EventsForDay(e.ctx, now)
	require.ErrorIs(e.T(), err, app.ErrInvalidInitialization)

	_, err = a.EventsForWeek(e.ctx, now)
	require.ErrorIs(e.T(), err, app.ErrInvalidInitialization)

	_, err = a.EventsForMonth(e.ctx, now)
	require.ErrorIs(e.T(), err, app.ErrInvalidInitialization)
}

func (e *AppSuite) TestGetDoneFromContextGetErr() {
	createDto := app.EventDTO{}
	dto := app.UpdateEventDTO{}
	now := time.Now().UTC()

	e.T().Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := e.app.CreateEvent(ctx, createDto)
		require.ErrorIs(e.T(), err, context.Canceled)

		_, err = e.app.UpdateEvent(ctx, dto)
		require.ErrorIs(e.T(), err, context.Canceled)

		err = e.app.DeleteEvent(ctx, "")
		require.ErrorIs(e.T(), err, context.Canceled)

		_, err = e.app.EventsForDay(ctx, now)
		require.ErrorIs(e.T(), err, context.Canceled)

		_, err = e.app.EventsForWeek(ctx, now)
		require.ErrorIs(e.T(), err, context.Canceled)

		_, err = e.app.EventsForMonth(ctx, now)
		require.ErrorIs(e.T(), err, context.Canceled)
	})

	e.T().Run("deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*0)
		defer cancel()

		_, err := e.app.CreateEvent(ctx, createDto)
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		_, err = e.app.UpdateEvent(ctx, dto)
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		err = e.app.DeleteEvent(ctx, "")
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		_, err = e.app.EventsForDay(ctx, now)
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		_, err = e.app.EventsForWeek(ctx, now)
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		_, err = e.app.EventsForMonth(ctx, now)
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)
	})
}

func (e *AppSuite) TestCreateEventWithoutErr() {
	testCases := []struct {
		name string
		in   app.EventDTO
	}{
		{
			name: "1",
			in: app.EventDTO{
				UserID:       "",
				Title:        "test",
				Description:  storage.EventDescription{},
				Datetime:     time.Now().UTC().Add(time.Minute * 2),
				Duration:     time.Hour,
				RemindBefore: time.Hour,
			},
		},
		{
			name: "2",
			in: app.EventDTO{
				UserID:       "",
				Title:        "test",
				Description:  storage.EventDescription{},
				Datetime:     time.Now().UTC().Add(time.Minute * 4),
				Duration:     time.Hour,
				RemindBefore: time.Hour,
			},
		},
		{
			name: "3",
			in: app.EventDTO{
				UserID:       "",
				Title:        "test",
				Description:  storage.EventDescription{},
				Datetime:     time.Now().UTC().Add(time.Minute * 6),
				Duration:     time.Hour,
				RemindBefore: time.Hour,
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		e.T().Run(tc.name, func(t *testing.T) {
			_, err := e.app.CreateEvent(e.ctx, tc.in)
			require.NoError(e.T(), err)
		})
	}
}

func (e *AppSuite) TestCreateEventWithTheSameDatesGetErr() {
	now := time.Now().UTC()
	dto := app.EventDTO{
		UserID:       "",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     now,
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	_, err := e.app.CreateEvent(e.ctx, dto)
	require.NoError(e.T(), err)
	dto = app.EventDTO{
		UserID:       "",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     now,
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	_, err = e.app.CreateEvent(e.ctx, dto)
	require.ErrorIs(e.T(), err, app.ErrDateBusy)
}

func (e *AppSuite) TestUpdateEventWithoutErr() {
	now := time.Now().UTC()
	createDTO := app.EventDTO{
		UserID:       "",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     now,
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	event, err := e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	testCases := []struct {
		name     string
		id       string
		userID   string
		datetime time.Time
	}{
		{name: "id1 userId1", id: event.ID(), userID: "1", datetime: now},
		{name: "id1 userId1", id: event.ID(), userID: "1", datetime: now.Add(time.Second)},
		{name: "id1 userId1", id: event.ID(), userID: "1", datetime: now.Add(time.Second * 2)},
		{name: "id1 userId1", id: event.ID(), userID: "1", datetime: now.Add(time.Second * 3)},
	}
	for _, tc := range testCases {
		tc := tc

		dto := app.UpdateEventDTO{
			ID:           tc.id,
			UserID:       tc.userID,
			Title:        "test",
			Description:  storage.EventDescription{},
			Datetime:     tc.datetime,
			Duration:     time.Hour,
			RemindBefore: time.Hour,
		}
		_, err = e.app.UpdateEvent(e.ctx, dto)
		require.NoError(e.T(), err)
	}
}

func (e *AppSuite) TestUpdateEventNotFoundGetErr() {
	dto := app.UpdateEventDTO{
		ID:           "1",
		UserID:       "1",
		Title:        "",
		Description:  storage.EventDescription{},
		Datetime:     time.Now().UTC(),
		Duration:     0,
		RemindBefore: 0,
	}
	_, err := e.app.UpdateEvent(e.ctx, dto)
	require.ErrorIs(e.T(), err, app.ErrNotFound)
}

func (e *AppSuite) TestUpdateEventWithBusyDateGetErr() {
	now := time.Now().UTC()
	createDTO := app.EventDTO{
		UserID:       "1",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     now,
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	event, err := e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	createDTO = app.EventDTO{
		UserID:       "1",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     now.Add(time.Hour),
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	_, err = e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	dto := app.UpdateEventDTO{
		ID:           event.ID(),
		UserID:       "1",
		Title:        "",
		Description:  storage.EventDescription{},
		Datetime:     now.Add(time.Hour),
		Duration:     0,
		RemindBefore: 0,
	}
	_, err = e.app.UpdateEvent(e.ctx, dto)
	require.ErrorIs(e.T(), err, app.ErrDateBusy)
}

func (e *AppSuite) TestDeleteEventWithoutErr() {
	now := time.Now().UTC()
	createDTO := app.EventDTO{
		UserID:       "1",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     now,
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	event, err := e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	err = e.app.DeleteEvent(e.ctx, event.ID())
	require.NoError(e.T(), err)
}

func (e *AppSuite) TestDeleteEventNotFountGetErr() {
	err := e.app.DeleteEvent(e.ctx, "")
	require.ErrorIs(e.T(), err, app.ErrNotFound)
}

func (e *AppSuite) TestEventsForPeriodNotFoundGetErr() {
	_, err := e.app.EventsForDay(e.ctx, time.Now().UTC())
	require.ErrorIs(e.T(), err, app.ErrNotFound)

	_, err = e.app.EventsForWeek(e.ctx, time.Now().UTC())
	require.ErrorIs(e.T(), err, app.ErrNotFound)

	_, err = e.app.EventsForMonth(e.ctx, time.Now().UTC())
	require.ErrorIs(e.T(), err, app.ErrNotFound)
}

func (e *AppSuite) TestEventsForDayWithoutErr() {
	specificDate := time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)
	createDTO := app.EventDTO{
		UserID:       "",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     specificDate,
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	_, err := e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	createDTO = app.EventDTO{
		UserID:       "",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     specificDate.AddDate(0, 0, 2),
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	_, err = e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	events := make(map[string]app.Event)
	for i := 1; i < 4; i++ {
		createDTO = app.EventDTO{
			UserID:       "",
			Title:        "test",
			Description:  storage.EventDescription{},
			Datetime:     specificDate.AddDate(0, 0, 1).Add(time.Hour * time.Duration(i)),
			Duration:     time.Hour,
			RemindBefore: time.Hour,
		}
		createdEvent, err := e.app.CreateEvent(e.ctx, createDTO)
		require.NoError(e.T(), err)
		events[createdEvent.ID()] = createdEvent
	}

	res, err := e.app.EventsForDay(e.ctx, specificDate.AddDate(0, 0, 1))
	require.NoError(e.T(), err)
	require.Equal(e.T(), events, res)
}

func (e *AppSuite) TestEventsForWeekWithoutErr() {
	specificDate := time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)
	createDTO := app.EventDTO{
		UserID:       "",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     specificDate.AddDate(0, 0, 2),
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	_, err := e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	createDTO = app.EventDTO{
		UserID:       "",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     specificDate.AddDate(0, 0, 10),
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	_, err = e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	events := make(map[string]app.Event)
	for i := 3; i < 10; i++ {
		createDTO = app.EventDTO{
			UserID:       "",
			Title:        "test",
			Description:  storage.EventDescription{},
			Datetime:     specificDate.AddDate(0, 0, i),
			Duration:     time.Hour,
			RemindBefore: time.Hour,
		}
		createdEvent, err := e.app.CreateEvent(e.ctx, createDTO)
		require.NoError(e.T(), err)
		events[createdEvent.ID()] = createdEvent
	}

	res, err := e.app.EventsForWeek(e.ctx, specificDate.AddDate(0, 0, 3))
	require.NoError(e.T(), err)
	require.Equal(e.T(), events, res)
}

func (e *AppSuite) TestEventsForMonthWithoutErr() {
	specificDate := time.Date(2021, 10, 0, 0, 0, 0, 0, time.UTC)
	createDTO := app.EventDTO{
		UserID:       "",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     specificDate,
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	_, err := e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	createDTO = app.EventDTO{
		UserID:       "",
		Title:        "test",
		Description:  storage.EventDescription{},
		Datetime:     specificDate.AddDate(0, 0, 32),
		Duration:     time.Hour,
		RemindBefore: time.Hour,
	}
	_, err = e.app.CreateEvent(e.ctx, createDTO)
	require.NoError(e.T(), err)

	events := make(map[string]app.Event)
	for i := 1; i < 32; i++ {
		createDTO = app.EventDTO{
			UserID:       "",
			Title:        "test",
			Description:  storage.EventDescription{},
			Datetime:     specificDate.AddDate(0, 0, i),
			Duration:     time.Hour,
			RemindBefore: time.Hour,
		}
		createdEvent, err := e.app.CreateEvent(e.ctx, createDTO)
		require.NoError(e.T(), err)
		events[createdEvent.ID()] = createdEvent
	}

	res, err := e.app.EventsForMonth(e.ctx, specificDate.AddDate(0, 0, 1))
	require.NoError(e.T(), err)
	require.Equal(e.T(), events, res)
}
