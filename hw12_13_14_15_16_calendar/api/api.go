package api

import (
	"time"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	UserId      uuid.UUID `json:"user_id"`
}
