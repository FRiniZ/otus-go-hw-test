package app

import (
	"context"
	"os"
	"testing"
	"time"

	logger "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/memory"
	internalrmq "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/transport/rabbitmq"
	"github.com/stretchr/testify/require"
)

func TestScheduler(t *testing.T) {
	ctx := context.Background()

	db := memorystorage.New()
	require.NotNil(t, db)
	log := logger.NewLogger("DEBUG", os.Stdout)
	producer := internalrmq.NewDummyProducer()
	scheduler := Scheduler{log: log, storage: db, producer: producer}

	t.Run("test_send_notification", func(t *testing.T) {
		currTime := time.Now()
		userID := int64(100)
		event := model.Event{
			UserID:      userID,
			Title:       "TitleN1",
			Description: "DescriptionN1",
			OnTime:      currTime.AddDate(0, 0, 1),
			OffTime:     currTime.AddDate(0, 0, 7),
			NotifyTime:  currTime.AddDate(0, 0, 0),
		}

		err := db.InsertEvent(ctx, &event)
		require.NoError(t, err)

		n, err := scheduler.SendNotification(ctx, currTime)
		require.NoError(t, err)
		require.EqualValues(t, int64(1), n)

		db.DeleteEvent(ctx, event.ID)
		require.NoError(t, err)
	})

	t.Run("test_delete_old_events", func(t *testing.T) {
		currTime := time.Now()
		userID := int64(200)
		event := model.Event{
			UserID:      userID,
			Title:       "TitleN1",
			Description: "DescriptionN1",
			OnTime:      currTime.AddDate(0, 0, 1),
			OffTime:     currTime.AddDate(0, 0, 7),
			NotifyTime:  currTime.AddDate(0, 0, 0),
		}

		err := db.InsertEvent(ctx, &event)
		require.NoError(t, err)

		eventID := event.ID
		require.NotEmpty(t, eventID)

		n, err := scheduler.DeleteEventsOlderDate(ctx, currTime.AddDate(1, 0, 1))
		require.NoError(t, err)
		require.EqualValues(t, int64(1), n)

		_, err = db.LookupEvent(ctx, eventID)
		require.ErrorIs(t, err, memorystorage.ErrEventNotFound)
	})
}
