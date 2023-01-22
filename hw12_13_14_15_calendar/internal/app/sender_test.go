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

func TestSender(t *testing.T) {
	db := memorystorage.New()
	require.NotNil(t, db)
	log := logger.NewLogger("DEBUG", os.Stdout)
	notifyChannel := make(chan model.NotificationMsg, 1)
	consumer := internalrmq.NewDummyConsumer(notifyChannel)
	sender := Sender{log: log, storage: db, consumer: consumer}

	t.Run("test_receive_notification", func(t *testing.T) {
		currTime := time.Now()
		eventID := int64(100)
		userID := int64(200)

		notify := model.NotificationMsg{
			ID:     eventID,
			Title:  "TitleN1",
			Date:   currTime,
			UserID: userID,
		}

		notifyChannel <- notify

		msg := <-sender.consumer.NotifyChannel()
		require.EqualValues(t, msg.ID, notify.ID)
	})

	t.Run("test_update_events", func(t *testing.T) {
		ctx := context.Background()
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
		require.EqualValues(t, false, event.Notified)

		err = sender.storage.UpdateEventNotified(ctx, event.ID)
		require.NoError(t, err)

		event, err = db.LookupEvent(ctx, eventID)
		require.NoError(t, err)
		require.EqualValues(t, true, event.Notified)
	})
}
