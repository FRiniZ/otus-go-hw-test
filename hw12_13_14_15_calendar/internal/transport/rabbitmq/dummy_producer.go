package internalrmq

import (
	"context"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
)

type DummyProducer struct {
}

func NewDummyProducer() *DummyProducer {
	return &DummyProducer{}
}

func (c *DummyProducer) Connect(ctx context.Context) error {
	return nil
}

func (c *DummyProducer) Close(ctx context.Context) error {
	return nil
}

func (c *DummyProducer) SendNotification(ctx context.Context, event *model.Event) error {
	return nil
}
