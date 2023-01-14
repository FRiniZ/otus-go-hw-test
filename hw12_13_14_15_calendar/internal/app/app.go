package app

import (
	"context"
	"errors"

	"github.com/FRiniZ/otus-go-hw-test/hw12_13_14_15_calendar/internal/storage"
)

var ErrDateBusy = errors.New("date is busy")

type App struct {
	log     Logger
	storage Storage
}

type Logger interface {
	Panicf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warningf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
}

type Storage interface {
	InsertEvent(*storage.Event) error
	UpdateEvent(*storage.Event) error
	DeleteEvent(*storage.Event) error
	LookupEvent(eID int64) (storage.Event, error)
	ListEvents(userID int64) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{log: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
