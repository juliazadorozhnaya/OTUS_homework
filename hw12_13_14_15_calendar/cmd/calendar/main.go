package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/app"
	_ "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server"
	servergrpc "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server/grpc"
	serverhttp "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/storage/sql"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"

	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/logger"
)

var (
	configPath  string
	storageType string

	ErrorInvalidStorageType = errors.New("invalid storage type")
)

func init() {
	defaultConfigPath := path.Join("configs", "config.toml")
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to configuration file")

	flag.StringVar(&storageType, "storage", "sql", "Type of storage. Expected values: \"mem\" || \"sql\"")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf, err := config.New(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	log := logger.New(conf.Logger)

	var application *app.Calendar
	switch storageType {
	case "memory":
		storage := memorystorage.New()
		application = app.New(storage)
	case "sql":
		dbConn := conf.Database
		connString := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
			dbConn.Prefix, dbConn.UserName, dbConn.Password, dbConn.Host, dbConn.Port, dbConn.DatabaseName)
		storage := sqlstorage.New(connString)
		application = app.New(storage)
	default:
		log.Error(ErrorInvalidStorageType.Error())
		os.Exit(1)
	}

	httpServer := serverhttp.NewServer(log, application, conf.HTTPServer)
	grpcServer := servergrpc.NewServer(log, application, conf.GRPCServer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := httpServer.Start(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			log.Error("failed to start HTTP server: " + err.Error())
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		if err := grpcServer.Start(); err != nil {
			log.Error("failed to start gRPC server: " + err.Error())
			cancel()
		}
	}()

	go func() {
		<-ctx.Done()
		log.Info("shutting down servers...")

		if err := httpServer.Stop(); err != nil {
			log.Error("failed to stop HTTP server: " + err.Error())
		}

		if err := grpcServer.Stop(); err != nil {
			log.Error("failed to stop gRPC server: " + err.Error())
		}
	}()

	log.Info("app is running...")
	wg.Wait()
	log.Info("servers closed")
}
