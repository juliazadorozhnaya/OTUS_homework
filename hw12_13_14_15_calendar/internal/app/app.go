package app

import (
	"context"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/model"
	"sync"
)

type Calendar struct {
	storage Storage
	mutex   sync.RWMutex
}

type Storage interface {
	CreateUser(context.Context, model.User) error
	SelectUsers(ctx context.Context) ([]model.User, error)
	DeleteUser(context.Context, string) error

	CreateEvent(context.Context, model.Event) error
	SelectEvents(context.Context) ([]model.Event, error)
	UpdateEvent(context.Context, model.Event) error
	DeleteEvent(context.Context, string) error
}

func New(storage Storage) *Calendar {
	return &Calendar{
		storage: storage,
		mutex:   sync.RWMutex{},
	}
}

func (calendar *Calendar) CreateUser(ctx context.Context, User model.User) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	user := model.User{
		ID:        User.ID,
		FirstName: User.FirstName,
		LastName:  User.LastName,
		Email:     User.Email,
		Age:       User.Age,
	}

	return calendar.storage.CreateUser(ctx, user)
}

func (calendar *Calendar) SelectUsers(ctx context.Context) ([]model.User, error) {
	calendar.mutex.RLock()
	defer calendar.mutex.RUnlock()

	users := make([]model.User, 0)

	storageUsers, err := calendar.storage.SelectUsers(ctx)
	if err != nil {
		return users, err
	}

	for _, storageUser := range storageUsers {
		users = append(users, model.User{
			ID:        storageUser.ID,
			FirstName: storageUser.FirstName,
			LastName:  storageUser.LastName,
			Email:     storageUser.Email,
			Age:       storageUser.Age,
		})
	}

	return users, nil
}

func (calendar *Calendar) DeleteUser(ctx context.Context, id string) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	return calendar.storage.DeleteUser(ctx, id)
}

func (calendar *Calendar) CreateEvent(ctx context.Context, Event model.Event) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	event := model.Event{
		ID:           Event.ID,
		Title:        Event.Title,
		Description:  Event.Description,
		UserID:       Event.UserID,
		Beginning:    Event.Beginning,
		Finish:       Event.Finish,
		Notification: Event.Notification,
	}

	return calendar.storage.CreateEvent(ctx, event)
}

func (calendar *Calendar) UpdateEvent(ctx context.Context, Event model.Event) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	event := model.Event{
		ID:           Event.ID,
		Title:        Event.Title,
		Description:  Event.Description,
		UserID:       Event.UserID,
		Beginning:    Event.Beginning,
		Finish:       Event.Finish,
		Notification: Event.Notification,
	}

	return calendar.storage.UpdateEvent(ctx, event)
}

func (calendar *Calendar) DeleteEvent(ctx context.Context, id string) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	return calendar.storage.DeleteEvent(ctx, id)
}

func (calendar *Calendar) SelectEvents(ctx context.Context) ([]model.Event, error) {
	calendar.mutex.RLock()
	defer calendar.mutex.RUnlock()

	events := make([]model.Event, 0)

	storageEvents, err := calendar.storage.SelectEvents(ctx)
	if err != nil {
		return events, err
	}

	for _, storageEvent := range storageEvents {
		events = append(events, model.Event{
			ID:           storageEvent.ID,
			Title:        storageEvent.Title,
			Description:  storageEvent.Description,
			Beginning:    storageEvent.Beginning,
			Notification: storageEvent.Notification,
			UserID:       storageEvent.UserID,
		})
	}

	return events, nil
}
