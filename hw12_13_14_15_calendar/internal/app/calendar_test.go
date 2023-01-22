package app

import (
	"context"
	"os"
	"testing"
	"time"

	logger "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestCalendar(t *testing.T) {
	ctx := context.Background()

	db := memorystorage.New()
	require.NotNil(t, db)
	log := logger.NewLogger("DEBUG", os.Stdout)
	calendar := Calendar{log: log, storage: db}

	t.Run("test_rules", func(t *testing.T) {
		currTime := time.Now()
		event := model.Event{
			Title:       string(make([]byte, 151)),
			Description: "DescriptionN1",
			OnTime:      time.Time{},
			OffTime:     time.Time{},
			NotifyTime:  time.Time{},
		}
		err := calendar.checkBasicRules(&event, true)
		require.ErrorIs(t, err, ErrID)

		err = calendar.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrUserID)

		event.UserID = 1
		err = calendar.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrTitle)

		event.Title = "TitleN1"
		err = calendar.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOnTime)

		event.OnTime = currTime
		err = calendar.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOffTime)

		event.OffTime = currTime
		err = calendar.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOffTime)

		event.OffTime = currTime.AddDate(0, 0, -1)
		err = calendar.checkBasicRules(&event, false)
		require.ErrorIs(t, err, ErrOffTime)

		event.OffTime = currTime.AddDate(0, 0, 7)
		err = calendar.checkBasicRules(&event, false)
		require.NoError(t, err)

		err = calendar.InsertEvent(ctx, &event)
		require.NoError(t, err)

		eventCopy := model.Event{
			ID:          0,
			UserID:      event.UserID,
			Title:       event.Title,
			Description: event.Description,
			OnTime:      event.OnTime,
			OffTime:     event.OffTime,
			NotifyTime:  event.NotifyTime,
		}
		err = calendar.InsertEvent(ctx, &eventCopy)
		require.ErrorIs(t, err, memorystorage.ErrDataRangeIsBusy)
	})

	t.Run("test_api", func(t *testing.T) {
		currTime := time.Now()
		userID := int64(100)
		event := model.Event{
			UserID:      userID,
			Title:       "TitleN1",
			Description: "DescriptionN1",
			OnTime:      currTime.AddDate(0, 0, 1),
			OffTime:     currTime.AddDate(0, 0, 7),
			NotifyTime:  time.Time{},
		}
		err := calendar.InsertEvent(ctx, &event)
		require.NoError(t, err)

		event.OnTime = currTime.AddDate(0, 0, 2)
		err = calendar.UpdateEvent(ctx, &event)
		require.NoError(t, err)

		eventFound, err := calendar.LookupEvent(ctx, event.ID)
		require.NoError(t, err)
		require.EqualValues(t, userID, eventFound.UserID)

		events, err := calendar.ListEvents(ctx, event.UserID)
		require.NoError(t, err)
		require.EqualValues(t, int(1), len(events))

		err = calendar.DeleteEvent(ctx, event.ID)
		require.NoError(t, err)
	})
}
