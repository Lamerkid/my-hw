package internalhttp

import "time"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// DTO
type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	UserID      string    `json:"user_id"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	UserID      string    `json:"user_id"`
}

type EventResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	UserID      string    `json:"user_id"`
}

type EventsListResponse struct {
	Events []EventResponse `json:"events"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
