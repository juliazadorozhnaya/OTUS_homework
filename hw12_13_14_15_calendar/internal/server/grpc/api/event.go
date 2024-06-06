package api

import (
	"context"
	"fmt"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type EventServer struct {
	app    server.Application
	logger server.Logger
}

func NewEventServer(logger server.Logger, app server.Application) *EventServer {
	return &EventServer{
		app:    app,
		logger: logger,
	}
}

func (serv *EventServer) SelectEvents(void *Void, selectEvents EventService_SelectEventsServer) error {
	defer func(start time.Time) {
		duration := time.Since(start)
		serv.Log(selectEvents.Context(), start, duration, "SelectEvents")
	}(time.Now())

	events, err := serv.app.SelectEvents(selectEvents.Context())
	if err != nil {
		return err
	}

	for _, event := range events {
		e := Event{
			ID:            event.GetID(),
			Title:         event.GetTitle(),
			Description:   event.GetDescription(),
			UserID:        event.GetUserID(),
			BeginningT:    timestamppb.New(event.GetBeginning()),
			FinishT:       timestamppb.New(event.GetFinish()),
			NotificationT: timestamppb.New(event.GetNotification()),
		}

		err := selectEvents.Send(&e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (serv *EventServer) CreateEvent(ctx context.Context, event *Event) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		serv.Log(ctx, start, duration, "CreateEvent")
	}(time.Now())

	err := serv.app.CreateEvent(ctx, event)
	return &Void{}, err
}

func (serv *EventServer) UpdateEvent(ctx context.Context, event *Event) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		serv.Log(ctx, start, duration, "UpdateEvent")
	}(time.Now())

	err := serv.app.UpdateEvent(ctx, event)
	return &Void{}, err
}

func (serv *EventServer) DeleteEvent(ctx context.Context, event *Event) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		serv.Log(ctx, start, duration, "DeleteEvent")
	}(time.Now())

	err := serv.app.DeleteEvent(ctx, event.ID)
	return &Void{}, err
}

func (serv *EventServer) mustEmbedUnimplementedEventServiceServer() {}

func (x *Event) GetBeginning() time.Time {
	return x.BeginningT.AsTime()
}

func (x *Event) GetFinish() time.Time {
	return x.FinishT.AsTime()
}

func (x *Event) GetNotification() time.Time {
	return x.NotificationT.AsTime()
}

func (serv *EventServer) Log(ctx context.Context, start time.Time, duration time.Duration, funcName string) {
	ip := ""

	if p, ok := peer.FromContext(ctx); ok {
		ip = p.Addr.String()
	}

	logMessage := fmt.Sprintf("%s [%s] %s %s",
		ip, start, funcName, duration)
	serv.logger.Info(logMessage)
}
