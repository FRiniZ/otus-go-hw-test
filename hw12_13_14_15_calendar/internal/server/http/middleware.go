package internalhttp

import (
	"net/http"
)

type MiddlewareLogger struct {
	Logger
}

func NewMiddlewareLogger(log Logger) *MiddlewareLogger {
	return &MiddlewareLogger{
		Logger: log,
	}
}

func (m *MiddlewareLogger) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rwc := NewResponseWriterCounter(w, r)
		next.ServeHTTP(rwc, r)
		m.Debugf(rwc.String())
	})
}
