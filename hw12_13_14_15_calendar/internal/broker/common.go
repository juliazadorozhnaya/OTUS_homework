package broker

import (
	"context"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker interface {
	Start() error
	Stop() error
	QueueDeclare(config config.QueueConfig) error
	Consume(config config.ConsumeConfig) (<-chan amqp.Delivery, error)
	PublishWithContext(ctx context.Context, config config.PublishConfig, body []byte) error
}
