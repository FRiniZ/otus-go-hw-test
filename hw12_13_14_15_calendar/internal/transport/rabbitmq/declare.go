package internalrmq

import "github.com/rabbitmq/amqp091-go"

func getQueueDeclated() (string, bool, bool, bool, bool, amqp091.Table) {
	return "notification", false, false, false, false, nil
}
