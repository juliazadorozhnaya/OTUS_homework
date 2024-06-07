package serverhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
)

// Server представляет HTTP сервер с логированием и обработкой запросов.
type Server struct {
	app    server.Application
	logger server.Logger
	srv    *http.Server
}

// NewServer создает новый HTTP сервер с указанным логгером, приложением и конфигурацией.
func NewServer(logger server.Logger, app server.Application, config server.Config) *Server {
	handler := newHandler(logger, app)

	mux := http.NewServeMux()
	mux.HandleFunc("/create/user", handler.createUser)
	mux.HandleFunc("/select/users", handler.selectUsers)
	mux.HandleFunc("/delete/user/", handler.deleteUser)
	mux.HandleFunc("/create/event", handler.createEvent)

	mux.HandleFunc("/select/events", handler.selectEvents)
	mux.HandleFunc("/update/event", handler.updateEvent)
	mux.HandleFunc("/delete/event/", handler.deleteEvent)

	mux.HandleFunc("/select/events/day", handler.selectEventsForDay)
	mux.HandleFunc("/select/events/week", handler.selectEventsForWeek)
	mux.HandleFunc("/select/events/month", handler.selectEventsForMonth)

	middleWare := newMiddleware(logger, mux).logging()

	return &Server{
		logger: logger,
		app:    app,
		srv: &http.Server{
			Addr:              net.JoinHostPort(config.GetHost(), config.GetPort()),
			Handler:           middleWare.Handler,
			ReadHeaderTimeout: 10 * time.Second,
		},
	}
}

// Start запускает HTTP сервер.
func (s *Server) Start() error {
	s.logger.Info(fmt.Sprintf("HTTP server listening: %s", s.srv.Addr))

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal(fmt.Sprintf("HTTP server failed to start: %s", err))
		return err
	}

	s.logger.Debug("HTTP server started successfully")
	return nil
}

// Stop останавливает HTTP сервер.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("HTTP server shutting down...")

	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.Error(fmt.Sprintf("HTTP server shutdown failed: %s", err))
		return err
	}

	s.logger.Debug("HTTP server stopped successfully")
	s.logger.Info("HTTP server stopped")
	return nil
}
