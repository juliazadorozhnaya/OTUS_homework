package internalhttp

import (
	"context"
	"fmt"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/model"
	handler "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server/handler"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server/router"
	"net"
	"net/http"
)

type Server struct {
	app    Application
	logger Logger
	srv    *http.Server
}

type Logger interface {
	Fatal(string)
	Error(string)
	Warn(string)
	Info(string)
	Debug(string)
}

type Application interface {
	CreateUser(ctx context.Context, User model.User) error
	SelectUsers(ctx context.Context) ([]model.User, error)
	DeleteUser(ctx context.Context, id string) error

	CreateEvent(ctx context.Context, Event model.Event) error
	SelectEvents(ctx context.Context) ([]model.Event, error)
	UpdateEvent(ctx context.Context, Event model.Event) error
	DeleteEvent(ctx context.Context, id string) error
}

func NewServer(logger Logger, app Application, config config.ServerConfig) *Server {
	hand := handler.NewHandler(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hand.GetHello)

	middleware := router.NewMiddleware(logger, mux)
	middleware.Logging()

	return &Server{
		logger: logger,
		app:    app,
		srv: &http.Server{
			Addr:    net.JoinHostPort(config.Host, config.Port),
			Handler: middleware.Handler,
		},
	}
}

func (s *Server) Start() error {
	s.logger.Info(fmt.Sprintf("server listening: %s", s.srv.Addr))

	if err := s.srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	if err := s.srv.Close(); err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}
