package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
)

var (
	ErrID             = errors.New("wrong ID")
	ErrUserID         = errors.New("wrong UserID")
	ErrTitle          = errors.New("wrong Title")
	ErrDescription    = errors.New("wrong Description")
	ErrOnTime         = errors.New("wrong OnTime")
	ErrOffTime        = errors.New("wrong OffTime")
	ErrNotifyTime     = errors.New("wrong NotifyTime")
	ErrEventNotFound  = errors.New("event not found")
	ErrDateBusy       = errors.New("date is busy")
	ErrTooLongCloseDB = errors.New("too long close db")
)

type App struct {
	log     Logger
	storage Storage
}

type Logger interface {
	Fatalf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warningf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
}

type Storage interface {
	Connect(context.Context) error
	Close(context.Context) error
	InsertEvent(context.Context, *storage.Event) error
	UpdateEvent(context.Context, *storage.Event) error
	DeleteEvent(context.Context, *storage.Event) error
	LookupEvent(context.Context, int64) (storage.Event, error)
	ListEvents(context.Context, int64) ([]storage.Event, error)
	EmptyDate(context.Context, int64, time.Time) (bool, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{log: logger, storage: storage}
}

func (a *App) CheckBasicRules(e *storage.Event, checkID bool) error {
	if checkID && e.ID == 0 {
		return fmt.Errorf("%w(ID is zero)", ErrID)
	}

	if e.UserID == 0 {
		return fmt.Errorf("%w(UserID is %v)", ErrUserID, e.UserID)
	}

	if len(e.Title) > 150 {
		return fmt.Errorf("%w(len %v, must be <=150)", ErrTitle, len(e.Title))
	}

	if e.OnTime.IsZero() {
		return fmt.Errorf("%w(empty OnTime)", ErrOnTime)
	}

	if !e.OffTime.IsZero() {
		if e.OffTime.Before(e.OnTime) {
			return fmt.Errorf("%w(OffTime before OnTime)", ErrOffTime)
		}
		if e.OffTime.Equal(e.OnTime) {
			return fmt.Errorf("%w(OffTime equal OnTime)", ErrOffTime)
		}
	}

	if !e.NotifyTime.IsZero() {
		fmt.Printf("\n\nNotifyTime:%v\n\n", e.NotifyTime)
		if e.NotifyTime.After(e.OffTime) {
			return fmt.Errorf("%w(NotifyTime after OffTime)", ErrNotifyTime)
		}
		if e.NotifyTime.Before(e.OnTime) {
			return fmt.Errorf("%w(NotifyTime before OnTime)", ErrNotifyTime)
		}
	}

	// TODO Add checking ErrDateBusy

	return nil
}

func (a *App) EmptyDate(ctx context.Context, userID int64, date time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.EmptyDate(ctx, userID, date)
}

func (a *App) Close(ctx context.Context) error {
	a.log.Infof("App closed\n")
	return a.storage.Close(ctx)
}

func (a *App) InsertEvent(ctx context.Context, event *storage.Event) error {
	if err := a.CheckBasicRules(event, false); err != nil {
		return err
	}

	if empty, err := a.EmptyDate(ctx, event.UserID, event.OnTime); err != nil {
		return err
	} else if !empty {
		return ErrDateBusy
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.InsertEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	if err := a.CheckBasicRules(event, true); err != nil {
		return err
	}

	if empty, err := a.EmptyDate(ctx, event.UserID, event.OnTime); err != nil {
		return err
	} else if !empty {
		return ErrDateBusy
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, event *storage.Event) error {
	if err := a.CheckBasicRules(event, true); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.DeleteEvent(ctx, event)
}

func (a *App) LookupEvent(ctx context.Context, id int64) (storage.Event, error) {
	if id == 0 {
		return storage.Event{}, ErrID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.LookupEvent(ctx, id)
}

func (a *App) ListEvents(ctx context.Context, userID int64) ([]storage.Event, error) {
	if userID == 0 {
		return []storage.Event{}, ErrUserID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.ListEvents(ctx, userID)
}
