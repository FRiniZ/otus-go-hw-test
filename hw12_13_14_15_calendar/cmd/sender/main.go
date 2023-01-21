package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	internalrmq "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/transport/rabbitmq"
)

func main() {
	conf := NewConfig().SenderConf
	storage := storage.NewStorage(conf.Storage)
	logger := logger.NewLogger(conf.Logger.Level, os.Stdout)
	consumer := internalrmq.NewConsumer(logger, conf.URLRMQ)
	sender := app.NewSender(logger, conf, storage, consumer)

	sender.Run()
	filename := filepath.Base(os.Args[0])
	fmt.Printf("%s stopped\n", filename)
}
