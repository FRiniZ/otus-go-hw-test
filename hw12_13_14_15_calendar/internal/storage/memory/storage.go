package memorystorage

import (
	"sync"

	"github.com/FRiniZ/otus-go-hw-test/hw12_13_14_15_calendar/internal/storage"
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

func (s *Storage) InsertEvent(e *storage.Event) error {
	if err := storage.CheckingEvent(e); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e.ID = getNewIDUnsafe()
	s.data[e.ID] = e
	return nil
}

func (s *Storage) UpdateEvent(e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[e.ID]; ok {
		s.data[e.ID] = e
	} else {
		return nil
	}

	return nil
}

func (s *Storage) DeleteEvent(e *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, e.ID)
	return nil
}

func (s *Storage) ListEvents(userID int64) ([]storage.Event, error) {
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

func (s *Storage) LookupEvent(eID int64) (event storage.Event, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.data[eID]; ok {
		event = *e
	}

	return event, err
}
