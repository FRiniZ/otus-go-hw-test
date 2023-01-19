package app

import (
	"context"
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
		err := app.checkBasicRules(&event, true)
		require.ErrorIs(t, err, ErrID)

		err = app.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrUserID)

		event.UserID = 1
		err = app.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrTitle)

		event.Title = "TitleN1"
		err = app.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOnTime)

		event.OnTime = currTime
		err = app.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOffTime)

		event.OffTime = currTime
		err = app.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOffTime)

		event.OffTime = currTime.AddDate(0, 0, -1)
		err = app.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOffTime)

		event.OffTime = currTime.AddDate(0, 0, 7)
		err = app.checkBasicRules(&event, false)
		require.NoError(t, err)

		err = app.InsertEvent(ctx, &event)
		require.NoError(t, err)

		eventCopy := storage.Event{
			ID:          0,
			UserID:      event.UserID,
			Title:       event.Title,
			Description: event.Description,
			OnTime:      event.OnTime,
			OffTime:     event.OffTime,
			NotifyTime:  event.NotifyTime,
		}
		err = app.InsertEvent(ctx, &eventCopy)
		require.ErrorIs(t, err, memorystorage.ErrDataRangeIsBusy)
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

		event.OnTime = currTime.AddDate(0, 0, 2)
		err = app.UpdateEvent(ctx, &event)
		require.NoError(t, err)

		eventFound, err := app.LookupEvent(ctx, event.ID)
		require.NoError(t, err)
		require.EqualValues(t, userID, eventFound.UserID)

		events, err := app.ListEvents(ctx, event.UserID)
		require.NoError(t, err)
		require.EqualValues(t, int(1), len(events))

		err = app.DeleteEvent(ctx, event.ID)
		require.NoError(t, err)
	})
}
