package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"
)

type ctxKeyID int

const (
	KeyLoggerID ctxKeyID = iota
)

type ServerConf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type Server struct {
	srv  http.Server
	app  Application
	log  Logger
	conf ServerConf
}

type Logger interface {
	Fatalf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warningf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
}

type Application interface{}

func NewServer(logger Logger, app Application, conf ServerConf) *Server {
	return &Server{log: logger, app: app, conf: conf}
}

func (s *Server) doNothing(w http.ResponseWriter, r *http.Request) {
	// empty function
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello-world"))
}

// Временное решение для проверки отработки закрытия.
func (s *Server) Close(w http.ResponseWriter, r *http.Request) {
	s.Stop(r.Context())
	w.Write([]byte("close"))
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.conf.Host, s.conf.Port)

	midLogger := NewMiddlewareLogger()
	mux := http.NewServeMux()
	mux.Handle("/", midLogger.loggingMiddleware(http.HandlerFunc(s.ServeHTTP)))
	mux.Handle("/close", midLogger.loggingMiddleware(http.HandlerFunc(s.Close)))

	// to avoid twice handling
	mux.HandleFunc("/favicon.ico", s.doNothing)

	s.srv = http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			bCtx := context.WithValue(ctx, KeyLoggerID, s.log)
			return bCtx
		},
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}
	s.log.Infof("HTTP server Shutdown\n")
	return nil
}
