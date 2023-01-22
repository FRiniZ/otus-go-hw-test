package internalrmq

import amqp091 "github.com/rabbitmq/amqp091-go"

type Logger interface {
	Fatalf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warningf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
}

func getQueueDeclated() (string, bool, bool, bool, bool, amqp091.Table) {
	return "notification", false, false, false, false, nil
}
