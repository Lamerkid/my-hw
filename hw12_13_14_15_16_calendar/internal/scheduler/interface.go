package scheduler

import (
	"context"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Storage interface {
	EventsForNotification(ctx context.Context) ([]domain.Event, error)
	MarkNotified(ctx context.Context, id uuid.UUID) error
}
