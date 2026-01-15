package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := New()
	if err := config.ReadConfig(configFile); err != nil {
		fmt.Println("Error readin config file: ", err)
		os.Exit(1)
	}
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg := logger.New(config.Logger.Level)
	storage, cleanup, err := setupStorage(ctx, *config, *logg)
	if err != nil {
		logg.Error("Failed to setup storage: " + err.Error())
		os.Exit(1)
	}
	defer cleanup()

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func setupStorage(ctx context.Context, config Config, logger logger.Logger) (app.Storage, func(), error) {
	switch config.Storage {
	case "inMemory":
		return memorystorage.New(), func() {}, nil
	case "Postgres":
		storage := sqlstorage.New()
		if err := storage.Connect(ctx, config.DBConnetion); err != nil {
			return nil, nil, err
		}

		cleanup := func() {
			if err := storage.Close(ctx); err != nil {
				logger.Error("Failed to close db connection: " + err.Error())
			}
		}
		return storage, cleanup, nil
	default:
		return nil, nil, fmt.Errorf("unknown storage type: %s", config.Storage)
	}
}
