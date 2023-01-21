package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
	memorystorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/memory"
	sqlstorage "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage/sql"
)

type Conf struct {
	DB  string `toml:"db"`
	DSN string `toml:"dsn"`
}

type Storage interface {
	Connect(context.Context) error
	Close(context.Context) error
	InsertEvent(context.Context, *model.Event) error
	UpdateEvent(context.Context, *model.Event) error
	DeleteEvent(context.Context, int64) error
	LookupEvent(context.Context, int64) (model.Event, error)
	ListEvents(context.Context, int64) ([]model.Event, error)
	ListEventsRange(context.Context, int64, time.Time, time.Time) ([]model.Event, error)
	IsBusyDateTimeRange(context.Context, int64, int64, time.Time, time.Time) error

	// for producers
	ListEventsDayOfNotice(context.Context, time.Time) ([]model.Event, error)
	DeleteEventsOlderDate(context.Context, time.Time) (int64, error)

	// for consumers
	UpdateEventNotified(context.Context, int64) error
}

func NewStorage(conf Conf) Storage {
	switch conf.DB {
	case "in-memory":
		return memorystorage.New()
	case "sql":
		return sqlstorage.New(conf.DSN)
	}

	fmt.Fprintln(os.Stderr, "wrong DB")
	os.Exit(1)
	return nil
}
