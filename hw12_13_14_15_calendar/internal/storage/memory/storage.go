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
	data  mapEvent
	mu    sync.RWMutex
	genID int64
}

var (
	ErrEventNotFound   = errors.New("event not found")
	ErrDataRangeIsBusy = errors.New("data is busy")
)

func (s *Storage) getNewIDUnsafe() int64 {
	ret := s.genID
	s.genID++
	return ret
}

func New() *Storage {
	return &Storage{data: make(mapEvent), mu: sync.RWMutex{}, genID: 1}
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
	case check.After(start) && check.Before(end):
		return true
	}
	return false
}

func (s *Storage) IsBusyDateTimeRange(ctx context.Context, id, userID int64, onTime, offTime time.Time) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.data {
		if v.UserID == userID && v.ID != id &&
			(s.inTimeSpan(v.OnTime, v.OffTime, onTime) ||
				s.inTimeSpan(v.OnTime, v.OffTime, offTime)) {
			return ErrDataRangeIsBusy
		}
	}
	return nil
}

func (s *Storage) InsertEvent(ctx context.Context, e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e.ID = s.getNewIDUnsafe()
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

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, id)
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

func (s *Storage) ListEventsRange(ctx context.Context, userID int64, begin, end time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sliceE := []storage.Event{}

	for _, v := range s.data {
		if v.UserID == userID &&
			(s.inTimeSpan(begin, end, v.OnTime) ||
				s.inTimeSpan(begin, end, v.OffTime)) {
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

func (s *Storage) ListEventsDayOfNotice(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sliceE := []storage.Event{}

	for _, v := range s.data {
		if !v.Notified && (v.NotifyTime.Before(date) || v.NotifyTime.Equal(date)) {
			sliceE = append(sliceE, *v)
		}
	}
	return sliceE, nil
}

func (s *Storage) UpdateEventNotified(ctx context.Context, eventid int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[eventid]; !ok {
		return ErrEventNotFound
	}
	s.data[eventid].Notified = true
	return nil
}
