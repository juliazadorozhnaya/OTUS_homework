package servergrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/server"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/server/grpc/api"
	"google.golang.org/grpc"
)

type Server struct {
	app     server.Application
	logger  server.Logger
	address string
	srv     *grpc.Server
}

func NewServer(logger server.Logger, app server.Application, config server.Config) *Server {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor(logger)),
	)

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

// Start запускает gRPC сервер.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		s.logger.Fatal(fmt.Sprintf("Failed to start gRPC server: %s", err))
		return err
	}

	s.logger.Info(fmt.Sprintf("gRPC server listening: %s", s.address))
	if err := s.srv.Serve(listener); err != nil {
		s.logger.Fatal(fmt.Sprintf("gRPC server failed to serve: %s", err))
		return err
	}

	s.logger.Debug("gRPC server started successfully")
	return nil
}

// Stop останавливает gRPC сервер.
func (s *Server) Stop(_ context.Context) error {
	s.logger.Info("gRPC server shutting down...")

	s.srv.GracefulStop()
	s.logger.Debug("gRPC server stopped gracefully")

	return nil
}
