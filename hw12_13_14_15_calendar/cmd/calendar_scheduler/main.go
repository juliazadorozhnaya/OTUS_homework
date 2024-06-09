package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/broker/rabbitmq"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/juliazadorozhnaya/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/pressly/goose/v3"
)

var (
	configPath  string
	storageType string

	ErrorInvalidStorageType = errors.New("invalid storage type")
)

func init() {
	defaultConfigPath := path.Join("config", "scheduler_config.toml")
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to configuration file")
	flag.StringVar(&storageType, "storage", "sql", "Type of storage. Expected values: \"memory\" || \"sql\"")
}

func main() {
	flag.Parse()

	log.Println("Loading configuration...")
	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	conf := config.Get()
	l := logger.New(conf.Logger)

	var storage app.Storage
	switch storageType {
	case "memory":
		storage = memorystorage.New()
	case "sql":
		dbConn := conf.Database
		connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			dbConn.UserName, dbConn.Password, dbConn.Host, dbConn.Port, dbConn.DatabaseName)

		l.Info("Running migrations...")
		if err := runMigrations(connString); err != nil {
			l.Error("Failed to run migrations: " + err.Error())
			return
		}

		var err error
		storage, err = sqlstorage.New(connString)
		if err != nil {
			l.Error("Failed to create SQL storage: " + err.Error())
			return
		}
	default:
		l.Error(ErrorInvalidStorageType.Error())
		os.Exit(1)
	}

	l.Info("Connecting to RabbitMQ...")
	rabbit := rabbitmq.New(*conf.RabbitMQ.Connection)
	if err := rabbit.Start(); err != nil {
		l.Error("Error connecting to RabbitMQ: " + err.Error())
		return
	}
	defer func() {
		l.Info("Stopping RabbitMQ connection...")
		if err := rabbit.Stop(); err != nil {
			l.Error("Error stopping RabbitMQ: " + err.Error())
		}
	}()

	l.Info("Declaring RabbitMQ queue...")
	err := rabbit.QueueDeclare(*conf.RabbitMQ.Queue)
	if err != nil {
		l.Error("Error declaring queue: " + err.Error())
		return
	}

	l.Info("Creating new calendar app...")
	calendarApp := app.New(storage)
	scheduler := app.NewScheduler(calendarApp, &rabbit, l, conf.RabbitMQ.Consume.Interval)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		l.Info("Received interrupt signal, shutting down...")
		cancel()
	}()

	l.Info("Starting scheduler...")
	if err := scheduler.Start(ctx); err != nil {
		l.Error("Scheduler error: " + err.Error())
	}
}

func runMigrations(connString string) error {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.Up(db, "migrations", goose.WithAllowMissing()); err != nil {
		return err
	}
	return nil
}
