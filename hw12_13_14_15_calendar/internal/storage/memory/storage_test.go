package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/model"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("user case", func(t *testing.T) {
		s := New()
		ctx := context.Background()

		users := []model.User{
			{
				FirstName: "Yuliya",
				LastName:  "Zadorozhnaya",
				Email:     "yzadorozhnaya@mail.ru",
				Age:       22,
			},
			{
				FirstName: "Tom",
				LastName:  "Holland",
				Email:     "tomholland@gmail.ru",
				Age:       25,
			},
		}

		for _, user := range users {
			require.Nil(t, s.CreateUser(ctx, user))
		}

		selectedUsers, err := s.SelectUsers(ctx)
		require.Nil(t, err)

		for _, selectedUser := range selectedUsers {
			require.True(t, containsUser(users, selectedUser))
			require.Nil(t, s.DeleteUser(ctx, selectedUser.ID))
		}

		selectedUsers, err = s.SelectUsers(ctx)
		require.Nil(t, err)
		require.Len(t, selectedUsers, 0)
	})

	t.Run("event case", func(t *testing.T) {
		s := New()
		ctx := context.Background()

		user := model.User{
			FirstName: "Yuliya",
			LastName:  "Zadorozhnaya",
			Email:     "yzadorozhnaya@mail.ru",
			Age:       22,
		}

		require.Nil(t, s.CreateUser(ctx, user))
		selectedUsers, err := s.SelectUsers(ctx)
		require.Nil(t, err)

		user = selectedUsers[0]

		events := []model.Event{
			{
				Title:        "Просмотр фильма",
				Description:  "Человек паук",
				Beginning:    time.Date(2024, time.May, 22, 18, 0, 0, 0, time.UTC),
				Finish:       time.Date(2024, time.May, 23, 0o5, 0, 0, 0, time.UTC),
				Notification: time.Date(2024, time.May, 21, 0o3, 0, 0, 0, time.UTC),
				UserID:       user.ID,
			},
			{
				Title:        "Свадьба",
				Description:  "Свадьба",
				Beginning:    time.Date(2024, time.April, 22, 20, 0, 0, 0, time.UTC),
				Finish:       time.Date(2024, time.April, 22, 22, 30, 0, 0, time.UTC),
				Notification: time.Date(2024, time.April, 22, 23, 0, 0, 0, time.UTC),
				UserID:       user.ID,
			},
		}

		for _, event := range events {
			require.Nil(t, s.CreateEvent(ctx, event))
		}

		selectedEvents, err := s.SelectEvents(ctx)
		require.Nil(t, err)

		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
		}

		events = selectedEvents
		events[0].Description = "Прийти на кв к Жеке, не забыть вкусняшки и заранее заказать пиццу"
		require.Nil(t, s.UpdateEvent(ctx, events[0]))

		selectedEvents, err = s.SelectEvents(ctx)
		require.Nil(t, err)

		for _, selectedEvent := range selectedEvents {
			require.True(t, containsEvent(events, selectedEvent))
			require.Nil(t, s.DeleteEvent(ctx, selectedEvent.ID))
		}

		require.Nil(t, s.DeleteUser(ctx, user.ID))
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
		if event.UserID == e.UserID &&
			event.Finish == e.Finish &&
			event.Notification == e.Notification &&
			event.Beginning == e.Beginning &&
			event.Description == e.Description &&
			event.Title == e.Title {
			return true
		}
	}
	return false
}
