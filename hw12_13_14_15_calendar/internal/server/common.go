package server

import (
	"context"
	"time"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/model"
)

type Server interface {
	Start() error
	Stop() error
}

type Config interface {
	GetPort() string
	GetHost() string
}

type Logger interface {
	Fatal(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
}

type Application interface {
	CreateUser(context.Context, model.IUser) error
	SelectUsers(context.Context) ([]model.IUser, error)
	DeleteUser(context.Context, string) error

	CreateEvent(context.Context, model.IEvent) error
	SelectEvents(context.Context) ([]model.IEvent, error)
	UpdateEvent(context.Context, model.IEvent) error
	DeleteEvent(context.Context, string) error

	SelectEventsForDay(context.Context, time.Time) ([]model.IEvent, error)
	SelectEventsForWeek(context.Context, time.Time) ([]model.IEvent, error)
	SelectEventsForMonth(context.Context, time.Time) ([]model.IEvent, error)
}
