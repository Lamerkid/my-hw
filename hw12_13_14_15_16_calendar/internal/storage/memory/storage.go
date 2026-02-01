package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
)

type Storage struct {
	data map[uuid.UUID]domain.Event
	mu   sync.RWMutex
}

func New() *Storage {
	return &Storage{
		data: make(map[uuid.UUID]domain.Event),
	}
}

func (s *Storage) Close() error {
	s.data = nil
	return nil
}

func (s *Storage) Write(ctx context.Context, event domain.Event) error {
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

func (s *Storage) Update(ctx context.Context, event domain.Event) error {
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

func (s *Storage) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, exists := s.data[id]
		if !exists {
			return fmt.Errorf("entry is not present in storage")
		}
		delete(s.data, id)
		return nil
	}
}

func (s *Storage) GetEvent(ctx context.Context, id uuid.UUID) (domain.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	select {
	case <-ctx.Done():
		return domain.Event{}, ctx.Err()
	default:
		event, exists := s.data[id]
		if !exists {
			return domain.Event{}, fmt.Errorf("entry is not present in storage")
		}
		return event, nil
	}
}

func (s *Storage) EventsByDay(ctx context.Context, date string) ([]domain.Event, error) {
	return CollectEvents(ctx, s, date, 0, 1)
}

func (s *Storage) EventsByWeek(ctx context.Context, date string) ([]domain.Event, error) {
	return CollectEvents(ctx, s, date, 0, 7)
}

func (s *Storage) EventsByMonth(ctx context.Context, date string) ([]domain.Event, error) {
	return CollectEvents(ctx, s, date, 1, 0)
}

func CollectEvents(ctx context.Context, s *Storage, date string, month, day int) ([]domain.Event, error) {
	var events []domain.Event
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
