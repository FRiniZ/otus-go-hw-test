package internalrmq

import (
	"context"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
)

type DummyProducer struct {
	log  Logger
	conf Conf
}

func NewDummyProducer(log Logger, conf Conf) *DummyProducer {
	return &DummyProducer{log: log, conf: conf}
}

func (c *DummyProducer) Connect(ctx context.Context) error {
	return nil
}

func (c *DummyProducer) Close(ctx context.Context) error {
	return nil
}

func (c *DummyProducer) SendNotification(ctx context.Context, event *storage.Event) error {
	return nil
}
