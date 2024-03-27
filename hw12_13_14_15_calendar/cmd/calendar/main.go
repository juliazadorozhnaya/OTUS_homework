package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/app"
	sqlstorage "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/storage/sql"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	internalhttp "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/api/server"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/config"
	"github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/juliazadorozhnaya/hw12_13_14_15_calendar/internal/storage/memory"
)

var (
	configPath    string
	configStorage string
)

func init() {
	defaultConfigPath := path.Join("configs", "config.toml")
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to configuration file")

	flag.StringVar(&configStorage, "storage", "mem", "Type of storage")
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

	if configStorage == "mem" {
		storage := memorystorage.New()
		application = app.New(storage)
	} else {
		storage := sqlstorage.New(conf.Database)
		application = app.New(storage)
	}

	server := internalhttp.NewServer(log, application, conf.Server)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		if err := server.Stop(); err != nil {
			log.Error("failed to stop api api: " + err.Error())
		}
	}()

	log.Info("app is running...")

	if err := server.Start(); !errors.Is(err, http.ErrServerClosed) && err != nil {
		log.Error("failed to start api api: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	log.Info("api closed")
}
