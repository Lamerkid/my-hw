package internalgrpc

import (
	context "context"
	"net"
	"testing"
	"time"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestGRPCServer(t *testing.T) {
	logg := logger.NewLogger("DEBUG")
	storage := memorystorage.New()
	calendar := app.NewApp(logg, storage)
	service := NewEventService(logg, calendar)
	grpcSrv := NewServer(logg, service)

	ctx := context.Background()

	host := "127.0.0.1"
	port := "8081"
	timeout := 30 * time.Second

	go func() {
		if err := grpcSrv.Start(ctx, host, port, timeout); err != nil {
			t.Errorf("Server error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.NewClient(net.JoinHostPort(host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := NewEventServiceClient(conn)

	t.Run("CreateEvent", func(t *testing.T) {
		resp, err := client.CreateEvent(ctx, &CreateEventRequest{
			Title:        "Test Event",
			Description:  "Test Description",
			StartTime:    timestamppb.New(time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)),
			EndTime:      timestamppb.New(time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)),
			UserId:       "40ab8025-d98d-41ed-a151-b6a9ee23ccb3",
			NotifyBefore: "15m",
		})

		require.NoError(t, err)
		require.Equal(t, "Test Event", resp.GetEvent().GetTitle())
		require.NotEmpty(t, resp.GetEvent().GetId())
	})

	t.Run("SelectEventByDay", func(t *testing.T) {
		targetDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

		arrayResp, err := client.SelectEventByDay(ctx, &SelectEventRequest{
			Date: timestamppb.New(targetDate),
		})

		require.NoError(t, err)
		require.Len(t, arrayResp.GetEvents(), 1)
	})
}
