package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	data map[uuid.UUID]storage.Event
	mu   sync.RWMutex
}

func New() *Storage {
	return &Storage{
		data: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) Close() error {
	s.data = nil
	return nil
}

func (s *Storage) Write(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, exists := s.data[event.ID]
		if exists {
			return fmt.Errorf("entry already exists in storage")
		}
		s.data[event.ID] = event
		return nil
	}
}

func (s *Storage) Update(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, exists := s.data[event.ID]
		if !exists {
			return fmt.Errorf("entry is not present in storage")
		}
		s.data[event.ID] = event
		return nil
	}
}

func (s *Storage) Delete(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, exists := s.data[event.ID]
		if !exists {
			return fmt.Errorf("entry is not present in storage")
		}
		delete(s.data, event.ID)
		return nil
	}
}

func (s *Storage) EventsByDay(ctx context.Context, date string) ([]storage.Event, error) {
	return CollectEvents(ctx, s, date, 0, 1)
}

func (s *Storage) EventsByWeek(ctx context.Context, date string) ([]storage.Event, error) {
	return CollectEvents(ctx, s, date, 0, 7)
}

func (s *Storage) EventsByMonth(ctx context.Context, date string) ([]storage.Event, error) {
	return CollectEvents(ctx, s, date, 1, 0)
}

func CollectEvents(ctx context.Context, s *Storage, date string, month, day int) ([]storage.Event, error) {
	var events []storage.Event
	parsedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, err
	}
	endDate := parsedDate.AddDate(0, month, day)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.data {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if !event.StartTime.Before(parsedDate) && event.StartTime.Before(endDate) {
				events = append(events, event)
			}
		}
	}
	return events, nil
}
