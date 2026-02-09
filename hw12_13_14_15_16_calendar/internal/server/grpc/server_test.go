package internalgrpc

import (
	context "context"
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

func TestGRPC_CreateEvent(t *testing.T) {
	logg := logger.NewLogger("DEBUG")
	storage := memorystorage.New()
	calendar := app.NewApp(logg, storage)
	service := NewEventService(logg, calendar)
	srv := NewServer(logg, service)

	go func() {
		if err := srv.Start(context.Background(), "127.0.0.1", "8080", 30*time.Second); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := grpc.NewClient("127.0.0.1:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := NewEventServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CreateEvent(ctx, &CreateEventRequest{
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   timestamppb.New(time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)),
		EndTime:     timestamppb.New(time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)),
		UserId:      "40ab8025-d98d-41ed-a151-b6a9ee23ccb3",
	})

	require.NoError(t, err)
	require.Equal(t, "Test Event", resp.GetEvent().GetTitle())
	require.NotEmpty(t, resp.GetEvent().GetId())

	targetDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

	arrayResp, err := client.SelectEventByDay(ctx, &SelectEventRequest{
		Date: timestamppb.New(targetDate),
	})

	require.NoError(t, err)
	require.Len(t, arrayResp.GetEvents(), 1)
}
