package servergrpc

import (
	"fmt"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server/grpc/api"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	app     server.Application
	logger  server.Logger
	address string
	srv     *grpc.Server
}

func NewServer(logger server.Logger, app server.Application, config server.Config) *Server {
	srv := grpc.NewServer()

	eventServer := api.NewEventServer(logger, app)
	api.RegisterEventServiceServer(srv, eventServer)

	userServer := api.NewUserServer(logger, app)
	api.RegisterUserServiceServer(srv, userServer)

	return &Server{
		logger:  logger,
		app:     app,
		srv:     srv,
		address: net.JoinHostPort(config.GetHost(), config.GetPort()),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	s.logger.Info(fmt.Sprintf("grpc server listening: %s", s.address))
	if err := s.srv.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.srv.Stop()
	return nil
}
