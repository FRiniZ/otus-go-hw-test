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
	DeleteEvent(context.Context, int64) error
	LookupEvent(context.Context, int64) (storage.Event, error)
	ListEvents(context.Context, int64) ([]storage.Event, error)
	ListEventsDay(context.Context, int64, time.Time) ([]storage.Event, error)
	ListEventsWeek(context.Context, int64, time.Time) ([]storage.Event, error)
	ListEventsMonth(context.Context, int64, time.Time) ([]storage.Event, error)
}

type reqByID struct {
	ID int64 `json:"id"`
}

type reqByUser struct {
	UserID int64 `json:"userid"`
}

type reqByUserByDate struct {
	UserID int64     `json:"userid"`
	Date   time.Time `json:"date"`
}

func New(log Logger, app Application, conf Conf, cancel context.CancelFunc) *Server {
	return &Server{log: log, app: app, conf: conf, cancel: cancel}
}

func (s *Server) doNothing(w http.ResponseWriter, r *http.Request) {
	// empty function
}

func (s *Server) helperDecode(stream io.ReadCloser, w http.ResponseWriter, data interface{}) error {
	decoder := json.NewDecoder(stream)
	if err := decoder.Decode(&data); err != nil {
		s.log.Errorf("Can't decode json:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't decode json:%v\"}\n", err)))
		return err
	}
	return nil
}

func (s *Server) InsertEvent(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var req reqByID
	if err := s.helperDecode(r.Body, w, &req); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.app.DeleteEvent(ctx, req.ID)
	if err != nil {
		s.log.Errorf("DeleteEvent:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't DeleteEvent:%v\"}\n", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"msg\": \"Deleted\"}\n"))
}

func (s *Server) LookupEvent(w http.ResponseWriter, r *http.Request) {
	var req reqByID
	if err := s.helperDecode(r.Body, w, &req); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	eventFound, err := s.app.LookupEvent(ctx, req.ID)
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

func (s *Server) ListEvents(w http.ResponseWriter, r *http.Request) {
	var req reqByUser
	if err := s.helperDecode(r.Body, w, &req); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	eventsFound, err := s.app.ListEvents(ctx, req.UserID)
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
			s.log.Errorf("ListEvents:%v\n", err)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't ListEvents:%v\"},\n", err)))
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

func (s *Server) ListEventsDay(w http.ResponseWriter, r *http.Request) { //nolint
	// lines are duplicate of `internal/server/http/server.go:233-263` (dupl)
	var req reqByUserByDate
	if err := s.helperDecode(r.Body, w, &req); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	eventsFound, err := s.app.ListEventsDay(ctx, req.UserID, req.Date)
	if err != nil {
		s.log.Errorf("ListEventsDay:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't ListEventsDay:%v\"}\n", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[\n"))
	for i, e := range eventsFound {
		jevent, err := json.Marshal(e)
		if err != nil {
			s.log.Errorf("ListEventsDay:%v\n", err)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't ListEventsDay:%v\"},\n", err)))
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

func (s *Server) ListEventsWeek(w http.ResponseWriter, r *http.Request) { //nolint
	var req reqByUserByDate
	if err := s.helperDecode(r.Body, w, &req); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	eventsFound, err := s.app.ListEventsWeek(ctx, req.UserID, req.Date)
	if err != nil {
		s.log.Errorf("ListEventsWeek:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't ListEventsWeek:%v\"}\n", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[\n"))
	for i, e := range eventsFound {
		jevent, err := json.Marshal(e)
		if err != nil {
			s.log.Errorf("ListEventsWeek:%v\n", err)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't ListEventsWeek:%v\"},\n", err)))
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

func (s *Server) ListEventsMonth(w http.ResponseWriter, r *http.Request) { //nolint
	var req reqByUserByDate
	if err := s.helperDecode(r.Body, w, &req); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	eventsFound, err := s.app.ListEventsMonth(ctx, req.UserID, req.Date)
	if err != nil {
		s.log.Errorf("ListEventsMonth:%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't ListEventsMonth:%v\"}\n", err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[\n"))
	for i, e := range eventsFound {
		jevent, err := json.Marshal(e)
		if err != nil {
			s.log.Errorf("ListEventsMonth:%v\n", err)
			w.Write([]byte(fmt.Sprintf("{\"error\": \"Can't ListEventsMonth:%v\"},\n", err)))
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

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.conf.Host, s.conf.Port)
	midLogger := NewMiddlewareLogger()
	mux := http.NewServeMux()

	mux.Handle("/InsertEvent", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.InsertEvent))))
	mux.Handle("/UpdateEvent", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.UpdateEvent))))
	mux.Handle("/DeleteEvent", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.DeleteEvent))))
	mux.Handle("/LookupEvent", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.LookupEvent))))
	mux.Handle("/ListEvents", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.ListEvents))))
	mux.Handle("/ListEventsDay", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.ListEventsDay))))
	mux.Handle("/ListEventsWeek", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.ListEventsWeek))))
	mux.Handle("/ListEventsMonth", midLogger.setCommonHeadersMiddleware(
		midLogger.loggingMiddleware(http.HandlerFunc(s.ListEventsMonth))))

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
