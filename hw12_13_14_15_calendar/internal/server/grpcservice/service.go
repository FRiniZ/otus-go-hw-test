package grpcservice

import (
	context "context"
	"net"
	"time"

	api "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/api/stub"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ctxKeyID int

const (
	KeyMethodID ctxKeyID = iota
)

type Conf struct {
	Port string `toml:"port"`
	Host string `toml:"host"`
}

type Logger interface {
	Fatalf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warningf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
}

type Application interface {
	InsertEvent(context.Context, *storage.Event) error
	UpdateEvent(context.Context, *storage.Event) error
	DeleteEvent(context.Context, int64) error
	LookupEvent(context.Context, int64) (storage.Event, error)
	ListEvents(context.Context, int64) ([]storage.Event, error)
	ListEventsDay(context.Context, int64, time.Time) ([]storage.Event, error)
	ListEventsWeek(context.Context, int64, time.Time) ([]storage.Event, error)
	ListEventsMonth(context.Context, int64, time.Time) ([]storage.Event, error)
}

type Service struct {
	log     Logger
	conf    Conf
	app     Application
	basesrv *grpc.Server
	api.UnimplementedCalendarServer
}

func (Service) APIEventFromEvent(event *storage.Event) *api.Event {
	return &api.Event{
		ID:          &event.ID,
		UserID:      &event.UserID,
		Title:       &event.Title,
		Description: &event.Description,
		OnTime:      timestamppb.New(event.OnTime),
		OffTime:     timestamppb.New(event.OffTime),
		NotifyTime:  timestamppb.New(event.NotifyTime),
	}
}

func (Service) EventFromAPIEvent(apiEvent *api.Event) *storage.Event {
	event := storage.Event{}

	event.ID = *apiEvent.ID
	event.UserID = *apiEvent.UserID
	event.Title = *apiEvent.Title
	event.Description = *apiEvent.Description
	if err := apiEvent.OnTime.CheckValid(); err == nil {
		event.OnTime = apiEvent.OnTime.AsTime().Local()
	}
	if err := apiEvent.OffTime.CheckValid(); err == nil {
		event.OffTime = apiEvent.OffTime.AsTime().Local()
	}
	if err := apiEvent.NotifyTime.CheckValid(); err == nil {
		event.NotifyTime = apiEvent.NotifyTime.AsTime().Local()
	}

	return &event
}

func (s Service) InsertEvent(ctx context.Context, req *api.ReqByEvent) (*api.RepID, error) {
	event := s.EventFromAPIEvent(req.Event)
	if err := s.app.InsertEvent(ctx, event); err != nil {
		return nil, err
	}

	return &api.RepID{ID: &event.ID}, nil
}

func (s Service) UpdateEvent(ctx context.Context, req *api.ReqByEvent) (*emptypb.Empty, error) {
	event := s.EventFromAPIEvent(req.Event)
	if err := s.app.UpdateEvent(ctx, event); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s Service) DeleteEvent(ctx context.Context, req *api.ReqByID) (*emptypb.Empty, error) {
	if err := s.app.DeleteEvent(ctx, *req.ID); err != nil {
		return nil, err
	}
	return new(emptypb.Empty), nil
}

func (s Service) LookupEvent(ctx context.Context, req *api.ReqByID) (*api.RepEvents, error) {
	event, err := s.app.LookupEvent(ctx, *req.ID)
	_ = event // to avoid lint err: event declared but not used (typecheck)
	if err != nil {
		return nil, err
	}

	rep := api.RepEvents{}
	rep.Event = append(rep.Event, s.APIEventFromEvent(&event))
	return &rep, nil
}

func (s Service) ListEvents(ctx context.Context, req *api.ReqByUser) (*api.RepEvents, error) {
	events, err := s.app.ListEvents(ctx, *req.UserID)
	if err != nil {
		return nil, err
	}

	rep := api.RepEvents{}
	rep.Event = make([]*api.Event, len(events))
	for i, event := range events {
		event := event
		rep.Event[i] = s.APIEventFromEvent(&event)
	}
	return &rep, nil
}

func (s Service) ListEventsDay(ctx context.Context, req *api.ReqByUserByDate) (*api.RepEvents, error) {
	events, err := s.app.ListEventsDay(ctx, *req.UserID, req.Date.AsTime().Local())
	if err != nil {
		return nil, err
	}

	rep := api.RepEvents{}
	rep.Event = make([]*api.Event, len(events))
	for i, event := range events {
		event := event
		rep.Event[i] = s.APIEventFromEvent(&event)
	}
	return &rep, nil
}

func (s Service) ListEventsWeek(ctx context.Context, req *api.ReqByUserByDate) (*api.RepEvents, error) {
	events, err := s.app.ListEventsWeek(ctx, *req.UserID, req.Date.AsTime().Local())
	if err != nil {
		return nil, err
	}

	rep := api.RepEvents{}
	rep.Event = make([]*api.Event, len(events))
	for i, event := range events {
		event := event
		rep.Event[i] = s.APIEventFromEvent(&event)
	}
	return &rep, nil
}

func (s Service) ListEventsMonth(ctx context.Context, req *api.ReqByUserByDate) (*api.RepEvents, error) {
	events, err := s.app.ListEventsMonth(ctx, *req.UserID, req.Date.AsTime().Local())
	if err != nil {
		return nil, err
	}

	rep := api.RepEvents{}
	rep.Event = make([]*api.Event, len(events))
	for i, event := range events {
		event := event
		rep.Event[i] = s.APIEventFromEvent(&event)
	}
	return &rep, nil
}

func New(log Logger, app Application, conf Conf, basesrv *grpc.Server) *Service {
	return &Service{app: app, conf: conf, log: log, basesrv: basesrv}
}

func (s *Service) Start(context.Context) error {
	addr := net.JoinHostPort(s.conf.Host, s.conf.Port)
	dial, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	api.RegisterCalendarServer(s.basesrv, s)

	s.log.Infof("GRPC-server started on:%v\n", addr)
	if err := s.basesrv.Serve(dial); err != nil {
		return err
	}

	return nil
}

func (s *Service) Stop(context.Context) error {
	s.basesrv.GracefulStop()
	s.log.Infof("GRPC-server shutdown\n")
	return nil
}
