package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"path"

	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/broker/rabbitmq"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/logger"
)

var (
	configPath  string
	storageType string

	ErrorInvalidStorageType = errors.New("invalid storage type")
)

func init() {
	defaultConfigPath := path.Join("config", "sender_config.toml")
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to configuration file")
	flag.StringVar(&storageType, "storage", "sql", "Type of storage. Expected values: \"memory\" || \"sql\"")
}

func main() {
	flag.Parse()

	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	conf := config.Get()
	l := logger.New(conf.Logger)

	l.Info("Connecting to RabbitMQ...")
	rabbit := rabbitmq.New(*conf.RabbitMQ.Connection)
	if err := rabbit.Start(); err != nil {
		l.Fatal("Error connecting to RabbitMQ: %v", err)
		return
	}
	defer func() {
		l.Info("Stopping RabbitMQ connection...")
		if err := rabbit.Stop(); err != nil {
			l.Error("Error stopping RabbitMQ: %v", err)
		}
	}()

	l.Info("Starting to consume messages from RabbitMQ...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sender := app.NewSender(&rabbit, l)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		l.Info("Received interrupt signal, shutting down...")
		cancel()
		sender.Stop()
	}()

	l.Info("Starting sender...")
	go func() {
		if err := sender.Start(ctx); err != nil {
			l.Fatal("Error starting sender: %v", err)
		}
	}()

	l.Info("Waiting for messages...")

	<-ctx.Done()
	l.Info("Context cancelled, exiting main function...")
}
