package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

const (
	LevelError = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

var ErrLogLevel = errors.New("unrecognized log_level")

type Conf struct {
	Level string `toml:"level"`
}

type Logger struct {
	logLevel int
	out      io.Writer
	mu       *sync.Mutex
}

func NewLogger(level string, out io.Writer) *Logger {
	switch strings.ToUpper(level) {
	case "ERROR":
		return &Logger{logLevel: LevelError, mu: &sync.Mutex{}, out: out}
	case "WARN":
		return &Logger{logLevel: LevelWarn, mu: &sync.Mutex{}, out: out}
	case "INFO":
		return &Logger{logLevel: LevelInfo, mu: &sync.Mutex{}, out: out}
	case "DEBUG":
		return &Logger{logLevel: LevelDebug, mu: &sync.Mutex{}, out: out}
	default:
		fmt.Fprintln(os.Stderr, "unrecognized log_level")
		os.Exit(1)
	}
	return nil
}

func (l *Logger) printf(format string, a ...interface{}) {
	l.mu.Lock()
	_, err := fmt.Fprintf(l.out, format, a...)
	l.mu.Unlock()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: Fprintf : %v", err)
		os.Exit(1)
	}
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.printf("Fatal:"+format, a)
	os.Exit(1)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	if l.logLevel >= LevelError {
		l.printf("ERROR:"+format, a...)
	}
}

func (l *Logger) Warningf(format string, a ...interface{}) {
	if l.logLevel >= LevelWarn {
		l.printf("WARN:"+format, a...)
	}
}

func (l *Logger) Infof(format string, a ...interface{}) {
	if l.logLevel >= LevelInfo {
		l.printf("INFO:"+format, a...)
	}
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	if l.logLevel >= LevelDebug {
		l.printf("DEBUG:"+format, a...)
	}
}
