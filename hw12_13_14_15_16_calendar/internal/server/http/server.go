package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
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
	CreateEvent(ctx context.Context, id, title string) error
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
	mux.HandleFunc("/hello", s.Hello)

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
