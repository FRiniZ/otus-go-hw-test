package logger

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
)

const (
	L0 = iota // Error
	L1        // Warn
	L2        // Info
	L3        // Debug
)

type Logger struct {
	logLevel int
	out      io.Writer
	mu       *sync.Mutex
	cancel   context.CancelFunc
}

func New(level string, out io.Writer, cancel context.CancelFunc) *Logger {
	switch strings.ToUpper(level) {
	case "ERROR":
		return &Logger{logLevel: L0, mu: &sync.Mutex{}, out: out, cancel: cancel}
	case "WARN":
		return &Logger{logLevel: L1, mu: &sync.Mutex{}, out: out, cancel: cancel}
	case "INFO":
		return &Logger{logLevel: L2, mu: &sync.Mutex{}, out: out, cancel: cancel}
	case "DEBUG":
		return &Logger{logLevel: L3, mu: &sync.Mutex{}, out: out, cancel: cancel}
	default:
		panic("unrecognized log_level")
	}
}

func (l *Logger) printf(format string, a ...interface{}) {
	l.mu.Lock()
	_, err := fmt.Fprintf(l.out, format, a...)
	l.mu.Unlock()
	if err != nil {
		l.Panicf("Panic: %v", err)
	}
}

func (l *Logger) Panicf(format string, a ...interface{}) {
	l.printf("PANIC:"+format, a)
	l.cancel()
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	if l.logLevel >= L0 {
		l.printf("ERROR:"+format, a...)
	}
}

func (l *Logger) Warningf(format string, a ...interface{}) {
	if l.logLevel >= L1 {
		l.printf("WARN:"+format, a...)
	}
}

func (l *Logger) Infof(format string, a ...interface{}) {
	if l.logLevel >= L2 {
		l.printf("INFO:"+format, a...)
	}
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.logLevel >= L3 {
		l.printf("DEBUG:"+format, a...)
	}
}
