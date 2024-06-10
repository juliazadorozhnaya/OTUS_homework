package integration_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/broker/rabbitmq"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/model"
	memorystorage "github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func setupRabbitMQ(t *testing.T, log *logger.Logger) *rabbitmq.BrokerRabbit {
	t.Helper()

	log.Info("Setting up RabbitMQ...")
	rabbitConfig := config.ConnectionConfig{
		Login:    "guest",
		Password: "guest",
		Host:     "localhost",
		Port:     "5672",
	}

	broker := rabbitmq.New(rabbitConfig)
	err := broker.Start()
	require.NoError(t, err)
	log.Info("RabbitMQ started")

	queueConfig := config.QueueConfig{
		Name:       "test_queue",
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
	}
	err = broker.QueueDeclare(queueConfig)
	require.NoError(t, err)
	log.Info("Queue declared: %v", queueConfig)

	require.NoError(t, err)

	return &broker
}

func TestScheduler(t *testing.T) {
	err := config.LoadConfig("../config/sender_config.toml")
	require.NoError(t, err)
	t.Log("Config loaded")

	conf := config.Get()
	log := logger.New(conf.Logger)
	t.Log("Logger set up")

	memoryStorage := memorystorage.New()
	application := app.New(memoryStorage, *log)
	t.Log("Memory storage and application set up")

	broker := setupRabbitMQ(t, log)
	defer func() {
		log.Info("Stopping RabbitMQ...")
		err := broker.Stop()
		require.NoError(t, err)
		log.Info("RabbitMQ stopped")
	}()
	t.Log("RabbitMQ set up")

	scheduler := app.NewScheduler(application, broker, log, 1*time.Second)
	t.Log("Scheduler created")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t.Log("Context created")

	go func() {
		log.Info("Starting Scheduler...")
		if err := scheduler.Start(ctx); err != nil {
			log.Error("Scheduler error: %v", err)
		}
	}()
	defer func() {
		log.Info("Stopping Scheduler...")
		scheduler.Stop()
		log.Info("Scheduler stopped")
	}()
	t.Log("Scheduler started")

	user := &model.User{
		FirstName: "testuser",
		LastName:  "last",
		Email:     "test@test.com",
		Age:       30,
	}
	err = application.CreateUser(ctx, user)
	require.NoError(t, err)
	t.Logf("User created: %v", user)

	users, err := application.SelectUsers(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	userID := users[0].GetID()
	t.Logf("Users selected, userID: %v", userID)

	now := time.Now().UTC().Truncate(time.Second)
	event := &model.Event{
		Title:        "notifyevent",
		Description:  "desc",
		Beginning:    now.Add(6 * time.Hour),
		Finish:       now.Add(10 * time.Hour),
		Notification: now.Add(3*time.Hour + 5*time.Second),
		UserID:       userID,
	}
	err = application.CreateEvent(ctx, event)
	require.NoError(t, err)
	log.Info("Event created: %v", event)

	time.Sleep(10 * time.Second)
	newEvent, _ := application.SelectEventsByTime(ctx, event.Notification)

	body, _ := json.Marshal(newEvent)
	publishConfig := config.PublishConfig{
		Exchange:    "test_exchange",
		Key:         "special_consumer_key",
		Mandatory:   false,
		Immediate:   false,
		ContentType: "application/json",
	}
	err = broker.PublishWithContext(ctx, publishConfig, body)
	require.NoError(t, err)

	msgs, err := broker.Consume(config.ConsumeConfig{
		Queue:     "test_queue",
		Consumer:  "special_consumer_key",
		AutoAck:   true,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
	})
	if err != nil {
		log.Error("Failed to start consuming messages: %v", err)
	}
	require.NotNil(t, msgs, "Message channel should not be nil")

	select {
	case msg, ok := <-msgs:
		if !ok {
			t.Fatal("Message channel closed")
		}
		t.Log("Received message:", string(msg.Body))
	case <-time.After(30 * time.Second):
		t.Fatal("Did not receive notification in time")
	}
}
