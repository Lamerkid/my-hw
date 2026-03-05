package domain

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID     `db:"id"`
	Title        string        `db:"title"`
	Description  string        `db:"description,omitempty"`
	StartTime    time.Time     `db:"start_time"`
	EndTime      time.Time     `db:"end_time"`
	UserID       uuid.UUID     `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
	Notified     bool          `db:"notified"`
}

func (e Event) IsValid() bool {
	return e.Title != "" && e.StartTime.Before(e.EndTime)
}
