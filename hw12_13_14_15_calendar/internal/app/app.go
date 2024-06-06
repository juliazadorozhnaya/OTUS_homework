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
	CreateUser(ctx context.Context, User model.User) error
	SelectUsers(ctx context.Context) ([]model.User, error)
	DeleteUser(ctx context.Context, id string) error

	CreateEvent(ctx context.Context, Event model.Event) error
	SelectEvents(ctx context.Context) ([]model.Event, error)
	UpdateEvent(ctx context.Context, Event model.Event) error
	DeleteEvent(ctx context.Context, id string) error
}

func New(storage Storage) *Calendar {
	return &Calendar{
		storage: storage,
		mutex:   sync.RWMutex{},
	}
}

func (calendar *Calendar) CreateUser(ctx context.Context, user model.IUser) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	storageUser := model.User{
		ID:        user.GetID(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Email:     user.GetEmail(),
		Age:       user.GetAge(),
	}

	return calendar.storage.CreateUser(ctx, storageUser)
}

func (calendar *Calendar) SelectUsers(ctx context.Context) ([]model.IUser, error) {
	calendar.mutex.RLock()
	defer calendar.mutex.RUnlock()

	users := make([]model.IUser, 0)

	storageUsers, err := calendar.storage.SelectUsers(ctx)
	if err != nil {
		return users, err
	}

	for _, storageUser := range storageUsers {
		user := storageUser
		users = append(users, &user)
	}

	return users, nil
}

func (calendar *Calendar) DeleteUser(ctx context.Context, id string) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	return calendar.storage.DeleteUser(ctx, id)
}

func (calendar *Calendar) CreateEvent(ctx context.Context, event model.IEvent) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	storageEvent := model.Event{
		ID:           event.GetID(),
		Title:        event.GetTitle(),
		Description:  event.GetDescription(),
		UserID:       event.GetUserID(),
		Beginning:    event.GetBeginning(),
		Finish:       event.GetFinish(),
		Notification: event.GetNotification(),
	}

	return calendar.storage.CreateEvent(ctx, storageEvent)
}

func (calendar *Calendar) UpdateEvent(ctx context.Context, event model.IEvent) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	storageEvent := model.Event{
		ID:           event.GetID(),
		Title:        event.GetTitle(),
		Description:  event.GetDescription(),
		UserID:       event.GetUserID(),
		Beginning:    event.GetBeginning(),
		Finish:       event.GetFinish(),
		Notification: event.GetNotification(),
	}

	return calendar.storage.UpdateEvent(ctx, storageEvent)
}

func (calendar *Calendar) DeleteEvent(ctx context.Context, id string) error {
	calendar.mutex.Lock()
	defer calendar.mutex.Unlock()

	return calendar.storage.DeleteEvent(ctx, id)
}

func (calendar *Calendar) SelectEvents(ctx context.Context) ([]model.IEvent, error) {
	calendar.mutex.RLock()
	defer calendar.mutex.RUnlock()

	events := make([]model.IEvent, 0)

	storageEvents, err := calendar.storage.SelectEvents(ctx)
	if err != nil {
		return events, err
	}

	for _, storageEvent := range storageEvents {
		event := storageEvent
		events = append(events, &event)
	}

	return events, nil
}
