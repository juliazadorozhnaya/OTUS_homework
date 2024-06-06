package api

import (
	"context"
	"time"

	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
)

type UserServer struct {
	UnimplementedEventServiceServer
	app    server.Application
	logger server.Logger
}

func NewUserServer(logger server.Logger, app server.Application) *UserServer {
	return &UserServer{
		app:    app,
		logger: logger,
	}
}

func (s *UserServer) SelectUsers(_ *Void, stream UserService_SelectUsersServer) error {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("SelectUsers", stream.Context(), start, duration) //nolint: staticcheck
	}(time.Now())

	users, err := s.app.SelectUsers(stream.Context()) //nolint: staticcheck
	if err != nil {
		return err
	}

	for _, user := range users {
		e := User{
			ID:        user.GetID(),
			FirstName: user.GetFirstName(),
			LastName:  user.GetLastName(),
			Email:     user.GetEmail(),
			Age:       user.GetAge(),
		}

		err := stream.Send(&e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserServer) CreateUser(ctx context.Context, user *User) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("CreateUser", ctx, start, duration)
	}(time.Now())

	err := s.app.CreateUser(ctx, user)
	return &Void{}, err
}

func (s *UserServer) DeleteUser(ctx context.Context, user *User) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("DeleteUser", ctx, start, duration)
	}(time.Now())

	err := s.app.DeleteUser(ctx, user.ID)
	return &Void{}, err
}

func (s *UserServer) mustEmbedUnimplementedUserServiceServer() {}
