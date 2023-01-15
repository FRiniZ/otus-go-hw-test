package internalhttp

import (
	"net/http"
)

type MiddlewareLogger struct{}

func NewMiddlewareLogger() *MiddlewareLogger {
	return &MiddlewareLogger{}
}

func (m *MiddlewareLogger) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := r.Context().Value(KeyLoggerID).(Logger)
		rwc := NewResponseWriterCounter(w, r)
		next.ServeHTTP(rwc, r)
		l.Debugf(rwc.String())
	})
}
