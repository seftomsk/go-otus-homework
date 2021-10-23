// +build postgres_test

package sqlstorage_test

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/app"
	"github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/seftomsk/go-otus-homework/hw12_13_14_15_calendar/internal/storage/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	user   = "admin"
	pwd    = "admin"
	host   = "localhost"
	port   = "5432"
	dbname = "mydb"
)

type StorageSuite struct {
	suite.Suite
	mu  sync.Mutex
	ctx context.Context
	rep *sqlstorage.InPostgres
	ids []string
}

func (e *StorageSuite) SetupSuite() {
	rep, err := sqlstorage.New(user, pwd, host, port, dbname)
	require.NoError(e.T(), err)
	e.rep = rep
}

func (e *StorageSuite) SetupTest() {
	e.ctx = context.Background()
}

func (e *StorageSuite) TearDownTest() {
	e.mu.Lock()
	for _, id := range e.ids {
		_ = e.rep.Delete(context.Background(), id)
	}
	e.mu.Unlock()
}

func (e *StorageSuite) TearDownSuite() {
	err := e.rep.Close()
	require.NoError(e.T(), err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}

func (e *StorageSuite) TestInvalidInitializationGetErr() {
	rep := sqlstorage.InPostgres{}

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
	e.mu.Lock()
	e.ids = append(e.ids, event.ID())
	e.mu.Unlock()

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
	now := time.Now().UTC()
	event := storage.NewEvent("1", "1", "", storage.EventDescription{}, now, time.Second, time.Second)
	err := e.rep.Create(e.ctx, event)
	require.NoError(e.T(), err)
	e.mu.Lock()
	e.ids = append(e.ids, event.ID())
	e.mu.Unlock()

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
		e.mu.Lock()
		e.ids = append(e.ids, event.ID())
		e.mu.Unlock()

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
func (e *StorageSuite) TestCreateWhenDescriptionIsNotNilWithoutErr() {
	now := time.Now().UTC()
	testCases := []struct {
		name     string
		id       string
		userID   string
		datetime time.Time
	}{
		{name: "id1 userId1", id: "1", userID: "1", datetime: now.Add(time.Second)},
	}
	for _, tc := range testCases {
		tc := tc
		e.T().Run(tc.name, func(t *testing.T) {
			event := storage.NewEvent(
				tc.id,
				tc.userID,
				"",
				storage.EventDescription{Data: "Desc", Valid: true},
				tc.datetime,
				time.Second,
				time.Second)
			err := e.rep.Create(e.ctx, event)
			require.NoError(t, err)
			e.mu.Lock()
			e.ids = append(e.ids, tc.id)
			e.mu.Unlock()
		})
	}
}

func (e *StorageSuite) TestCreateWithoutErr() {
	now := time.Now().UTC()
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
			e.mu.Lock()
			e.ids = append(e.ids, event.ID())
			e.mu.Unlock()
		})
	}
}

// Update.
func (e *StorageSuite) TestUpdateWithoutErr() {
	now := time.Now().UTC()
	event := storage.NewEvent("1", "1", "", storage.EventDescription{}, now, time.Second, time.Second)
	err := e.rep.Create(e.ctx, event)
	require.NoError(e.T(), err)
	e.mu.Lock()
	e.ids = append(e.ids, event.ID())
	e.mu.Unlock()

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
	now := time.Now().UTC()
	event := storage.NewEvent("1", "1", "", storage.EventDescription{}, now, time.Second, time.Second)
	err := e.rep.Update(e.ctx, event.ID(), event)
	require.ErrorIs(e.T(), err, storage.ErrNotFound)
}

// Delete.
func (e *StorageSuite) TestDeleteWithoutErr() {
	now := time.Now().UTC()
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
			err := e.rep.Delete(e.ctx, tc.id)
			require.NoError(t, err)
		})
	}
}

func (e *StorageSuite) TestDeleteNotFoundEventsGetErr() {
	err := e.rep.Delete(e.ctx, "1")
	require.ErrorIs(e.T(), err, storage.ErrNotFound)
}
