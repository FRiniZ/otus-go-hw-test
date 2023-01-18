package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
)

type mapEvent map[int64]*storage.Event

type Storage struct {
	data mapEvent
	mu   sync.RWMutex
}

var (
	GenID            = int64(1)
	ErrEventNotFound = errors.New("event not found")
)

func getNewIDUnsafe() int64 {
	ret := GenID
	GenID++
	return ret
}

func New() *Storage {
	return &Storage{data: make(mapEvent), mu: sync.RWMutex{}}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) inTimeSpan(start, end, check time.Time) bool {
	switch {
	case check.Equal(start):
		return true
	case check.Equal(end):
		return true
	case check.After(start) && check.Before(start):
		return true
	}
	return false
}

func (s *Storage) IsBusyDateTimeRange(ctx context.Context, userID int64, onTime, offTime time.Time) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.data {
		// Если дата начала события или дата окончания события входят в
		// диапазон другого события или даты равны то считаем что время занято
		if v.UserID == userID &&
			(s.inTimeSpan(v.OnTime, v.OffTime, onTime) ||
				s.inTimeSpan(v.OnTime, v.OffTime, offTime)) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) InsertEvent(ctx context.Context, e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e.ID = getNewIDUnsafe()
	s.data[e.ID] = e
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[e.ID]; !ok {
		return ErrEventNotFound
	}
	s.data[e.ID] = e

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, e.ID)
	return nil
}

func (s *Storage) ListEvents(ctx context.Context, userID int64) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sliceE := []storage.Event{}
	for _, v := range s.data {
		if v.UserID == userID {
			sliceE = append(sliceE, *v)
		}
	}

	return sliceE, nil
}

func (s *Storage) firstDayOfWeek(t time.Time) time.Time {
	for t.Weekday() != time.Monday {
		t = t.AddDate(0, 0, -1)
	}
	return t
}

func (s *Storage) firstDayOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 0, -t.Day()+1)
}

func (s *Storage) ListEventsWeek(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	return s.ListEventsDay(ctx, userID, s.firstDayOfWeek(date))
}

func (s *Storage) ListEventsMonth(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	return s.ListEventsDay(ctx, userID, s.firstDayOfMonth(date))
}

func (s *Storage) ListEventsDay(ctx context.Context, userID int64, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sliceE := []storage.Event{}
	for _, v := range s.data {
		if v.UserID == userID && s.inTimeSpan(date, v.OnTime, v.OffTime) {
			sliceE = append(sliceE, *v)
		}
	}

	return sliceE, nil
}

func (s *Storage) LookupEvent(ctx context.Context, eID int64) (storage.Event, error) {
	var event storage.Event
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.data[eID]; ok {
		event = *e
		return event, nil
	}

	return event, ErrEventNotFound
}
