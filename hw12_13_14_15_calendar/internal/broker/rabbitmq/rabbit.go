package rabbitmq

import (
	"context"
	"fmt"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type BrokerRabbit struct {
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(connConfig config.ConnectionConfig) BrokerRabbit {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s",
		connConfig.Login, connConfig.Password, connConfig.Host, connConfig.Port)
	return BrokerRabbit{url: url}
}

// Start подключается к RabbitMQ и открывает канал.
func (b *BrokerRabbit) Start() error {
	var err error
	b.conn, err = amqp.Dial(b.url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	b.ch, err = b.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	return nil
}

// QueueDeclare создает (декларирует) очередь в RabbitMQ.
func (b *BrokerRabbit) QueueDeclare(config config.QueueConfig) error {
	_, err := b.ch.QueueDeclare(
		config.Name,
		config.Durable,
		config.AutoDelete,
		config.Exclusive,
		config.NoWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	return nil
}

// Stop закрывает канал и соединение с RabbitMQ.
func (b *BrokerRabbit) Stop() error {
	if err := b.ch.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	if err := b.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection to RabbitMQ: %w", err)
	}

	return nil
}

// Consume регистрирует потребителя для указанной очереди и возвращает канал для получения сообщений.
func (b *BrokerRabbit) Consume(config config.ConsumeConfig) (<-chan amqp.Delivery, error) {
	delivery, err := b.ch.Consume(
		config.Queue,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	return delivery, nil
}

// PublishWithContext публикует сообщение в RabbitMQ с использованием контекста.
func (b *BrokerRabbit) PublishWithContext(ctx context.Context, config config.PublishConfig, body []byte) error {
	err := b.ch.PublishWithContext(ctx,
		config.Exchange,
		config.Key,
		config.Mandatory,
		config.Immediate,
		amqp.Publishing{
			ContentType: config.ContentType,
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
