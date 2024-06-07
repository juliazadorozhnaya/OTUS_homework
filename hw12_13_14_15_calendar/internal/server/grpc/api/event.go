package api

import (
	"context"
	"time"

	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type EventServer struct {
	UnimplementedEventServiceServer
	logger server.Logger
	app    server.Application
}

func NewEventServer(logger server.Logger, app server.Application) *EventServer {
	return &EventServer{logger: logger, app: app}
}

func (s *EventServer) SelectEvents(ctx context.Context, _ *Void) (*Events, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("SelectEvents", ctx, start, duration)
	}(time.Now())

	events, err := s.app.SelectEvents(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select events: %v", err)
	}

	var result Events
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
		result.Events = append(result.Events, &e)
	}

	return &result, nil
}

func (s *EventServer) CreateEvent(ctx context.Context, event *Event) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("CreateEvent", ctx, start, duration)
	}(time.Now())

	err := s.app.CreateEvent(ctx, event)
	return &Void{}, err
}

func (s *EventServer) UpdateEvent(ctx context.Context, event *Event) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("UpdateEvent", ctx, start, duration)
	}(time.Now())

	if err := s.app.UpdateEvent(ctx, event); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}
	return &Void{}, nil
}

func (s *EventServer) DeleteEvent(ctx context.Context, event *Event) (*Void, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("DeleteEvent", ctx, start, duration)
	}(time.Now())

	if err := s.app.DeleteEvent(ctx, event.ID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}
	return &Void{}, nil
}

func (s *EventServer) SelectEventsForDay(ctx context.Context, req *DateRequest) (*Events, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("SelectEventsForDay", ctx, start, duration)
	}(time.Now())

	events, err := s.app.SelectEventsForDay(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select events for day: %v", err)
	}

	var result Events
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
		result.Events = append(result.Events, &e)
	}

	return &result, nil
}

func (s *EventServer) SelectEventsForWeek(ctx context.Context, req *DateRequest) (*Events, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("SelectEventsForWeek", ctx, start, duration)
	}(time.Now())

	events, err := s.app.SelectEventsForWeek(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select events for week: %v", err)
	}

	var result Events
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
		result.Events = append(result.Events, &e)
	}

	return &result, nil
}

func (s *EventServer) SelectEventsForMonth(ctx context.Context, req *DateRequest) (*Events, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Info("SelectEventsForMonth", ctx, start, duration)
	}(time.Now())

	events, err := s.app.SelectEventsForMonth(ctx, req.Date.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to select events for month: %v", err)
	}

	var result Events
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
		result.Events = append(result.Events, &e)
	}

	return &result, nil
}

func (s *EventServer) mustEmbedUnimplementedEventServiceServer() {}

func (x *Event) GetBeginning() time.Time {
	return x.BeginningT.AsTime()
}

func (x *Event) GetFinish() time.Time {
	return x.FinishT.AsTime()
}

func (x *Event) GetNotification() time.Time {
	return x.NotificationT.AsTime()
}
