package sender

import (
	"context"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/rmq"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type NotificationHandler interface {
	Send(ctx context.Context, msg rmq.Notification) error
}
