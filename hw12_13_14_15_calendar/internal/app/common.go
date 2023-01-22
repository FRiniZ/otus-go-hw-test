package app

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrID             = errors.New("wrong ID")
	ErrUserID         = errors.New("wrong UserID")
	ErrTitle          = errors.New("wrong Title")
	ErrDescription    = errors.New("wrong Description")
	ErrOnTime         = errors.New("wrong OnTime")
	ErrOffTime        = errors.New("wrong OffTime")
	ErrNotifyTime     = errors.New("wrong NotifyTime")
	ErrEventNotFound  = errors.New("event not found")
	ErrTooLongCloseDB = errors.New("too long close db")
)

type Logger interface {
	Fatalf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warningf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
}

func exitfail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
