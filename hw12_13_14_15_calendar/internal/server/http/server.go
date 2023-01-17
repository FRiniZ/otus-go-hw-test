package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
)

type ctxKeyID int

const (
	KeyLoggerID ctxKeyID = iota
)

type Conf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type ResponceErr struct {
	Msg string
}

type Server struct {
	srv    http.Server
	log    Logger
	app    Application
	conf   Conf
	cancel context.CancelFunc
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

func NewServer(log Logger, app Application, conf Conf, cancel context.CancelFunc) *Server {
	return &Server{log: log, app: app, conf: conf, cancel: cancel}
}

func (s *Server) doNothing(w http.ResponseWriter, r *http.Request) {
	// empty function
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello-world"))
}

func (s *Server) helperDecode(stream io.ReadCloser, w http.ResponseWriter, event *storage.Event) error {
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&event); err != nil {
		s.log.Errorf("Can't decode json:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't decode json:%v\"}\n", err)))
		return err
	}
	return nil
}

func (s *Server) InsertEventV1(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	if err := s.helperDecode(r.Body, w, &event); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.app.InsertEvent(ctx, &event)
	if err != nil {
		s.log.Errorf("InsertEvent:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't InsertEvent:%v\"}\n", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"msg\": \"Inserted\"}\n"))
}

func (s *Server) UpdateEventV1(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	if err := s.helperDecode(r.Body, w, &event); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.app.UpdateEvent(ctx, &event)
	if err != nil {
		s.log.Errorf("UpdateEvent:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't UpdateEvent:%v\"}\n", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"msg\": \"Updated\"}\n"))
}

func (s *Server) DeleteEventV1(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	if err := s.helperDecode(r.Body, w, &event); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.app.DeleteEvent(ctx, &event)
	if err != nil {
		s.log.Errorf("DeleteEvent:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't DeleteEvent:%v\"}\n", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"msg\": \"Deleted\"}\n"))
}

func (s *Server) LookupEventV1(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	if err := s.helperDecode(r.Body, w, &event); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	eventFound, err := s.app.LookupEvent(ctx, event.ID)
	if err != nil {
		s.log.Errorf("LookupEvent:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't LookupEvent:%v\"}\n", err)))
		return
	}
	jevent, err := json.Marshal(eventFound)
	if err != nil {
		s.log.Errorf("LookupEvent:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't LookupEvent:%v\"}\n", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jevent)
	w.Write([]byte("\n"))
}

func (s *Server) ListEventsV1(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	if err := s.helperDecode(r.Body, w, &event); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	eventsFound, err := s.app.ListEvents(ctx, event.UserID)
	if err != nil {
		s.log.Errorf("ListEvents:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't ListEvents:%v\"}\n", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[\n"))
	for i, e := range eventsFound {
		jevent, err := json.Marshal(e)
		if err != nil {
			s.log.Errorf("LookupEvent:%v\n", err)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't LookupEvent:%v\"},\n", err)))
		}
		w.Write(jevent)
		if i+1 == len(eventsFound) {
			w.Write([]byte("\n"))
		} else {
			w.Write([]byte(",\n"))
		}
	}
	w.Write([]byte("]\n"))
}

// Временное решение для проверки отработки закрытия.
func (s *Server) Shutdown(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("close"))
	s.cancel()
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.conf.Host, s.conf.Port)
	midLogger := NewMiddlewareLogger()
	mux := http.NewServeMux()

	mux.Handle("/", midLogger.loggingMiddleware(http.HandlerFunc(s.ServeHTTP)))
	mux.Handle("/shutdown", midLogger.loggingMiddleware(http.HandlerFunc(s.Shutdown)))
	mux.Handle("/InsertEventV1", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.InsertEventV1))))
	mux.Handle("/UpdateEventV1", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.UpdateEventV1))))
	mux.Handle("/DeleteEventV1", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.DeleteEventV1))))
	mux.Handle("/LookupEventV1", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.LookupEventV1))))
	mux.Handle("/ListEventsV1", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.ListEventsV1))))

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

	s.log.Infof("HTTP-server started on:%v\n", addr)

	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Infof("HTTP-server shutdown\n")
	return s.srv.Shutdown(ctx)
}
