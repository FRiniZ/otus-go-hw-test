package memorystorage

import (
	"testing"
	"time"

	storage "github.com/FRiniZ/otus-go-hw-test/hw12_13_14_15_calendar/internal/storage"
	faker "github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Skip()
	db := New()
	num := 10000

	events := make([]storage.Event, num)
	for i := 0; i < num; i++ {
		faker.FakeData(&events[i])
	}

	t.Parallel()
	for i := 0; i < num; i++ {
		t.Run("Insert_And_Lookup", func(t *testing.T) {
			ev := events[i]
			err := db.InsertEvent(&ev)
			require.NotEqual(t, int64(0), ev.ID, "e.ID must be not zero")
			require.Equal(t, nil, err, "not equal result")
			ev2, err := db.LookupEvent(ev.ID)
			require.NoError(t, err, "err is not nil")
			require.Equal(t, ev.ID, ev2.ID, "not equal result")
			err = db.DeleteEvent(&ev)
			require.NoError(t, err, "err is not nil")
			ev2, err = db.LookupEvent(ev.ID)
			require.NoError(t, err, "err is not nil")
			require.Equal(t, int64(0), ev2.ID, "not equal result")
		})
	}
}

func TestStorageRules(t *testing.T) {
	db := New()

	t.Parallel()
	t.Run("Checking userID", func(t *testing.T) {
		var ev storage.Event
		faker.FakeData(&ev)
		ev.UserID = 0
		err := db.InsertEvent(&ev)
		require.ErrorIs(t, err, storage.ErrUserID, "expected err message")
	})

	t.Run("Checking userTitle", func(t *testing.T) {
		var ev storage.Event
		faker.FakeData(&ev)
		ev.Title = string(make([]byte, 200))
		err := db.InsertEvent(&ev)
		require.ErrorIs(t, err, storage.ErrTitle, "expected err message")
	})

	t.Run("Checking userOnTime", func(t *testing.T) {
		var ev storage.Event
		faker.FakeData(&ev)
		ev.OnTime = time.Time{}
		err := db.InsertEvent(&ev)
		require.ErrorIs(t, err, storage.ErrOnTime, "expected err message")
	})

	t.Run("Checking userOffTime", func(t *testing.T) {
		var ev storage.Event
		faker.FakeData(&ev)
		ev.OnTime = time.Now()
		ev.OffTime = time.Now().AddDate(0, 0, -1)
		err := db.InsertEvent(&ev)
		require.ErrorIs(t, err, storage.ErrOffTime, "expected err message")

		ev.OnTime = time.Now()
		ev.OffTime = ev.OnTime
		err = db.InsertEvent(&ev)
		require.ErrorIs(t, err, storage.ErrOffTime, "expected err message")
	})

	t.Run("Checking userNotifyTime", func(t *testing.T) {
		var ev storage.Event
		faker.FakeData(&ev)
		ev.OnTime = time.Now()
		ev.OffTime = time.Now().AddDate(0, 0, +7)
		ev.NotifyTime = time.Now().AddDate(0, 0, +8)
		err := db.InsertEvent(&ev)
		require.ErrorIs(t, err, storage.ErrNotifyTime, "expected err message")

		ev.OnTime = time.Now()
		ev.OffTime = time.Now().AddDate(0, 0, +7)
		ev.NotifyTime = time.Now().AddDate(0, 0, -1)
		err = db.InsertEvent(&ev)
		require.ErrorIs(t, err, storage.ErrNotifyTime, "expected err message")
	})
}
