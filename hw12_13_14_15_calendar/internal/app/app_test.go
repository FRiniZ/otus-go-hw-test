package app

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	ctx := context.Background()
	db := memorystorage.New()
	log, err := logger.New("DEBUG", os.Stdout)
	app := New(log, db)
	require.NoError(t, err)

	t.Run("test_rules", func(t *testing.T) {
		currTime := time.Now()
		event := storage.Event{
			Title:       string(make([]byte, 151)),
			Description: "DescriptionN1",
			OnTime:      time.Time{},
			OffTime:     time.Time{},
			NotifyTime:  time.Time{},
		}
		err := app.CheckBasicRules(&event, true)
		require.ErrorIs(t, err, ErrID)

		err = app.CheckBasicRules(&event, false)
		require.ErrorIs(t, err, ErrUserID)

		event.UserID = 1
		err = app.CheckBasicRules(&event, false)
		require.ErrorIs(t, err, ErrTitle)

		event.Title = "TitleN1"
		err = app.CheckBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOnTime)

		event.OnTime = currTime
		err = app.CheckBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOffTime)

		event.OffTime = currTime
		err = app.CheckBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOffTime)

		event.OffTime = currTime.AddDate(0, 0, -1)
		err = app.CheckBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOffTime)

		event.OffTime = currTime.AddDate(0, 0, 7)
		err = app.CheckBasicRules(&event, false)
		require.NoError(t, err)

		err = app.InsertEvent(ctx, &event)
		require.NoError(t, err)
		err = app.InsertEvent(ctx, &event)
		require.ErrorIs(t, err, ErrDateBusy)
	})

	t.Run("test_api", func(t *testing.T) {
		currTime := time.Now()
		userID := int64(100)
		event := storage.Event{
			UserID:      userID,
			Title:       "TitleN1",
			Description: "DescriptionN1",
			OnTime:      currTime.AddDate(0, 0, 1),
			OffTime:     currTime.AddDate(0, 0, 7),
			NotifyTime:  time.Time{},
		}
		err := app.InsertEvent(ctx, &event)
		require.NoError(t, err)
		fmt.Println(event.OnTime)

		event.OnTime = currTime.AddDate(0, 0, 2)
		err = app.UpdateEvent(ctx, &event)
		fmt.Println(event.OnTime)
		require.NoError(t, err)

		eventFound, err := app.LookupEvent(ctx, event.ID)
		require.NoError(t, err)
		require.EqualValues(t, userID, eventFound.UserID)

		events, err := app.ListEvents(ctx, event.UserID)
		require.NoError(t, err)
		require.EqualValues(t, int(1), len(events))

		err = app.DeleteEvent(ctx, &event)
		require.NoError(t, err)
	})
}
