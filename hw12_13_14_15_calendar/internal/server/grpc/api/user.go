package api

import (
	"context"
	"fmt"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc/peer"
	"time"
)

type UserServer struct {
	app    server.Application
	logger server.Logger
}

func NewUserServer(logger server.Logger, app server.Application) *UserServer {
	return &UserServer{
		app:    app,
		logger: logger,
	}
}

func (serv *UserServer) SelectUsers(void *Void, selectUsers UserService_SelectUsersServer) error {
	defer func(start time.Time) {
		duration := time.Since(start)
		serv.Log(selectUsers.Context(), start, duration, "SelectUsers")
	}(time.Now())

	users, err := serv.app.SelectUsers(selectUsers.Context())
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

		err := selectUsers.Send(&e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (serv *UserServer) CreateUser(ctx context.Context, user *User) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		serv.Log(ctx, start, duration, "CreateUser")
	}(time.Now())

	err := serv.app.CreateUser(ctx, user)
	return &Void{}, err
}

func (serv *UserServer) DeleteUser(ctx context.Context, user *User) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		serv.Log(ctx, start, duration, "DeleteUser")
	}(time.Now())

	err := serv.app.DeleteUser(ctx, user.ID)
	return &Void{}, err
}

func (serv *UserServer) mustEmbedUnimplementedUserServiceServer() {}

func (serv *UserServer) Log(ctx context.Context, start time.Time, duration time.Duration, funcName string) {
	ip := ""

	if p, ok := peer.FromContext(ctx); ok {
		ip = p.Addr.String()
	}

	logMessage := fmt.Sprintf("%s [%s] %s %s",
		ip, start, funcName, duration)
	serv.logger.Info(logMessage)
}
