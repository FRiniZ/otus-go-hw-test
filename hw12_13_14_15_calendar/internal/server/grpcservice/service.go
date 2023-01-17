package grpcservice

import (
	context "context"
	"net"

	api "github.com/FRiniZ/otus-go-hw-test/hw12_calendar/api/stub"
	"google.golang.org/grpc"
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

type GRPCService struct {
	log     Logger
	conf    Conf
	basesrv *grpc.Server
	api.UnimplementedCalendarServer
}

// DeleteEventV1 implements api.CalendarServer.
func (GRPCService) DeleteEventV1(context.Context, *api.RequestV1) (*api.ReplyV1, error) {
	panic("unimplemented")
}

// InserEventV1 implements api.CalendarServer.
func (GRPCService) InserEventV1(context.Context, *api.RequestV1) (*api.ReplyV1, error) {
	panic("unimplemented")
}

// ListEventsV1 implements api.CalendarServer.
func (GRPCService) ListEventsV1(context.Context, *api.RequestV1) (*api.ReplyV1, error) {
	panic("unimplemented")
}

// LookupEventV1 implements api.CalendarServer.
func (GRPCService) LookupEventV1(context.Context, *api.RequestV1) (*api.ReplyV1, error) {
	panic("unimplemented")
}

// UpdateEventV1 implements api.CalendarServer.
func (GRPCService) UpdateEventV1(context.Context, *api.RequestV1) (*api.ReplyV1, error) {
	panic("unimplemented")
}

func NewGRPCService(log Logger, conf Conf) *GRPCService {
	return &GRPCService{conf: conf, log: log, basesrv: grpc.NewServer()}
}

func (s *GRPCService) Start(context.Context) error {
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

func (s *GRPCService) Stop(context.Context) error {
	s.basesrv.GracefulStop()
	s.log.Infof("GRPC-server shutdown\n")
	return grpc.ErrServerStopped
}
