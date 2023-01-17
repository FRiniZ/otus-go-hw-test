package memorystorage

import (
	"context"
	"sync"

	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/app"
	"github.com/FRiniZ/otus-go-hw-test/hw12_calendar/internal/storage"
)

type mapEvent map[int64]*storage.Event

type Storage struct {
	data mapEvent
	mu   sync.RWMutex
}

var GenID = int64(1)

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
		return app.ErrEventNotFound
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

func (s *Storage) LookupEvent(ctx context.Context, eID int64) (storage.Event, error) {
	var event storage.Event
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.data[eID]; ok {
		event = *e
		return event, nil
	}

	return event, app.ErrEventNotFound
}
