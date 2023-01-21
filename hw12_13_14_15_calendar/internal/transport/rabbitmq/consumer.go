package internalrmq

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	log             Logger
	url             string
	channel         *amqp091.Channel
	connect         *amqp091.Connection
	queue           amqp091.Queue
	deliveryChannel <-chan amqp091.Delivery
	notifyChannel   chan model.NotificationMsg
}

var ErrCantRecvMsg = errors.New("can't receive message")

func NewConsumer(log Logger, url string) *Consumer {
	notifyChannel := make(chan model.NotificationMsg, 1)
	return &Consumer{log: log, url: url, notifyChannel: notifyChannel}
}

func (c *Consumer) Connect(ctx context.Context) error {
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

	c.deliveryChannel, err = c.channel.Consume(c.queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range c.deliveryChannel {
			notify, err := c.unpackMsg(msg)
			if err != nil {
				c.log.Errorf("RecvNotification:%v\n", err)
			} else {
				c.notifyChannel <- notify
				c.log.Debugf("Received notification: %v\n", notify)
			}
		}
	}()

	c.log.Debugf("Connected to RabbitMQ\n")

	return nil
}

func (c *Consumer) NotifyChannel() <-chan model.NotificationMsg {
	return c.notifyChannel
}

func (c *Consumer) Close(ctx context.Context) error {
	c.connect.Close()
	c.channel.Close()
	return nil
}

func (c *Consumer) unpackMsg(msg amqp091.Delivery) (model.NotificationMsg, error) {
	var notify model.NotificationMsg

	if err := json.Unmarshal(msg.Body, &notify); err != nil {
		return notify, err
	}
	return notify, nil
}
