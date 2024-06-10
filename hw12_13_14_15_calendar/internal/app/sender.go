package app

import (
	"context"
	"fmt"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/broker"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/logger"
)

// Sender отвечает за чтение сообщений из очереди RabbitMQ
// и логирование их.
type Sender struct {
	broker   broker.Broker
	logger   *logger.Logger
	stopChan chan struct{}
}

func NewSender(broker broker.Broker, logger *logger.Logger) *Sender {
	return &Sender{
		broker:   broker,
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

// Start запускает процесс рассыльщика, который читает сообщения из очереди RabbitMQ.
func (s *Sender) Start(ctx context.Context) error {
	s.logger.Info("Sender started")

	if s.broker == nil {
		return fmt.Errorf("broker is not initialized")
	}

	msgs, err := s.broker.Consume(*config.Get().RabbitMQ.Consume)
	if err != nil {
		return fmt.Errorf("failed to start consuming messages: %w", err)
	}

	for {
		select {
		case msg := <-msgs:
			s.logger.Info("Received message: %s", msg.Body)

			if !config.Get().RabbitMQ.Consume.AutoAck {
				if err := msg.Ack(false); err != nil {
					s.logger.Error("Failed to acknowledge message: %s", err)
				}
			}

		case <-ctx.Done():
			s.logger.Info("Sender stopped")
			return nil

		case <-s.stopChan:
			s.logger.Info("Sender stop signal received")
			return nil
		}
	}
}

// Stop останавливает рассыльщик, посылая сигнал остановки.
func (s *Sender) Stop() {
	s.logger.Info("Sender stopping...")
	close(s.stopChan)
}
