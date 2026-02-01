package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	UserID      uuid.UUID `json:"user_id"`
}

func (e Event) IsValid() bool {
	return e.Title != "" && e.StartTime.Before(e.EndTime)
}

func (e Event) Duration() time.Duration {
	return e.EndTime.Sub(e.StartTime)
}
