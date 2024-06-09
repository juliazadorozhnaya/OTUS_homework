package integration_test

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/model"
	sqlstorage "github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	mutex := sync.Mutex{}
	dsn := "postgresql://postgres:1234512345@postgres:5432/calendardb?sslmode=disable"

	t.Run("user case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		s, err := sqlstorage.New(dsn)
		if err != nil {
			log.Fatalf("failed to create storage: %v", err)
		}
		ctx := context.Background()

		users := []model.User{
			{
				FirstName: "Alice",
				LastName:  "Johnson",
				Email:     "alice.johnson@example.com",
				Age:       30,
			},
			{
				FirstName: "Bob",
				LastName:  "Smith",
				Email:     "bob.smith@example.com",
				Age:       25,
			},
		}

		for _, user := range users {
			require.Nil(t, s.CreateUser(ctx, user))
		}

		selectedUsers, err := s.SelectUsers(ctx)
		require.Nil(t, err)

		log.Printf("Selected Users: %+v", selectedUsers)

		for _, selectedUser := range selectedUsers {
			log.Printf("Checking user: %+v", selectedUser)
			require.True(t, containsUser(users, selectedUser))
			require.Nil(t, s.DeleteUser(ctx, selectedUser.ID))
		}

		selectedUsers, err = s.SelectUsers(ctx)
		require.Nil(t, err)
		require.Len(t, selectedUsers, 0)
	})

	t.Run("event case", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()

		s, err := sqlstorage.New(dsn)
		if err != nil {
			log.Fatalf("failed to create storage: %v", err)
		}
		ctx := context.Background()

		user := model.User{
			FirstName: "Charlie",
			LastName:  "Brown",
			Email:     "charlie.brown@example.com",
			Age:       35,
		}

		require.Nil(t, s.CreateUser(ctx, user))
		selectedUsers, err := s.SelectUsers(ctx)
		require.Nil(t, err)

		user = selectedUsers[0]

		events := []model.Event{
			{
				Title:        "Meeting with team",
				Description:  "",
				Beginning:    time.Date(2023, time.July, 15, 9, 0, 0, 0, time.UTC),
				Finish:       time.Date(2023, time.July, 15, 10, 0, 0, 0, time.UTC),
				Notification: time.Date(2023, time.July, 15, 8, 30, 0, 0, time.UTC),
				UserID:       user.ID,
			},
			{
				Title:        "Project presentation",
				Description:  "Prepare slides and speech",
				Beginning:    time.Date(2023, time.July, 15, 11, 0, 0, 0, time.UTC),
				Finish:       time.Date(2023, time.July, 15, 12, 0, 0, 0, time.UTC),
				Notification: time.Date(2023, time.July, 15, 10, 30, 0, 0, time.UTC),
				UserID:       user.ID,
			},
		}

		for _, event := range events {
			require.Nil(t, s.CreateEvent(ctx, event))
		}

		selectedEvents, err := s.SelectEventsByTime(ctx, events[0].Notification)
		require.Nil(t, err)

		log.Printf("Selected Events by Time: %+v", selectedEvents)

		for _, selectedEvent := range selectedEvents {
			log.Printf("Checking event: %+v", selectedEvent)
			require.True(t, containsEvent(events, selectedEvent))
		}

		selectedEvents, err = s.SelectEvents(ctx)
		require.Nil(t, err)

		log.Printf("Selected Events: %+v", selectedEvents)

		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
		}

		events = selectedEvents
		events[0].Description = "Come to the meeting room early"
		require.Nil(t, s.UpdateEvent(ctx, events[0]))

		selectedEvents, err = s.SelectEvents(ctx)
		require.Nil(t, err)

		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
			require.Nil(t, s.DeleteEvent(ctx, selectedEvent.ID))
		}

		selectedEvents, err = s.SelectEvents(ctx)
		require.Nil(t, err)
		require.Len(t, selectedEvents, 0)

		require.Nil(t, s.DeleteUser(ctx, user.ID))
		selectedUsers, err = s.SelectUsers(ctx)
		require.Nil(t, err)
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
