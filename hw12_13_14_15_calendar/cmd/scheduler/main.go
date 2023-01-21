package main

import (
	"fmt"
	"os"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	internalrmq "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/transport/rabbitmq"
)

func main() {
	conf := NewConfig().SchedulerConf
	storage := storage.NewStorage(conf.Storage)
	logger := logger.NewLogger(conf.Logger.Level, os.Stdout)
	producer := internalrmq.NewProducer(logger, conf.URLRMQ)
	scheduler := app.NewScheduler(logger, conf, storage, producer)

	scheduler.Run()
	fmt.Println("calendar_scheduler stopped")
}
