package rmq

import "time"

type Notification struct {
	EventID     string    `json:"eventId"`
	UserID      string    `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"eventTime"`
	NotifyAt    time.Time `json:"notifyAt"`
	Type        string    `json:"type"` // "upcoming", "started", etc.
}
