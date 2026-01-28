package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
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
	CreateEvent(ctx context.Context, event domain.Event) error
	UpdateEvent(ctx context.Context, event domain.Event) error
	DeleteEvent(ctx context.Context, event domain.Event) error
	SelectEventByDay(ctx context.Context, date string) ([]domain.Event, error)
	SelectEventByWeek(ctx context.Context, date string) ([]domain.Event, error)
	SelectEventByMonth(ctx context.Context, date string) ([]domain.Event, error)
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
	mux.HandleFunc("/api/v3/event", s.CreateEvent)
	mux.HandleFunc("/api/v3/event/", s.SelectEventByDay) // get, put, delete

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
