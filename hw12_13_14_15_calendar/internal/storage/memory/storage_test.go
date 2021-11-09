package memory_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageSuite struct {
	suite.Suite
	rep *memory.InMemory
	ctx context.Context
}

func (e *StorageSuite) SetupTest() {
	e.rep = memory.New()
	e.ctx = context.Background()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}

func (e *StorageSuite) TestInvalidInitializationGetErr() {
	rep := memory.InMemory{}

	_, err := rep.GetByID(context.Background(), "")
	require.ErrorIs(e.T(), err, storage.ErrInvalidInitialization)

	_, err = rep.GetByDatetime(context.Background(), time.Now())
	require.ErrorIs(e.T(), err, storage.ErrInvalidInitialization)

	err = rep.Create(context.Background(), nil)
	require.ErrorIs(e.T(), err, storage.ErrInvalidInitialization)

	err = rep.Delete(context.Background(), "")
	require.ErrorIs(e.T(), err, storage.ErrInvalidInitialization)

	err = rep.Update(context.Background(), "", nil)
	require.ErrorIs(e.T(), err, storage.ErrInvalidInitialization)

	_, err = rep.GetByDateRange(context.Background(), time.Now(), time.Now())
	require.ErrorIs(e.T(), err, storage.ErrInvalidInitialization)
}

func (e *StorageSuite) TestGetDoneFromContextGetErr() {
	e.T().Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := e.rep.GetByID(ctx, "")
		require.ErrorIs(e.T(), err, context.Canceled)

		_, err = e.rep.GetByDatetime(ctx, time.Now())
		require.ErrorIs(e.T(), err, context.Canceled)

		err = e.rep.Create(ctx, nil)
		require.ErrorIs(e.T(), err, context.Canceled)

		err = e.rep.Delete(ctx, "")
		require.ErrorIs(e.T(), err, context.Canceled)

		err = e.rep.Update(ctx, "", nil)
		require.ErrorIs(e.T(), err, context.Canceled)

		_, err = e.rep.GetByDateRange(ctx, time.Now(), time.Now())
		require.ErrorIs(e.T(), err, context.Canceled)
	})
	e.T().Run("deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*0)
		defer cancel()

		_, err := e.rep.GetByID(ctx, "")
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		_, err = e.rep.GetByDatetime(ctx, time.Now())
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		err = e.rep.Create(ctx, nil)
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		err = e.rep.Delete(ctx, "")
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		err = e.rep.Update(ctx, "", nil)
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)

		_, err = e.rep.GetByDateRange(ctx, time.Now(), time.Now())
		require.ErrorIs(e.T(), err, context.DeadlineExceeded)
	})
}

// GetByDatetime.
func (e *StorageSuite) TestGetByDatetimeWithoutErr() {
	now := time.Now().UTC()
	event := storage.NewEvent("1", "1", "", storage.EventDescription{}, now, time.Second, time.Second)
	err := e.rep.Create(e.ctx, event)
	require.NoError(e.T(), err)

	res, err := e.rep.GetByDatetime(e.ctx, now)
	require.NoError(e.T(), err)
	require.Equal(e.T(), event, res)
}

func (e *StorageSuite) TestGetByDatetimeEventNotFoundGetErr() {
	now := time.Now().UTC()
	_, err := e.rep.GetByDatetime(e.ctx, now)
	require.ErrorIs(e.T(), err, storage.ErrNotFound)
}

// GetByID.
func (e *StorageSuite) TestGetByIdWithoutErr() {
	event := storage.NewEvent("1", "1", "", storage.EventDescription{}, time.Now(), time.Second, time.Second)
	err := e.rep.Create(e.ctx, event)
	require.NoError(e.T(), err)

	res, err := e.rep.GetByID(e.ctx, event.ID())
	require.NoError(e.T(), err)
	require.Equal(e.T(), event, res)
}

func (e *StorageSuite) TestGetIdEventNotFoundGetErr() {
	_, err := e.rep.GetByID(e.ctx, "1")
	require.ErrorIs(e.T(), err, storage.ErrNotFound)
}

// GetByDateRange.
func (e *StorageSuite) TestGetByDateRangeWithoutErr() {
	specificDate := time.Date(2021, 10, 0, 0, 0, 0, 0, time.UTC)

	thirtyOneDays := make(map[string]app.Event)
	sevenDays := make(map[string]app.Event)
	fourteenDays := make(map[string]app.Event)
	oneDay := make(map[string]app.Event)

	for i := 1; i <= 31; i++ {
		event := storage.NewEvent(
			strconv.Itoa(i),
			"1",
			"",
			storage.EventDescription{},
			specificDate.AddDate(0, 0, i),
			time.Second,
			time.Second)
		err := e.rep.Create(e.ctx, event)
		require.NoError(e.T(), err)

		if i >= 1 && i <= 31 {
			thirtyOneDays[strconv.Itoa(i)] = event
		}
		if i >= 1 && i <= 7 {
			sevenDays[strconv.Itoa(i)] = event
		}
		if i >= 1 && i <= 14 {
			fourteenDays[strconv.Itoa(i)] = event
		}
		if i == 1 {
			oneDay[strconv.Itoa(i)] = event
		}
	}

	testCases := []struct {
		name  string
		start time.Time
		end   time.Time
		want  map[string]app.Event
	}{
		{
			name:  "all test range",
			start: specificDate.AddDate(0, 0, 1),
			end:   specificDate.AddDate(0, 0, 31),
			want:  thirtyOneDays,
		},
		{
			name:  "seven days",
			start: specificDate.AddDate(0, 0, 1),
			end:   specificDate.AddDate(0, 0, 7),
			want:  sevenDays,
		},
		{
			name:  "fourteen days",
			start: specificDate.AddDate(0, 0, 1),
			end:   specificDate.AddDate(0, 0, 14),
			want:  fourteenDays,
		},
		{
			name:  "one day",
			start: specificDate.AddDate(0, 0, 1),
			end:   specificDate.AddDate(0, 0, 1),
			want:  oneDay,
		},
	}

	for _, tc := range testCases {
		tc := tc
		e.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := e.rep.GetByDateRange(e.ctx, tc.start, tc.end)
			require.NoError(t, err)
			require.Equal(t, tc.want, res)
		})
	}
}

func (e *StorageSuite) TestGetByDateRangeStartDateLessThanEndDateGetErr() {
	start := time.Date(2021, 10, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 8, 0, 0, 0, 0, 0, time.UTC)

	_, err := e.rep.GetByDateRange(e.ctx, start, end)
	require.ErrorIs(e.T(), err, storage.ErrStartDateMoreThanEndDate)
}

func (e *StorageSuite) TestGetByDateRangeNotFoundGetErr() {
	start := time.Date(2021, 8, 0, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 10, 0, 0, 0, 0, 0, time.UTC)

	_, err := e.rep.GetByDateRange(e.ctx, start, end)
	require.ErrorIs(e.T(), err, storage.ErrNotFound)
}

// Create.
func (e *StorageSuite) TestCreateWithoutErr() {
	now := time.Now()
	testCases := []struct {
		name     string
		id       string
		userID   string
		datetime time.Time
	}{
		{name: "id1 userId1", id: "1", userID: "1", datetime: now.Add(time.Second)},
		{name: "id2 userId1", id: "2", userID: "1", datetime: now.Add(time.Second * 2)},
		{name: "id3 userId1", id: "3", userID: "1", datetime: now.Add(time.Second * 3)},
		{name: "id4 userId2", id: "4", userID: "2", datetime: now.Add(time.Second * 4)},
		{name: "id5 userId2", id: "5", userID: "2", datetime: now.Add(time.Second * 5)},
		{name: "id6 userId2", id: "6", userID: "2", datetime: now.Add(time.Second * 6)},
	}
	for _, tc := range testCases {
		tc := tc
		e.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()
			event := storage.NewEvent(
				tc.id,
				tc.userID,
				"",
				storage.EventDescription{},
				tc.datetime,
				time.Second,
				time.Second)
			err := e.rep.Create(e.ctx, event)
			require.NoError(t, err)
		})
	}
}

// Update.
func (e *StorageSuite) TestUpdateWithoutErr() {
	now := time.Now()
	event := storage.NewEvent("1", "1", "", storage.EventDescription{}, now, time.Second, time.Second)
	err := e.rep.Create(e.ctx, event)
	require.NoError(e.T(), err)

	testCases := []struct {
		name     string
		id       string
		userID   string
		datetime time.Time
	}{
		{name: "id1 userId1", id: event.ID(), userID: "1", datetime: now.Add(time.Second)},
		{name: "id1 userId1", id: event.ID(), userID: "1", datetime: now.Add(time.Second * 2)},
		{name: "id1 userId1", id: event.ID(), userID: "1", datetime: now.Add(time.Second * 3)},
	}

	for _, tc := range testCases {
		tc := tc
		e.T().Run(tc.name, func(t *testing.T) {
			event = storage.NewEvent(
				tc.id,
				tc.userID,
				"",
				storage.EventDescription{},
				tc.datetime,
				time.Second,
				time.Second)
			err = e.rep.Update(e.ctx, tc.id, event)
			require.NoError(t, err)
		})
	}
}

func (e *StorageSuite) TestUpdateEventNotFoundGetErr() {
	event := storage.NewEvent("1", "1", "", storage.EventDescription{}, time.Now(), time.Second, time.Second)
	err := e.rep.Update(e.ctx, event.ID(), event)
	require.ErrorIs(e.T(), err, storage.ErrNotFound)
}

// Delete.
func (e *StorageSuite) TestDeleteWithoutErr() {
	now := time.Now()
	testCases := []struct {
		name     string
		id       string
		userID   string
		datetime time.Time
	}{
		{name: "id1 userId1", id: "1", userID: "1", datetime: now.Add(time.Second)},
		{name: "id2 userId1", id: "2", userID: "1", datetime: now.Add(time.Second * 2)},
		{name: "id3 userId1", id: "3", userID: "1", datetime: now.Add(time.Second * 3)},
		{name: "id4 userId2", id: "4", userID: "2", datetime: now.Add(time.Second * 4)},
		{name: "id5 userId2", id: "5", userID: "2", datetime: now.Add(time.Second * 5)},
		{name: "id6 userId2", id: "6", userID: "2", datetime: now.Add(time.Second * 6)},
	}
	for _, tc := range testCases {
		tc := tc
		e.T().Run(tc.name, func(t *testing.T) {
			event := storage.NewEvent(
				tc.id,
				tc.userID,
				"",
				storage.EventDescription{},
				tc.datetime,
				time.Second,
				time.Second)
			err := e.rep.Create(e.ctx, event)
			require.NoError(t, err)
		})
	}
	for _, tc := range testCases {
		tc := tc
		e.T().Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := e.rep.Delete(e.ctx, tc.id)
			require.NoError(t, err)
		})
	}
}

func (e *StorageSuite) TestDeleteNotFoundEventsGetErr() {
	err := e.rep.Delete(e.ctx, "1")
	require.ErrorIs(e.T(), err, storage.ErrNotFound)
}
