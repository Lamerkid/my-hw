package rmq

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	EventID    uuid.UUID `json:"eventId"`
	EventTitle string    `json:"eventTitle"`
	EventTime  time.Time `json:"eventTime"`
	UserID     uuid.UUID `json:"userId"`
}
