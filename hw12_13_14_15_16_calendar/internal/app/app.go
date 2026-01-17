package app

import (
	"context"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	Close() error
	Write(ctx context.Context, event storage.Event) error
	Update(ctx context.Context, event storage.Event) error
	Delete(ctx context.Context, event storage.Event) error
	EventsByDay(ctx context.Context, date string) ([]storage.Event, error)
	EventsByWeek(ctx context.Context, date string) ([]storage.Event, error)
	EventsByMonth(ctx context.Context, date string) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
