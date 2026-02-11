package internalgrpc

import (
	"context"
	"net"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	logger  Logger
	service *Service
	grpc    *grpc.Server
}

func NewServer(logger Logger, service *Service) *Server {
	return &Server{
		logger:  logger,
		service: service,
	}
}

func (s *Server) Start(ctx context.Context, host, port string, timeout time.Duration) error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}

	s.grpc = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryServerRequestLoggerInterceptor(),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Timeout: timeout,
		}),
	)

	reflection.Register(s.grpc)

	RegisterEventServiceServer(s.grpc, s.service)

	s.logger.Info("starting server on %s", lsn.Addr().String())

	serverErr := make(chan error, 1)

	go func() {
		if err := s.grpc.Serve(lsn); err != nil {
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
	done := make(chan struct{})
	go func() {
		s.grpc.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.grpc.Stop()
		return ctx.Err()
	}
}
