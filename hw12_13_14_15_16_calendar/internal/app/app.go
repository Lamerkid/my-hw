package app

import (
	"context"
	"time"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
)

type App struct {
	logger  Logger
	storage Storage
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event domain.Event) error {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.Write(appCtx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event domain.Event) error {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.Update(appCtx, event)
}

func (a *App) DeleteEvent(ctx context.Context, event domain.Event) error {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.Delete(appCtx, event)
}

func (a *App) SelectEventByDay(ctx context.Context, date string) ([]domain.Event, error) {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.EventsByDay(appCtx, date)
}

func (a *App) SelectEventByWeek(ctx context.Context, date string) ([]domain.Event, error) {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.EventsByWeek(appCtx, date)
}

func (a *App) SelectEventByMonth(ctx context.Context, date string) ([]domain.Event, error) {
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.EventsByMonth(appCtx, date)
}
