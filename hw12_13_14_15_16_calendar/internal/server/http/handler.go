package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
)

type Handler struct {
	app    app.EventService
	logger Logger
}

func NewHandler(logger Logger, app app.EventService) *Handler {
	return &Handler{
		app:    app,
		logger: logger,
	}
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello to %s\n", r.Host)
}

func (h *Handler) EventByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	switch r.Method {
	case http.MethodPut:
		var req app.UpdateEventCommand
		h.UpdateEvent(w, r, req)
	case http.MethodDelete:
		h.DeleteEvent(w, r, idStr)
	}
}

func (h *Handler) SelectEventHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "date parameter is required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rangeParam := r.URL.Query().Get("range")
	switch rangeParam {
	case "day":
		h.SelectEventByDay(w, r, date)
	case "week":
		h.SelectEventByWeek(w, r, date)
	case "month":
		h.SelectEventByMonth(w, r, date)
	default:
		http.Error(w, "range must be 'day', 'week', or 'month'", http.StatusBadRequest)
		return
	}
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	var req app.CreateEventCommand
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.app.CreateEvent(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Contet-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request, req app.UpdateEventCommand) {
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.app.UpdateEvent(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Contet-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.app.DeleteEvent(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Contet-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *Handler) SelectEventByDay(w http.ResponseWriter, r *http.Request, date time.Time) {
	events, err := h.app.SelectEventByDay(r.Context(), date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		h.logger.Error("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) SelectEventByWeek(w http.ResponseWriter, r *http.Request, date time.Time) {
	events, err := h.app.SelectEventByWeek(r.Context(), date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Contet-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *Handler) SelectEventByMonth(w http.ResponseWriter, r *http.Request, date time.Time) {
	events, err := h.app.SelectEventByMonth(r.Context(), date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Contet-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
