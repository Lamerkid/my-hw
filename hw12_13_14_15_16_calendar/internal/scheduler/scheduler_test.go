package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/domain"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/rmq"
	memorystorage "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

var event1 = domain.Event{
	ID:           uuid.New(),
	Title:        "test from scheduler",
	StartTime:    time.Now().Add(5 * time.Minute),
	EndTime:      time.Now().Add(30 * time.Minute),
	Description:  "testing",
	UserID:       uuid.New(),
	NotifyBefore: 10 * time.Minute,
}

func TestScheduler(t *testing.T) {
	logg := logger.NewLogger("DEBUG")
	brocker := rmq.NewMockBroker()
	storage := memorystorage.New()
	scheduler := NewScheduler(logg, brocker, storage, 500*time.Millisecond)

	ctx := context.Background()

	storage.Write(ctx, event1)

	go func() {
		if err := scheduler.Start(ctx); err != nil {
			t.Errorf("Scheduler error: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	outEvent, err := storage.GetEvent(ctx, event1.ID)
	require.NoError(t, err)
	require.True(t, outEvent.Notified)

	notif := brocker.Messages
	require.NotEmpty(t, notif)
}
