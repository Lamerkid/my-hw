package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
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
	os.Exit(run())
}

func run() (exitCode int) {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return 0
	}

	config := New(configFile)
	if err := config.ValidateConfig(); err != nil {
		return 1
	}

	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := NewStorage(ctx, config)
	if err != nil {
		return 2
	}
	defer storage.Close()

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

	err = server.Start(config.Host, config.Port, config.Timeout)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		return 3
	}

	return 0
}

func NewStorage(ctx context.Context, cfg *Config) (app.Storage, error) {
	switch cfg.Storage.Type {
	case "inMemory":
		return memorystorage.New(), nil
	case "Postgres":
		sqlStorage := sqlstorage.New()
		if err := sqlStorage.Connect(ctx, cfg.Storage.DSN); err != nil {
			return nil, err
		}
		return sqlStorage, nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Storage.Type)
	}
}
