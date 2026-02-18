package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/configs"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/server"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar_config.yaml", "Path to configuration file")
}

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return 0
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	config, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	logg := logger.NewLogger(config.Calendar.Logger.Level)

	storage, err := storage.NewStorage(ctx, config.Calendar.Storage.Type,
		config.Calendar.Storage.DSN)
	if err != nil {
		logg.Error("failed to create storage: %v", err)
		return 2
	}
	defer storage.Close()

	calendar := app.NewApp(logg, storage)

	server, err := server.NewServer(logg, calendar, config.Calendar.Server.Type)
	if err != nil {
		logg.Error("failed to create server: %v", err)
		return 3
	}

	logg.Info("calendar is starting...")

	if err = server.Start(ctx,
		config.Calendar.Server.Host,
		config.Calendar.Server.Port,
		config.Calendar.Server.Timeout,
	); err != nil {
		logg.Error("failed to start server: %v", err)
		cancel()
		return 4
	}

	logg.Info("calendar has stopped running...")
	return 0
}
