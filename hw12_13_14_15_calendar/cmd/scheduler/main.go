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
	config := NewConfig()
	storage := storage.NewStorage(config.Storage)
	logger := logger.NewLogger(config.Logger.Level, os.Stdout)
	producer := internalrmq.NewProducer(logger, config.UrlRMQ)
	scheduler := app.NewScheduler(logger, config.SchedulerConf, storage, producer)

	scheduler.Run()
	fmt.Println("calendar_scheduler stopped")
}
