package internalrmq

import (
	"context"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
)

type DummyConsumer struct {
	notifyChannel chan model.NotificationMsg
}

func NewDummyConsumer(notifyChannel chan model.NotificationMsg) *DummyConsumer {
	return &DummyConsumer{notifyChannel: notifyChannel}
}

func (c *DummyConsumer) Connect(ctx context.Context) error {
	return nil
}

func (c *DummyConsumer) NotifyChannel() <-chan model.NotificationMsg {
	return c.notifyChannel
}

func (c *DummyConsumer) Close(ctx context.Context) error {
	return nil
}
