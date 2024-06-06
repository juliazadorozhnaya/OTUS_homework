package server

import (
	"context"

	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/model"
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
}
