package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/model"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]model.Event
	users  map[string]model.User
}

var (
	ErrEventNotFound = errors.New("event not found")
	ErrUserNotFound  = errors.New("user not found")
)

func New() *Storage {
	return &Storage{
		events: make(map[string]model.Event),
		users:  make(map[string]model.User),
	}
}

// CreateUser создает нового пользователя и добавляет его в map пользователей.
func (s *Storage) CreateUser(_ context.Context, user model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = uuid.New().String()
	s.users[user.ID] = user

	return nil
}

// DeleteUser удаляет пользователя по его идентификатору.
func (s *Storage) DeleteUser(_ context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[userID]; !ok {
		return ErrUserNotFound
	}

	delete(s.users, userID)
	return nil
}

// CreateEvent cоздает новое событие и добавляет его в map событий.
func (s *Storage) CreateEvent(_ context.Context, event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = uuid.New().String()
	s.events[event.ID] = event

	return nil
}

// DeleteEvent удаляет событие по его идентификатору.
func (s *Storage) DeleteEvent(_ context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventID]; !ok {
		return ErrEventNotFound
	}

	delete(s.events, eventID)
	return nil
}

// UpdateEvent обновляет существующее событие в map событий.
func (s *Storage) UpdateEvent(_ context.Context, event model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return ErrEventNotFound
	}

	s.events[event.ID] = event
	return nil
}

// SelectEvents возвращает все события.
func (s *Storage) SelectEvents(_ context.Context) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]model.Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}

	return events, nil
}

// SelectUsers возвращает всех пользователей.
func (s *Storage) SelectUsers(_ context.Context) ([]model.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]model.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	return users, nil
}

// SelectEventsForDay возвращает события на указанный день.
func (s *Storage) SelectEventsForDay(_ context.Context, date time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]model.Event, 0)
	for _, event := range s.events {
		if event.Beginning.Year() == date.Year() && event.Beginning.YearDay() == date.YearDay() {
			events = append(events, event)
		}
	}

	return events, nil
}

// SelectEventsForWeek возвращает события на неделю, начиная с указанной даты.
func (s *Storage) SelectEventsForWeek(_ context.Context, startDate time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]model.Event, 0)
	endDate := startDate.AddDate(0, 0, 7)
	for _, event := range s.events {
		if event.Beginning.After(startDate) && event.Beginning.Before(endDate) {
			events = append(events, event)
		}
	}

	return events, nil
}

// SelectEventsForMonth возвращает события на месяц, начиная с указанной даты.
func (s *Storage) SelectEventsForMonth(_ context.Context, startDate time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]model.Event, 0)
	endDate := startDate.AddDate(0, 1, 0)
	for _, event := range s.events {
		if event.Beginning.After(startDate) && event.Beginning.Before(endDate) {
			events = append(events, event)
		}
	}

	return events, nil
}

// SelectEventsByTime возвращает список событий, которые должны быть уведомлены в указанное время.
func (s *Storage) SelectEventsByTime(_ context.Context, t time.Time) ([]model.Event, error) {
	events := make([]model.Event, 0)

	for _, event := range s.events {
		if event.Notification.Equal(t) {
			events = append(events, event)
		}
	}

	return events, nil
}
