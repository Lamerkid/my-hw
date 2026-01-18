package storage

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	Title       string
	StartTime   time.Time
	EndTime     time.Time
	Description string
	UserID      uuid.UUID
}
