package model

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	StartTime   time.Time `db:"start_time"`
	EndTime     time.Time `db:"end_time"`
	Description string    `db:"description"`
	UserID      uuid.UUID `db:"user_id"`
}
