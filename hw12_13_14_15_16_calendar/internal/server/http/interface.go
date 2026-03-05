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
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	UserID       string    `json:"userId"`
	NotifyBefore string    `json:"notifyBefore"`
}

type UpdateEventRequest struct {
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	NotifyBefore string    `json:"notifyBefore"`
}

type EventResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	UserID       string    `json:"userId"`
	NotifyBefore string    `json:"notifyBefore"`
}

type EventsListResponse struct {
	Events []EventResponse
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
