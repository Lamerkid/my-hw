package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
)

type App struct {
	logger  Logger
	storage Storage
}

func NewApp(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, cmd CreateEventCommand) (domain.Event, error) {
	a.logger.Debug("creating event: %s", cmd.Title)

	id, err := uuid.NewRandom()
	if err != nil {
		return domain.Event{}, err
	}

	notifyBefore, err := time.ParseDuration(cmd.NotifyBefore)
	if err != nil {
		return domain.Event{}, err
	}

	event := domain.Event{
		ID:           id,
		Title:        cmd.Title,
		Description:  cmd.Description,
		StartTime:    cmd.StartTime,
		EndTime:      cmd.EndTime,
		UserID:       cmd.UserID,
		NotifyBefore: notifyBefore,
	}

	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = a.storage.Write(appCtx, event)
	if err != nil {
		return domain.Event{}, err
	}

	a.logger.Info("created event with id: %s", event.ID)
	return event, nil
}

func (a *App) UpdateEvent(ctx context.Context, cmd UpdateEventCommand) (domain.Event, error) {
	a.logger.Debug("updating event with id: %s", cmd.ID)

	notifyBefore, err := time.ParseDuration(cmd.NotifyBefore)
	if err != nil {
		return domain.Event{}, err
	}

	event := domain.Event{
		ID:           cmd.ID,
		Title:        cmd.Title,
		Description:  cmd.Description,
		StartTime:    cmd.StartTime,
		EndTime:      cmd.EndTime,
		NotifyBefore: notifyBefore,
	}

	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := a.storage.Update(appCtx, event); err != nil {
		return domain.Event{}, err
	}

	a.logger.Info("updated event with id: %s", event.ID)
	return event, nil
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) (domain.Event, error) {
	a.logger.Debug("deleting event with id: %s", id)

	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	event, err := a.storage.GetEvent(appCtx, id)
	if err != nil {
		return domain.Event{}, err
	}

	err = a.storage.Delete(appCtx, id)
	if err != nil {
		return domain.Event{}, err
	}

	a.logger.Info("deleted event with id: %s", id)
	return event, nil
}

func (a *App) SelectEventByDay(ctx context.Context, date time.Time) ([]domain.Event, error) {
	a.logger.Info("selecting event by day from: %s", date.String())
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.EventsByDay(appCtx, date.Format(time.DateOnly))
}

func (a *App) SelectEventByWeek(ctx context.Context, date time.Time) ([]domain.Event, error) {
	a.logger.Info("selecting event by week from: %s", date.String())
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.EventsByWeek(appCtx, date.Format(time.DateOnly))
}

func (a *App) SelectEventByMonth(ctx context.Context, date time.Time) ([]domain.Event, error) {
	a.logger.Info("selecting event by week from: %s", date.String())
	appCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.storage.EventsByMonth(appCtx, date.Format(time.DateOnly))
}
