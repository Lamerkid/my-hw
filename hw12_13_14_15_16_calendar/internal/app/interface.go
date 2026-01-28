package app

import (
	"context"
	"time"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	Close() error
	Write(ctx context.Context, event domain.Event) error
	Update(ctx context.Context, event domain.Event) error
	Delete(ctx context.Context, event domain.Event) error
	EventsByDay(ctx context.Context, date string) ([]domain.Event, error)
	EventsByWeek(ctx context.Context, date string) ([]domain.Event, error)
	EventsByMonth(ctx context.Context, date string) ([]domain.Event, error)
}

// Output interface
type EventService interface {
	CreateEvent(ctx context.Context, req CreateEventRequest) (domain.Event, error)
	UpdateEvent(ctx context.Context, req UpdateEventRequest) (domain.Event, error)
	DeleteEvent(ctx context.Context, req DeleteEventRequest) (domain.Event, error)
	SelectEventByDay(ctx context.Context, req SelectEventByDayRequest) ([]domain.Event, error)
	SelectEventByWeek(ctx context.Context, req SelectEventByWeekRequest) ([]domain.Event, error)
	SelectEventByMonth(ctx context.Context, req SelectEventByMonthRequest) ([]domain.Event, error)
}

type CreateEventRequest struct {
	Event domain.Event
}

type CreateEventResponse struct {
	Event domain.Event
}

type UpdateEventRequest struct {
	Event domain.Event
}

type UpdateEventResponse struct {
	Event domain.Event
}

type DeleteEventRequest struct {
	Event domain.Event
}

type DeleteEventResponse struct {
	Event domain.Event
}

type SelectEventByDayRequest struct {
	Date time.Time
}

type SelectEventByDayResponse struct {
	Event []domain.Event
}

type SelectEventByWeekRequest struct {
	Date time.Time
}

type SelectEventByWeekResponse struct {
	Event []domain.Event
}

type SelectEventByMonthRequest struct {
	Date time.Time
}

type SelectEventByMonthResponse struct {
	Event []domain.Event
}
