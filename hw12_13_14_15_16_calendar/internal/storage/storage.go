package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	config "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/configs"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
	memorystorage "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

type Type int

const (
	StorageTypeMemory Type = iota + 1
	StorageTypePostgres
)

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

func NewStorage(ctx context.Context, cfg config.Config) (Storage, error) {
	switch cfg.Storage.Type {
	case "inMemory":
		return memorystorage.New(), nil
	case "postgres":
		s := sqlstorage.New()
		if err := s.Connect(ctx, cfg.Storage.DSN); err != nil {
			return nil, fmt.Errorf("connect to postgres: %w", err)
		}

		return s, nil

	default:
		return nil, fmt.Errorf("unknown storage type: %s", cfg.Storage.Type)
	}
}
