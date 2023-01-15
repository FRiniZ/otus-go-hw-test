package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
	storage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func helperEvent(ev *storage.Event, i int) {
	ev.UserID = int64(i + 1)
	ev.Title = fmt.Sprintf("Title_N%v", i+1)
	ev.Description = fmt.Sprintf("Description_N%v", i+1)
	ev.OnTime = time.Now()
	ev.OffTime = time.Now().AddDate(0, 0, 7)
	ev.NotifyTime = time.Now().AddDate(0, 0, 6)
}

func TestStorage(t *testing.T) {
	db := New()
	num := 10000

	events := make([]storage.Event, num)
	for i := 0; i < num; i++ {
		helperEvent(&events[i], i)
	}

	for i := 0; i < num; i++ {
		t.Run("insert_lookup_update_delete_parallel", func(t *testing.T) {
			i := i
			t.Parallel()
			ev := events[i]
			err := db.InsertEvent(context.Background(), &ev)
			require.Equal(t, nil, err)
			require.NotEqual(t, int64(0), ev.ID)
			require.NoError(t, err)
			ev2, err := db.LookupEvent(context.Background(), ev.ID)
			require.NoError(t, err)
			require.Equal(t, ev.ID, ev2.ID)
			err = db.DeleteEvent(context.Background(), &ev)
			require.NoError(t, err)
			ev2, err = db.LookupEvent(context.Background(), ev.ID)
			require.NoError(t, err)
			require.Equal(t, int64(0), ev2.ID)
		})
	}
	t.Run("wrong_update", func(t *testing.T) {
		var ev storage.Event
		helperEvent(&ev, 1)
		ev.ID = -1
		err := db.UpdateEvent(context.Background(), &ev)
		require.ErrorIs(t, err, app.ErrEventNotFound)
	})
}

func TestStorageRules(t *testing.T) {
	db := New()

	t.Parallel()
	t.Run("Checking userID", func(t *testing.T) {
		var ev storage.Event
		helperEvent(&ev, 1)
		ev.UserID = 0
		err := db.InsertEvent(context.Background(), &ev)
		require.ErrorIs(t, err, app.ErrUserID, "expected err message")
	})

	t.Run("Checking userTitle", func(t *testing.T) {
		var ev storage.Event
		helperEvent(&ev, 1)
		ev.Title = string(make([]byte, 200))
		err := db.InsertEvent(context.Background(), &ev)
		require.ErrorIs(t, err, app.ErrTitle, "expected err message")
	})

	t.Run("Checking userOnTime", func(t *testing.T) {
		var ev storage.Event
		helperEvent(&ev, 1)
		ev.OnTime = time.Time{}
		err := db.InsertEvent(context.Background(), &ev)
		require.ErrorIs(t, err, app.ErrOnTime, "expected err message")
	})

	t.Run("Checking userOffTime", func(t *testing.T) {
		var ev storage.Event
		helperEvent(&ev, 1)
		ev.OnTime = time.Now()
		ev.OffTime = time.Now().AddDate(0, 0, -1)
		err := db.InsertEvent(context.Background(), &ev)
		require.ErrorIs(t, err, app.ErrOffTime, "expected err message")

		ev.OnTime = time.Now()
		ev.OffTime = ev.OnTime
		err = db.InsertEvent(context.Background(), &ev)
		require.ErrorIs(t, err, app.ErrOffTime, "expected err message")
	})

	t.Run("Checking userNotifyTime", func(t *testing.T) {
		var ev storage.Event
		helperEvent(&ev, 1)
		ev.OnTime = time.Now()
		ev.OffTime = time.Now().AddDate(0, 0, +7)
		ev.NotifyTime = time.Now().AddDate(0, 0, +8)
		err := db.InsertEvent(context.Background(), &ev)
		require.ErrorIs(t, err, app.ErrNotifyTime, "expected err message")

		ev.OnTime = time.Now()
		ev.OffTime = time.Now().AddDate(0, 0, +7)
		ev.NotifyTime = time.Now().AddDate(0, 0, -1)
		err = db.InsertEvent(context.Background(), &ev)
		require.ErrorIs(t, err, app.ErrNotifyTime, "expected err message")
	})
}
