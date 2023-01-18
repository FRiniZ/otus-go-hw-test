package grpcservice

import (
	context "context"
	"net"
	"strings"
	"time"

	api "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/api/stub"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
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
	DeleteEvent(context.Context, *storage.Event) error
	LookupEvent(context.Context, int64) (storage.Event, error)
	ListEvents(context.Context, int64) ([]storage.Event, error)
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
		event.OnTime = apiEvent.OnTime.AsTime()
	}
	if err := apiEvent.OffTime.CheckValid(); err == nil {
		event.OffTime = apiEvent.OffTime.AsTime()
	}
	if err := apiEvent.NotifyTime.CheckValid(); err == nil {
		event.NotifyTime = apiEvent.NotifyTime.AsTime()
	}

	return &event
}

// DeleteEventV1 implements api.CalendarServer.
func (Service) DeleteEventV1(context.Context, *api.RequestV1) (*api.ReplyV1, error) {
	panic("unimplemented")
}

// InsertEventV1 implements api.CalendarServer.
func (s Service) InsertEventV1(ctx context.Context, req *api.RequestV1) (*api.ReplyV1, error) {
	defer s.Log(ctx)

	event := s.EventFromAPIEvent(req.Event)
	if err := s.app.InsertEvent(ctx, event); err != nil {
		return &api.ReplyV1{}, err
	}
	rep := api.ReplyV1{}
	rep.Event = append(rep.Event, s.APIEventFromEvent(event))
	return &rep, nil
}

// UpdateEventV1 implements api.CalendarServer.
func (s Service) UpdateEventV1(ctx context.Context, req *api.RequestV1) (*api.ReplyV1, error) {
	defer s.Log(ctx)
	event := s.EventFromAPIEvent(req.Event)
	if err := s.app.UpdateEvent(ctx, event); err != nil {
		return &api.ReplyV1{}, err
	}

	rep := api.ReplyV1{}
	rep.Event = append(rep.Event, s.APIEventFromEvent(event))
	return &rep, nil
}

// ListEventsV1 implements api.CalendarServer.
func (s Service) ListEventsV1(ctx context.Context, req *api.RequestV1) (*api.ReplyV1, error) {
	defer s.Log(ctx)
	userID := *req.Event.UserID
	events, err := s.app.ListEvents(ctx, userID)
	if err != nil {
		return &api.ReplyV1{}, err
	}

	rep := api.ReplyV1{}
	rep.Event = make([]*api.Event, len(events))
	for i, event := range events {
		event := event
		rep.Event[i] = s.APIEventFromEvent(&event)
	}
	return &rep, nil
}

// LookupEventV1 implements api.CalendarServer.
func (s Service) LookupEventV1(ctx context.Context, req *api.RequestV1) (*api.ReplyV1, error) {
	defer s.Log(ctx)
	eventID := *req.Event.ID
	event, err := s.app.LookupEvent(ctx, eventID)
	if err != nil {
		return &api.ReplyV1{}, err
	}

	rep := api.ReplyV1{}
	rep.Event = append(rep.Event, s.APIEventFromEvent(&event))
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
	return grpc.ErrServerStopped
}

func UnaryLoggerEnricherInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) { //nolint
	ctxV := context.WithValue(ctx, KeyMethodID, info.FullMethod)
	// Calls the handler
	h, err := handler(ctxV, req)
	return h, err
}

func (s *Service) Log(ctx context.Context) {
	var b strings.Builder
	ip, _ := peer.FromContext(ctx)
	method := ctx.Value(KeyMethodID).(string)
	md, ok := metadata.FromIncomingContext(ctx)
	userAgent := "unknown"

	if ok {
		userAgent = md["user-agent"][0]
	}

	b.WriteString(ip.Addr.String())
	b.WriteString(" ")
	b.WriteString(time.Now().Format("02/Jan/2006:15:04:05 -0700"))
	b.WriteString(" ")
	b.WriteString(method)
	b.WriteString(" ")
	b.WriteString(userAgent)
	b.WriteString("\"\n")

	s.log.Infof(b.String())
}
