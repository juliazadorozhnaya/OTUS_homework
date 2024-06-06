package serverhttp

import (
	"fmt"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
	"net"
	"net/http"
)

type Server struct {
	app    server.Application
	logger server.Logger
	srv    *http.Server
}

func NewServer(logger server.Logger, app server.Application, config server.Config) *Server {
	handler := newHandler(logger, app)

	mux := http.NewServeMux()
	mux.HandleFunc("/create/user", handler.createUser)
	mux.HandleFunc("/select/users", handler.selectUsers)
	mux.HandleFunc("/delete/user", handler.deleteUser)
	mux.HandleFunc("/create/event", handler.createEvent)
	mux.HandleFunc("/select/events", handler.selectEvents)
	mux.HandleFunc("/update/event", handler.updateEvent)
	mux.HandleFunc("/delete/event", handler.deleteEvent)

	middleWare := newMiddleware(logger, mux)
	middleWare.logging()

	return &Server{
		logger: logger,
		app:    app,
		srv: &http.Server{
			Addr:    net.JoinHostPort(config.GetHost(), config.GetPort()),
			Handler: middleWare.Handler,
		},
	}
}

func (s *Server) Start() error {
	s.logger.Info(fmt.Sprintf("http server listening: %s", s.srv.Addr))

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
