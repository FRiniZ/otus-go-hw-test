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
	storage  SchedulerStorage
	producer SchedulerProducer
}

type SchedulerStorage interface {
	Connect(context.Context) error
	Close(context.Context) error
	ListEventsDayOfNotice(context.Context, time.Time) ([]storage.Event, error)
	DeleteEventsOlderDate(context.Context, time.Time) (int64, error)
}

type SchedulerProducer interface {
	Connect(context.Context) error
	Close(context.Context) error
	SendNotification(context.Context, *storage.Event) error
}

func NewScheduler(conf SchedulerConf) *Scheduler {
	var db SchedulerStorage

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
		producer: producer,
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
			ctxStop, cancelStop := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelStop()
			s.Stop(ctxStop)
			return

		case <-ticker.C:
			// TODO Должны ли мы завершить работу если процесс возвращает ошибку?
			date := time.Now()
			s.log.Debugf("Starting notification process...\n")
			if sent, err := s.SendNotification(ctx, date); err != nil {
				s.log.Errorf("%v", err)
				return
			} else {
				s.log.Debugf("Notifications  sent:%v\n", sent)
			}

			s.log.Debugf("Starting to remove events that are older than a year\n")
			if deleted, err := s.DeleteEventsOlderDate(ctx, date.AddDate(-1, 0, 0)); err != nil {
				s.log.Errorf("%v", err)
				return
			} else {
				s.log.Debugf("Old events deleted:%v\n", deleted)
			}
			s.log.Debugf("Notification process has finished\n")
		}
	}
}

func (s Scheduler) Stop(ctx context.Context) {
	s.producer.Close(ctx)
	s.log.Debugf("Producer closed\n")
	s.storage.Close(ctx)
	s.log.Debugf("Storage closed\n")
}

func (s Scheduler) DeleteEventsOlderDate(ctx context.Context, date time.Time) (int64, error) {
	return s.storage.DeleteEventsOlderDate(ctx, date)
}

func (s Scheduler) SendNotification(ctx context.Context, date time.Time) (int64, error) {
	sent := int64(0)
	events, err := s.storage.ListEventsDayOfNotice(ctx, date)
	if err != nil {
		return sent, err
	}

	for _, e := range events {
		err := s.producer.SendNotification(ctx, &e)
		if err != nil {
			return sent, fmt.Errorf("SendNotification:%w", err)
		}

		//		err = s.storage.UpdateEventNotified(ctx, e.ID)
		//		if err != nil {
		//			return fmt.Errorf("Process:%w", err)
		//		}
		sent++
	}
	return sent, nil
}
