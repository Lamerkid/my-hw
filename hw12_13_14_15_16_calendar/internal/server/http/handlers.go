package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
)

type Handler struct {
	eventService app.EventService
}

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello to %s\n", r.Host)
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req app.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	id, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := req.Event.UserID

	event := domain.Event{
		ID:          id,
		Title:       req.Event.Title,
		StartTime:   time.Now(),
		EndTime:     time.Now(),
		Description: req.Event.Description,
		UserID:      userId,
	}

	err = s.app.CreateEvent(ctx, event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Contet-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (s *Server) SelectEventByDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req app.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date := r.URL.Query().Get("date")

	// path := r.URL.Path
	// date := strings.TrimPrefix(path, "/api/v3/events")
	// date = strings.TrimSuffix(id, "/")

	ctx := r.Context()
	events, err := s.app.SelectEventByDay(ctx, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Contet-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
