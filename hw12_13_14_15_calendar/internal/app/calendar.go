package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/logger"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/model"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type CalendarConf struct {
	Logger     logger.Conf  `toml:"logger"`
	Storage    storage.Conf `toml:"storage"`
	HTTPServer struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	} `toml:"http-server"`
	GRPCServer struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
	} `toml:"grpc-server"`
}

type Calendar struct {
	conf    CalendarConf
	log     Logger
	storage CalendarStorage
}

type CalendarStorage interface {
	Connect(context.Context) error
	Close(context.Context) error
	InsertEvent(context.Context, *model.Event) error
	UpdateEvent(context.Context, *model.Event) error
	DeleteEvent(context.Context, int64) error
	LookupEvent(context.Context, int64) (model.Event, error)
	ListEvents(context.Context, int64) ([]model.Event, error)
	ListEventsRange(context.Context, int64, time.Time, time.Time) ([]model.Event, error)
	IsBusyDateTimeRange(context.Context, int64, int64, time.Time, time.Time) error
}

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

func (c *Calendar) checkBasicRules(e *model.Event, checkID bool) error {
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
	}

	return nil
}

func (c *Calendar) isBusyDateTimeRange(ctx context.Context, id, userID int64, onTime, offTime time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.storage.IsBusyDateTimeRange(ctx, id, userID, onTime, offTime)
}

func (c *Calendar) firstDayOfWeek(t time.Time) time.Time {
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, -1)
	}
	return t
}

func (c *Calendar) firstDayOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 0, -t.Day()+1)
}

func (c *Calendar) lastDayOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 1, -t.Day())
}

func (c *Calendar) Close(ctx context.Context) error {
	c.log.Infof("App closed\n")
	return c.storage.Close(ctx)
}

func (c *Calendar) InsertEvent(ctx context.Context, event *model.Event) error {
	if err := c.checkBasicRules(event, false); err != nil {
		return err
	}

	if err := c.isBusyDateTimeRange(ctx, event.ID, event.UserID, event.OnTime, event.OffTime); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return c.storage.InsertEvent(ctx, event)
}

func (c *Calendar) UpdateEvent(ctx context.Context, event *model.Event) error {
	if err := c.checkBasicRules(event, true); err != nil {
		return err
	}

	if err := c.isBusyDateTimeRange(ctx, event.ID, event.UserID, event.OnTime, event.OffTime); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.storage.UpdateEvent(ctx, event)
}

func (c *Calendar) DeleteEvent(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.storage.DeleteEvent(ctx, id)
}

func (c *Calendar) LookupEvent(ctx context.Context, id int64) (model.Event, error) {
	if id == 0 {
		return model.Event{}, ErrID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.storage.LookupEvent(ctx, id)
}

func (c *Calendar) ListEvents(ctx context.Context, userID int64) ([]model.Event, error) {
	if userID == 0 {
		return []model.Event{}, ErrUserID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.storage.ListEvents(ctx, userID)
}

func (c *Calendar) ListEventsDay(ctx context.Context, userID int64, date time.Time) ([]model.Event, error) {
	if userID == 0 {
		return []model.Event{}, ErrUserID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return c.storage.ListEventsRange(ctx, userID, date, date)
}

func (c *Calendar) ListEventsWeek(ctx context.Context, userID int64, date time.Time) ([]model.Event, error) {
	if userID == 0 {
		return []model.Event{}, ErrUserID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	monday := c.firstDayOfWeek(date)
	sunday := monday.AddDate(0, 0, 6)
	return c.storage.ListEventsRange(ctx, userID, monday, sunday)
}

func (c *Calendar) ListEventsMonth(ctx context.Context, userID int64, date time.Time) ([]model.Event, error) {
	if userID == 0 {
		return []model.Event{}, ErrUserID
	}
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	dayFirst := c.firstDayOfMonth(date)
	dayLast := c.lastDayOfMonth(date)
	return c.storage.ListEventsRange(ctx, userID, dayFirst, dayLast)
}

func NewCalendar(log Logger, conf CalendarConf, storage CalendarStorage) *Calendar {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := storage.Connect(ctx)
	if err != nil {
		exitfail(fmt.Sprintln("Can't connect to storage:%v", err))
	}

	return &Calendar{log: log, conf: conf, storage: storage}
}

func (c Calendar) Run(httpsrv Server, grpcsrv Server) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	g, ctxEG := errgroup.WithContext(ctx)

	func1 := func() error {
		return httpsrv.Start(ctxEG)
	}

	func2 := func() error {
		return grpcsrv.Start(ctxEG)
	}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := grpcsrv.Stop(ctx); err != nil {
			if !errors.Is(err, grpc.ErrServerStopped) {
				c.log.Errorf("failed to stop GRPC-server:%v\n", err)
			}
		}

		if err := httpsrv.Stop(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) &&
				!errors.Is(err, context.Canceled) {
				c.log.Errorf("failed to stop HTTP-server:%v\n", err)
			}
		}

		if err := c.storage.Close(ctx); err != nil {
			c.log.Errorf("failed to close db:%v\n", err)
		}
	}()

	g.Go(func1)
	g.Go(func2)

	if err := g.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) &&
			!errors.Is(err, grpc.ErrServerStopped) &&
			!errors.Is(err, context.Canceled) {
			c.log.Errorf("%v\n", err)
		}
	}
}
