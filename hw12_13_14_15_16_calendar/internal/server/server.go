package server

import (
	"context"
	"fmt"
	"time"

	config "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/configs"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/server/http"
)

type Type int

const (
	ServerTypeHTTP Type = iota + 1
	ServerTypeGRPC
)

type Server interface {
	Start(ctx context.Context, host, port string, timeout time.Duration) error
}

func NewServer(cfg config.Config, logg *logger.Logger, app *app.App) (Server, error) {
	switch cfg.Calendar.Server.Type {
	case "grpc":
		service := internalgrpc.NewEventService(logg, app)
		server := internalgrpc.NewServer(logg, service)

		return server, nil

	case "http":
		handler := internalhttp.NewHandler(logg, app)
		server := internalhttp.NewServer(logg, handler)

		return server, nil

	default:
		return nil, fmt.Errorf("unknown server type: %s", cfg.Calendar.Storage.Type)
	}
}
