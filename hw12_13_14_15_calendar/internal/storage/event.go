package storage

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrUserID      = errors.New("wrong UserID")
	ErrTitle       = errors.New("wrong Title")
	ErrDescription = errors.New("wrong Description")
	ErrOnTime      = errors.New("wrong OnTime")
	ErrOffTime     = errors.New("wrong OffTime")
	ErrNotifyTime  = errors.New("wrong NotifyTime")
)

type Event struct {
	ID          int64 `faker:"-"`
	UserID      int64
	Title       string `faker:"word"`
	Description string `faker:"sentence"`
	OnTime      time.Time
	OffTime     time.Time
	NotifyTime  time.Time
}

func CheckingEvent(e *Event) error {
	if e.UserID == 0 {
		return fmt.Errorf("%w(UserID is %v)", ErrUserID, e.UserID)
	}

	if len(e.Title) > 150 {
		return fmt.Errorf("%w(len %v, must be <=150)", ErrTitle, len(e.Title))
	}

	if e.OnTime.IsZero() {
		return fmt.Errorf("%w(empty OnTime)", ErrOnTime)
	}

	if !e.OffTime.IsZero() {
		if e.OffTime.Before(e.OnTime) {
			return fmt.Errorf("%w(OffTime before OnTime)", ErrOffTime)
		}
		if e.OffTime.Equal(e.OnTime) {
			return fmt.Errorf("%w(OffTime equal OnTime)", ErrOffTime)
		}
	}

	if !e.NotifyTime.IsZero() {
		if e.NotifyTime.After(e.OffTime) {
			return fmt.Errorf("%w(NotifyTime after OffTime)", ErrNotifyTime)
		}
		if e.NotifyTime.Before(e.OnTime) {
			return fmt.Errorf("%w(NotifyTime before OnTime)", ErrNotifyTime)
		}
	}

	return nil
}
