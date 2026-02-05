package internalgrpc

import (
	"context"
	"net"
	"time"

	grpc "google.golang.org/grpc"
)

type Server struct {
	server UnimplementedEventServiceServer
	logger Logger
	grpc   *grpc.Server
}

func NewServer(logger Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context, host, port string, timeout time.Duration) error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}

	s.grpc = grpc.NewServer(
		grpc.UnaryServerInterceptor(
			UnaryServerRequestLoggerInterceptor,
		),
	)

	RegisterEventServiceServer(s.grpc, &Server{})

	s.logger.Info("starting server on %s", lsn.Addr().String())

	if err := s.grpc.Serve(lsn); err != nil {
		return err
	}

	return nil
}

func (s *Server) CreateEvent(ctx context.Context, in *CreateEventRequest) (*EventResponse, error) {
	s.logger.Debug("Creating event with id: %s", in.Event.GetId())
	return s.server.CreateEvent(ctx, in)
}

func (s *Server) UpdateEvent(ctx context.Context, in *EventRequest) (*EventResponse, error) {
	s.logger.Debug("Updating event with id: %s", in.GetId())
	return s.server.UpdateEvent(ctx, in)
}

func (s *Server) DeleteEvent(ctx context.Context, in *EventRequest) (*EventResponse, error) {
	s.logger.Debug("Deleting event with id: %s", in.GetId())
	return s.server.DeleteEvent(ctx, in)
}

func (s *Server) SelectEventByDay(ctx context.Context, in *SelectEventRequest) (*ArrayEventResponse, error) {
	s.logger.Debug("Selecting events on date: %s, by range: %s", in.Date.String(), in.Range)
	return s.server.SelectEventByDay(ctx, in)
}

func (s *Server) SelectEventByWeek(ctx context.Context, in *SelectEventRequest) (*ArrayEventResponse, error) {
	s.logger.Debug("Selecting events on date: %s, by range: %s", in.Date.String(), in.Range)
	return s.server.SelectEventByWeek(ctx, in)
}

func (s *Server) SelectEventByMonth(ctx context.Context, in *SelectEventRequest) (*ArrayEventResponse, error) {
	s.logger.Debug("Selecting events on date: %s, by range: %s", in.Date.String(), in.Range)
	return s.server.SelectEventByMonth(ctx, in)
}
