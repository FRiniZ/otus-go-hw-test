package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/memory"
	sqlstorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/sql"
	internalrmq "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/transport/rabbitmq"
)

type SchedulerConf struct {
	Logger   logger.Conf      `toml:"logger"`
	Storage  storage.Conf     `toml:"storage"`
	RabbitMQ internalrmq.Conf `toml:"rabbitmq"`
	Period   time.Duration    `toml:"period"`
}

type Scheduler struct {
	conf     SchedulerConf
	log      Logger
	storage  Storage
	producer internalrmq.Producer
}

func NewScheduler(conf SchedulerConf) *Scheduler {
	var db Storage

	log, err := logger.New(conf.Logger.Level, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't allocate logger:%v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch conf.Storage.DB {
	case "in-memory":
		db = memorystorage.New()
	case "sql":
		db = sqlstorage.New(conf.Storage.DSN)
	}
	err = db.Connect(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't connect to storage:%v\n", err) //nolint:gocritic
		os.Exit(1)
	}

	producer := internalrmq.NewProducer(log, conf.RabbitMQ)

	if err := producer.Connect(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Can't connect to RabbitMQ:%v\n", err)
		os.Exit(1)
	}

	return &Scheduler{
		conf:     conf,
		log:      log,
		storage:  db,
		producer: *producer,
	}
}

func (s Scheduler) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	d := time.Duration(s.conf.Period)
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// TODO Должны ли мы завершить работу если процесс возвращает ошибку?
			s.log.Debugf("Starting notification process...\n")
			if err := s.Process(ctx); err != nil {
				s.log.Errorf("%v", err)
				return
			}
		}
	}
}

func (s Scheduler) Process(ctx context.Context) error {
	sent := 0
	events, err := s.storage.ListEventsDayOfNotice(ctx, time.Now())
	if err != nil {
		return err
	}

	for _, e := range events {
		err := s.producer.SendNotification(ctx, &e)
		if err != nil {
			return fmt.Errorf("Process:%w", err)
		}

		err = s.storage.UpdateEventNotified(ctx, e.ID)
		if err != nil {
			return fmt.Errorf("Process:%w", err)
		}
		sent++
	}

	s.log.Debugf("Notifications sent:%v\n", sent)
	s.log.Debugf("Notification process has finished\n")

	return nil
}
