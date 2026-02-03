package app

import (
	"context"
	"time"

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
	Close() error
	Write(ctx context.Context, event domain.Event) error
	Update(ctx context.Context, event domain.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetEvent(ctx context.Context, id uuid.UUID) (domain.Event, error)
	EventsByDay(ctx context.Context, date string) ([]domain.Event, error)
	EventsByWeek(ctx context.Context, date string) ([]domain.Event, error)
	EventsByMonth(ctx context.Context, date string) ([]domain.Event, error)
}

// EventService output interface.

type EventService interface {
	CreateEvent(ctx context.Context, cmd CreateEventCommand) (domain.Event, error)
	UpdateEvent(ctx context.Context, cmd UpdateEventCommand) (domain.Event, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) (domain.Event, error)
	SelectEventByDay(ctx context.Context, date time.Time) ([]domain.Event, error)
	SelectEventByWeek(ctx context.Context, date time.Time) ([]domain.Event, error)
	SelectEventByMonth(ctx context.Context, date time.Time) ([]domain.Event, error)
}

type CreateEventCommand struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	UserID      uuid.UUID `json:"userId"`
}

type UpdateEventCommand struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
}
