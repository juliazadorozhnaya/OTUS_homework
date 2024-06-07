package api

import (
	"context"
	"time"

	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	UnimplementedUserServiceServer
	app    server.Application
	logger server.Logger
}

func NewUserServer(logger server.Logger, app server.Application) *UserServer {
	return &UserServer{
		app:    app,
		logger: logger,
	}
}

func (s *UserServer) SelectUsers(ctx context.Context, _ *Void) (*Users, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("SelectUsers", ctx, start, duration)
	}(time.Now())

	users, err := s.app.SelectUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select users: %v", err)
	}

	protoUsers := make([]*User, len(users))
	for i, user := range users {
		protoUsers[i] = &User{
			ID:        user.GetID(),
			FirstName: user.GetFirstName(),
			LastName:  user.GetLastName(),
			Email:     user.GetEmail(),
			Age:       user.GetAge(),
		}
	}

	return &Users{Users: protoUsers}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, user *User) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("CreateUser", ctx, start, duration)
	}(time.Now())

	err := s.app.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}
	return &Void{}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, user *User) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("DeleteUser", ctx, start, duration)
	}(time.Now())

	err := s.app.DeleteUser(ctx, user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}
	return &Void{}, nil
}

func (s *UserServer) mustEmbedUnimplementedUserServiceServer() {}
