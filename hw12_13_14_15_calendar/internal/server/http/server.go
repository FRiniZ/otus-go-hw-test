package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	srv http.Server
	app Application
	log Logger
}

type Logger interface {
	Panicf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warningf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
}

type Application interface {
	// TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{log: logger, app: app}
}

func (s *Server) doNothing(w http.ResponseWriter, r *http.Request) {
	// empty function
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello-world"))
}

func (s *Server) Start(ctx context.Context, addr string) error {
	midLogger := NewMiddlewareLogger(s.log)
	mux := http.NewServeMux()
	mux.Handle("/", midLogger.loggingMiddleware(http.HandlerFunc(s.ServeHTTP)))

	// to avoid twice handling
	mux.HandleFunc("/favicon.ico", s.doNothing)

	s.srv = http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 2 * time.Second}
	s.srv.ListenAndServe()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		s.log.Errorf("HTTP server Shutdown: %w\n", err)
		return err
	}
	s.log.Infof("HTTP server Shutdown OK\n")
	return nil
}
