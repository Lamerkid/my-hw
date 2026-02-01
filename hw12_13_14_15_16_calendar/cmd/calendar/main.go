package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	config "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/configs"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage"
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

	config, err := config.LoadConfig(configFile)
	if err != nil {
		return 1
	}

	logg := logger.NewLogger(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := storage.NewStorage(ctx, config)
	if err != nil {
		return 2
	}
	defer storage.Close()

	calendar := app.NewApp(logg, storage)

	handler := internalhttp.NewHandler(logg, calendar)

	server := internalhttp.NewServer(logg, handler)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: %v", err)
		}
	}()

	logg.Info("calendar is running...")

	if err = server.Start(ctx, config.Server.Host, config.Server.Port, config.Server.Timeout); err != nil {
		logg.Error("failed to start http server: %v", err)
		cancel()
		return 3
	}

	logg.Info("calendar has stoped running...")
	return 0
}
