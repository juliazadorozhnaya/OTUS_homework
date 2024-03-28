package memorystorage

import (
	"context"
	"errors"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/model"
	"sync"

	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.Mutex
	events map[string]model.Event
	users  map[string]model.User
}

var ErrEventNotFound = errors.New("event not found")

func New() *Storage {
	return &Storage{
		events: make(map[string]model.Event),
		users:  make(map[string]model.User),
	}
}

func (s *Storage) CreateUser(ctx context.Context, user model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = uuid.New().String()
	s.users[user.ID] = user

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.users, userID)
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = uuid.New().String()
	s.events[event.ID] = event

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, eventID)
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return ErrEventNotFound
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) SelectEvents(ctx context.Context) ([]model.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := make([]model.Event, 0)
	for _, event := range s.events {
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) SelectUsers(ctx context.Context) ([]model.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := make([]model.User, 0)
	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}
