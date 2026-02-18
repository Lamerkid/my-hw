package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	config "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/configs"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/rmq"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	flag.Parse()

	config, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg := logger.NewLogger(config.Scheduler.Logger.Level)

	rabbit := rmq.NewAMQP(logg,
		config.Scheduler.AMQP.URL, config.Scheduler.AMQP.Topic)

	storage, err := storage.NewStorage(ctx, config.Scheduler.Storage.Type,
		config.Scheduler.Storage.DSN)
	if err != nil {
		logg.Error("failed to create storage: %v", err)
		return 2
	}
	defer storage.Close()

	scheduler := scheduler.NewScheduler(logg, rabbit,
		storage, config.Scheduler.Interval)

	err = scheduler.Start(ctx)
	if err != nil {
		logg.Error("failed to start scheduler: %v", err)
		cancel()
		return 3
	}

	logg.Info("scheduler has stopped running...")
	return 0
}
