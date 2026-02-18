package scheduler

import (
	"context"
	"time"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/rmq"
)

type Scheduler struct {
	storage  Storage
	producer rmq.Producer
	logger   Logger
	interval string
}

func NewScheduler(logg Logger, producer rmq.Producer, storage Storage, interval string) *Scheduler {
	return &Scheduler{
		storage:  storage,
		producer: producer,
		logger:   logg,
		interval: interval,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.logger.Info("startin scheduler...")

	schedInteval, err := time.ParseDuration(s.interval)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(schedInteval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.processEvents(ctx); err != nil {
				s.logger.Error("failed to process events: %v", err)
			}
		}
	}
}

func (s *Scheduler) processEvents(ctx context.Context) error {
	events, err := s.storage.EventsForNotification(ctx)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		s.logger.Debug("no events to notify")
		return nil
	}

	for _, event := range events {
		s.logger.Debug("publishing event: %s", event.ID)

		notification := rmq.Notification{
			EventID:    event.ID,
			EventTitle: event.Title,
			EventTime:  event.StartTime,
			UserID:     event.UserID,
		}

		if err := s.producer.Publish(ctx, notification); err != nil {
			s.logger.Error("failed to publish notification: %v", err)
			continue
		}

		if err := s.storage.MarkNotified(ctx, event.ID); err != nil {
			s.logger.Error("failed to mark event notified")
		}

		s.logger.Debug("marked notified event: %s", event.ID)
	}

	return nil
}
