package app

import (
	"context"
	"time"

	model "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage/models"
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
	Write(ctx context.Context, event model.Event) error
	Update(ctx context.Context, event model.Event) error
	Delete(ctx context.Context, event model.Event) error
	EventsByDay(ctx context.Context, date string) ([]model.Event, error)
	EventsByWeek(ctx context.Context, date string) ([]model.Event, error)
	EventsByMonth(ctx context.Context, date string) ([]model.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event model.Event) error {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.Write(appCtx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event model.Event) error {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.Update(appCtx, event)
}

func (a *App) DeleteEvent(ctx context.Context, event model.Event) error {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.Delete(appCtx, event)
}

func (a *App) SelectEventByDay(ctx context.Context, date string) ([]model.Event, error) {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.EventsByDay(appCtx, date)
}

func (a *App) SelectEventByWeek(ctx context.Context, date string) ([]model.Event, error) {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.EventsByWeek(appCtx, date)
}

func (a *App) SelectEventByMonth(ctx context.Context, date string) ([]model.Event, error) {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.EventsByMonth(appCtx, date)
}
