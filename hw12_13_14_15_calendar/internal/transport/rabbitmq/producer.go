package internalrmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	log     Logger
	url     string
	channel *amqp091.Channel
	connect *amqp091.Connection
	queue   amqp091.Queue
}

var ErrCantSendMsg = errors.New("can't send message")

func NewProducer(log Logger, url string) *Producer {
	return &Producer{log: log, url: url}
}

func (c *Producer) Connect(ctx context.Context) error {
	var err error
	c.log.Debugf("Connecting to RabbitMQ...\n")
	c.connect, err = amqp091.Dial(c.url)
	if err != nil {
		return err
	}

	c.channel, err = c.connect.Channel()
	if err != nil {
		return err
	}

	c.queue, err = c.channel.QueueDeclare(getQueueDeclated())
	if err != nil {
		return err
	}

	c.log.Debugf("Connected to RabbitMQ\n")

	return nil
}

func (c *Producer) Close(ctx context.Context) error {
	c.connect.Close()
	c.channel.Close()
	return nil
}

func (c *Producer) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (c *Producer) Stop(ctx context.Context) error {
	return nil
}

func (c *Producer) SendNotification(ctx context.Context, event *model.Event) error {
	msg := model.NotificationMsg{
		ID:     event.ID,
		Title:  event.Title,
		Date:   event.OnTime,
		UserID: event.UserID,
	}

	jdata, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	pub := amqp091.Publishing{
		ContentType: "application/json",
		Body:        jdata,
	}

	if err := c.channel.PublishWithContext(ctx, "", c.queue.Name, false, false, pub); err != nil {
		return fmt.Errorf("SendNotification: %w", err)
	}
	c.log.Debugf("sent notification msg\n")
	return nil
}
