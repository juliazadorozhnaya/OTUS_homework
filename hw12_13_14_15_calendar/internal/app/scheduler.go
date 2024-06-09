package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/broker"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/model"
)

// Scheduler отвечает за периодическое сканирование базы данных,
// отправку уведомлений через RabbitMQ и очистку старых событий.
type Scheduler struct {
	app      *Calendar
	broker   broker.Broker
	logger   *logger.Logger
	interval time.Duration
	stopChan chan struct{}
}

func NewScheduler(app *Calendar, broker broker.Broker, logger *logger.Logger, interval time.Duration) *Scheduler {
	return &Scheduler{
		app:      app,
		broker:   broker,
		logger:   logger,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// Start запускает процесс планировщика, который выполняется с заданным интервалом.
func (s *Scheduler) Start(ctx context.Context) error {
	s.logger.Info("Scheduler started")
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.logger.Info("Scheduler tick")
			s.HandleNotifications(ctx)
			s.cleanupOldEvents(ctx)

		case <-ctx.Done():
			s.logger.Info("Scheduler stopped")
			return nil

		case <-s.stopChan:
			s.logger.Info("Scheduler stop signal received")
			return nil
		}
	}
}

// HandleNotifications обрабатывает уведомления, выбирая события из базы данных
// и отправляя их через RabbitMQ.
func (s *Scheduler) HandleNotifications(ctx context.Context) {
	now := time.Now().Truncate(time.Second)
	formattedTime := now.Format("2006-01-02 15:04:05")
	t, err := time.Parse("2006-01-02 15:04:05", formattedTime)
	if err != nil {
		s.logger.Error("Error formatting time: %v", err)
		return
	}

	events, err := s.app.SelectEventsByTime(ctx, t)
	if err != nil {
		s.logger.Error("Error getting upcoming events: %v", err)
		return
	}

	for _, event := range events {
		notify := model.Event{
			ID:           event.GetID(),
			Title:        event.GetTitle(),
			Description:  event.GetDescription(),
			Notification: event.GetNotification(),
		}
		body, err := json.Marshal(notify)
		if err != nil {
			s.logger.Error("Error marshalling notification: %v", err)
			continue
		}

		if err := s.broker.PublishWithContext(ctx, *config.Get().RabbitMQ.Publish, body); err != nil {
			s.logger.Error("Error publishing message: %v", err)
		}
	}
}

// cleanupOldEvents удаляет события, которые произошли более года назад.
func (s *Scheduler) cleanupOldEvents(ctx context.Context) {
	oldEvents, err := s.app.SelectEventsForMonth(ctx, time.Now().AddDate(-1, 0, 0))
	if err != nil {
		s.logger.Error("Error selecting old events: %v", err)
		return
	}
	for _, event := range oldEvents {
		if err := s.app.DeleteEvent(ctx, event.GetID()); err != nil {
			s.logger.Error("Error deleting old event: %v", err)
		}
	}
}

// Stop останавливает планировщик, посылая сигнал остановки.
func (s *Scheduler) Stop() {
	s.logger.Info("Scheduler stopping...")
	close(s.stopChan)
}
