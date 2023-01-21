package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

func helperEvent(ev *model.Event, i int) {
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

	events := make([]model.Event, num)
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
			err = db.DeleteEvent(context.Background(), ev.ID)
			require.NoError(t, err)
			ev2, err = db.LookupEvent(context.Background(), ev.ID)
			require.ErrorIs(t, err, ErrEventNotFound)
			require.Equal(t, int64(0), ev2.ID)
		})
	}
	t.Run("wrong_update", func(t *testing.T) {
		var ev model.Event
		helperEvent(&ev, 1)
		ev.ID = -1
		err := db.UpdateEvent(context.Background(), &ev)
		require.ErrorIs(t, err, ErrEventNotFound)
	})
}
