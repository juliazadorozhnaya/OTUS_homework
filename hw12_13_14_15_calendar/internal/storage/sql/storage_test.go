package sqlstorage

import (
	"context"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

var (
	dsn = "postgresql://postgres:1234512345@postgres:/postgres?sslmode=disable"
)

func TestStorage(t *testing.T) {
	mutex := sync.Mutex{}
	s := New(dsn)

	t.Run("user case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		ctx := context.Background()

		users := []model.User{
			{
				FirstName: "Alice",
				LastName:  "Johnson",
				Email:     "alice.johnson@example.com",
				Age:       29,
			},
			{
				FirstName: "Bob",
				LastName:  "Smith",
				Email:     "bob.smith@example.com",
				Age:       34,
			},
		}

		for _, user := range users {
			require.NoError(t, s.CreateUser(ctx, user))
		}

		selectedUsers, err := s.SelectUsers(ctx)
		require.NoError(t, err)

		for _, selectedUser := range selectedUsers {
			require.True(t, containsUser(users, selectedUser))
			require.NoError(t, s.DeleteUser(ctx, selectedUser.ID))
		}

		selectedUsers, err = s.SelectUsers(ctx)
		require.NoError(t, err)
		require.Len(t, selectedUsers, 0)
	})

	t.Run("event case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		ctx := context.Background()

		user := model.User{
			FirstName: "Alice",
			LastName:  "Johnson",
			Email:     "alice.johnson@example.com",
			Age:       29,
		}

		require.NoError(t, s.CreateUser(ctx, user))
		selectedUsers, err := s.SelectUsers(ctx)
		require.NoError(t, err)

		user = selectedUsers[0]

		events := []model.Event{
			{
				Title:        "Project Meeting",
				Description:  "Discuss project requirements",
				Beginning:    time.Date(2023, time.June, 22, 10, 0, 0, 0, time.UTC),
				Finish:       time.Date(2023, time.June, 22, 11, 0, 0, 0, time.UTC),
				Notification: time.Date(2023, time.June, 22, 9, 0, 0, 0, time.UTC),
				UserID:       user.ID,
			},
			{
				Title:        "Team Building",
				Description:  "Outdoor team-building activities",
				Beginning:    time.Date(2023, time.June, 22, 14, 0, 0, 0, time.UTC),
				Finish:       time.Date(2023, time.June, 22, 17, 0, 0, 0, time.UTC),
				Notification: time.Date(2023, time.June, 22, 13, 0, 0, 0, time.UTC),
				UserID:       user.ID,
			},
		}

		for _, event := range events {
			require.NoError(t, s.CreateEvent(ctx, event))
		}

		selectedEvents, err := s.SelectEvents(ctx)
		require.NoError(t, err)

		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
		}

		events = selectedEvents
		events[0].Description = "Discuss updated project requirements"
		require.NoError(t, s.UpdateEvent(ctx, events[0]))

		selectedEvents, err = s.SelectEvents(ctx)
		require.NoError(t, err)

		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
			require.NoError(t, s.DeleteEvent(ctx, selectedEvent.ID))
		}

		selectedEvents, err = s.SelectEvents(ctx)
		require.NoError(t, err)
		require.Len(t, selectedEvents, 0)

		require.NoError(t, s.DeleteUser(ctx, user.ID))
		selectedUsers, err = s.SelectUsers(ctx)
		require.NoError(t, err)
		require.Len(t, selectedUsers, 0)
	})
}

func containsUser(users []model.User, u model.User) bool {
	for _, user := range users {
		if user.FirstName == u.FirstName &&
			user.LastName == u.LastName &&
			user.Email == u.Email &&
			user.Age == u.Age {
			return true
		}
	}
	return false
}

func containsEvent(events []model.Event, e model.Event) bool {
	for _, event := range events {
		if event.Finish == e.Finish &&
			event.Notification == e.Notification &&
			event.Beginning == e.Beginning &&
			event.Description == e.Description &&
			event.Title == e.Title {
			return true
		}
	}
	return false
}
