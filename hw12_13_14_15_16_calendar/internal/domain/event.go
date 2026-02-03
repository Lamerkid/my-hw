package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	UserID      uuid.UUID `json:"userId"`
}

func (e Event) IsValid() bool {
	return e.Title != "" && e.StartTime.Before(e.EndTime)
}
