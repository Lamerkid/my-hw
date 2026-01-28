package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
}

func (e Event) Valid() bool {
	return e.Title != "" && e.StartTime.Before(e.EndTime)
}
