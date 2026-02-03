package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	logger  Logger
	handler *Handler
	http    *http.Server
}

func NewServer(logger Logger, handler *Handler) *Server {
	return &Server{
		logger:  logger,
		handler: handler,
	}
}

func (s *Server) Start(ctx context.Context, host, port string, timeout time.Duration) error {
	url := net.JoinHostPort(host, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.handler.Hello)
	mux.HandleFunc("/api/v3/event", s.handler.CreateEvent)
	mux.HandleFunc("/api/v3/event/{id}", s.handler.EventHandler)
	mux.HandleFunc("/api/v3/event/select", s.handler.SelectEventHandler)

	handler := loggingMiddleware(mux)

	s.http = &http.Server{
		Addr:              url,
		Handler:           handler,
		ReadHeaderTimeout: timeout,
	}

	serverErr := make(chan error, 1)

	go func() {
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		} else {
			serverErr <- nil
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return s.Stop(shutdownCtx)

	case err := <-serverErr:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping server")

	if s.http != nil {
		s.http.SetKeepAlivesEnabled(false)
		if err := s.http.Shutdown(ctx); err != nil {
			<-ctx.Done()
			return err
		}
	}
	return nil
}
