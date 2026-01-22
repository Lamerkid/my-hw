package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/api"
	model "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage/models"
)

type Server struct {
	app    Application
	logger Logger
	http   *http.Server
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, event model.Event) error
	UpdateEvent(ctx context.Context, event model.Event) error
	DeleteEvent(ctx context.Context, event model.Event) error
	SelectEventByDay(ctx context.Context, date string) ([]model.Event, error)
	SelectEventByWeek(ctx context.Context, date string) ([]model.Event, error)
	SelectEventByMonth(ctx context.Context, date string) ([]model.Event, error)
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(host, port string, timeout time.Duration) error {
	url := net.JoinHostPort(host, port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Hello)
	// mux.HandleFunc("/api/v1/events", s.CreateEvent) // post, get
	// mux.HandleFunc("/api/v1/events/UUID", s.CreateEvent) // put, delete

	handler := loggingMiddleware(mux)

	s.http = &http.Server{
		Addr:              url,
		Handler:           handler,
		ReadHeaderTimeout: timeout,
	}

	return s.http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping server")
	s.http.SetKeepAlivesEnabled(false)
	if err := s.http.Shutdown(ctx); err != nil {
		<-ctx.Done()
		return err
	}
	return nil
}

func (s *Server) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello to %s\n", r.Host)
}

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req api.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	id, err := uuid.NewRandom()
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
	}

	event := model.Event{
		ID:          id,
		Title:       req.Title,
		StartTime:   req.StartTime,
		EndTime:     time.Now(),
		Description: *req.Description,
		UserID:      req.UserId,
	}

	err = s.app.CreateEvent(ctx, event)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	s.respond(w, r, http.StatusOK, event)
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, err)
}

func (s *Server) respond(w http.ResponseWriter, _ *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
