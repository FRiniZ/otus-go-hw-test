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
	DeleteEvent(context.Context, int64) error
	LookupEvent(context.Context, int64) (storage.Event, error)
	ListEvents(context.Context, int64) ([]storage.Event, error)
	ListEventsRange(context.Context, int64, time.Time, time.Time) ([]storage.Event, error)
	IsBusyDateTimeRange(context.Context, int64, int64, time.Time, time.Time) error
}

func New(logger Logger, storage Storage) *App {
	return &App{log: logger, storage: storage}
}

func (a *App) checkBasicRules(e *storage.Event, checkID bool) error {
	if checkID && e.ID == 0 {
		return fmt.Errorf("%w: zero", ErrID)
	}

	if e.UserID == 0 {
		return fmt.Errorf("%w: zero", ErrUserID)
	}

	if len(e.Title) > 150 {
		return fmt.Errorf("%w: must be <=150", ErrTitle)
	}

	if e.OnTime.IsZero() {
		return fmt.Errorf("%w: empty", ErrOnTime)
	}

	switch {
	case e.OffTime.IsZero():
		return fmt.Errorf("%w: empty", ErrOffTime)
	case e.OffTime.Before(e.OnTime):
		return fmt.Errorf("%w: before OnTime", ErrOffTime)
	case e.OffTime.Equal(e.OnTime):
		return fmt.Errorf("%w: equal OnTime", ErrOffTime)
	}

	if !e.NotifyTime.IsZero() {
		if e.NotifyTime.After(e.OffTime) {
			return fmt.Errorf("%w: after OffTime", ErrNotifyTime)
		}
		if e.NotifyTime.Before(e.OnTime) {
			return fmt.Errorf("%w: before OnTime", ErrNotifyTime)
		}
	}

	return nil
}

func (a *App) isBusyDateTimeRange(ctx context.Context, id, userID int64, onTime, offTime time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.IsBusyDateTimeRange(ctx, id, userID, onTime, offTime)
}

func (a *App) firstDayOfWeek(t time.Time) time.Time {
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, -1)
	}
	return t
}

func (a *App) firstDayOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 0, -t.Day()+1)
}

func (a *App) lastDayOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 1, -t.Day())
}

func (a *App) Close(ctx context.Context) error {
	a.log.Infof("App closed\n")
	return a.storage.Close(ctx)
}

func (a *App) InsertEvent(ctx context.Context, event *storage.Event) error {
	if err := a.checkBasicRules(event, false); err != nil {
		return err
	}

	if err := a.isBusyDateTimeRange(ctx, event.ID, event.UserID, event.OnTime, event.OffTime); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.InsertEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event *storage.Event) error {
	if err := a.checkBasicRules(event, true); err != nil {
		return err
	}

	if err := a.isBusyDateTimeRange(ctx, event.ID, event.UserID, event.OnTime, event.OffTime); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.DeleteEvent(ctx, id)
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

func (a *App) ListEventsDay(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	if userID == 0 {
		return []storage.Event{}, ErrUserID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return a.storage.ListEventsRange(ctx, userID, date, date)
}

func (a *App) ListEventsWeek(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	if userID == 0 {
		return []storage.Event{}, ErrUserID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	monday := a.firstDayOfWeek(date)
	sunday := monday.AddDate(0, 0, 6)
	return a.storage.ListEventsRange(ctx, userID, monday, sunday)
}

func (a *App) ListEventsMonth(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	if userID == 0 {
		return []storage.Event{}, ErrUserID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	dayFirst := a.firstDayOfMonth(date)
	dayLast := a.lastDayOfMonth(date)
	return a.storage.ListEventsRange(ctx, userID, dayFirst, dayLast)
}
