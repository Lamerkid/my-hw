package memorystorage

import (
	"context"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

var event1 = storage.Event{
	ID:          uuid.New(),
	Title:       "test1",
	StartTime:   time.Date(2026, 1, 15, 0, 30, 0, 0, time.UTC),
	EndTime:     time.Date(2026, 1, 15, 0, 30, 0, 0, time.UTC).Add(30 * time.Minute),
	Description: "testing",
	UserID:      uuid.New(),
}

var event2 = storage.Event{
	ID:          uuid.New(),
	Title:       "test2",
	StartTime:   time.Date(2026, 1, 15, 0, 30, 0, 0, time.UTC).AddDate(0, 0, 5),
	EndTime:     time.Date(2026, 1, 15, 0, 30, 0, 0, time.UTC).AddDate(0, 0, 5).Add(30 * time.Minute),
	Description: "testing",
	UserID:      uuid.New(),
}

var event3 = storage.Event{
	ID:          uuid.New(),
	Title:       "test3",
	StartTime:   time.Date(2026, 1, 15, 0, 30, 0, 0, time.UTC).AddDate(0, 0, 10),
	EndTime:     time.Date(2026, 1, 15, 0, 30, 0, 0, time.UTC).AddDate(0, 0, 10).Add(30 * time.Minute),
	Description: "testing",
	UserID:      uuid.New(),
}

var event4 = storage.Event{
	ID:          uuid.New(),
	Title:       "test4",
	StartTime:   time.Date(2026, 1, 15, 0, 30, 0, 0, time.UTC).AddDate(0, 1, 1),
	EndTime:     time.Date(2026, 1, 15, 0, 30, 0, 0, time.UTC).AddDate(0, 1, 1).Add(30 * time.Minute),
	Description: "testing",
	UserID:      uuid.New(),
}

func TestStorage(t *testing.T) {
	ctx := context.Background()
	storage := New()
	t.Run("basic functions", func(t *testing.T) {
		err := storage.Write(ctx, event1)
		require.NoError(t, err)

		err = storage.Write(ctx, event2)
		require.NoError(t, err)

		err = storage.Write(ctx, event3)
		require.NoError(t, err)

		err = storage.Write(ctx, event4)
		require.NoError(t, err)

		err = storage.Write(ctx, event2)
		require.Error(t, err)

		event1.Description = "more testing"
		err = storage.Update(ctx, event1)
		require.NoError(t, err)

		dayResult, err := storage.EventsByDay(ctx, "2026-01-15")
		require.NoError(t, err)
		require.Len(t, dayResult, 1)

		weekResult, err := storage.EventsByWeek(ctx, "2026-01-15")
		require.NoError(t, err)
		require.Len(t, weekResult, 2)

		monthResult, err := storage.EventsByMonth(ctx, "2026-01-15")
		require.NoError(t, err)
		require.Len(t, monthResult, 3)

		err = storage.Delete(ctx, event1)
		require.NoError(t, err)

		err = storage.Delete(ctx, event2)
		require.NoError(t, err)

		err = storage.Delete(ctx, event3)
		require.NoError(t, err)

		err = storage.Delete(ctx, event4)
		require.NoError(t, err)

		require.Len(t, storage.data, 0)

		err = storage.Delete(ctx, event3)
		require.Error(t, err)
	})
}

func BenchmarkStorage(b *testing.B) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	storage := New()

	b.ResetTimer()
	b.Run("write concurrently", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				storage.Write(ctx, event1)
				storage.Update(ctx, event1)
				storage.Delete(ctx, event1)
			}
		})
	})

	b.ResetTimer()
	b.Run("read concurrently", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				storage.EventsByDay(ctx, "2026-01-15")
				storage.EventsByWeek(ctx, "2026-01-15")
				storage.EventsByMonth(ctx, "2026-01-15")
			}
		})
	})

	b.ResetTimer()
	b.Run("rw concurently", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				storage.Write(ctx, event1)
				storage.EventsByDay(ctx, "2026-01-15")
				storage.Delete(ctx, event1)
			}
		})
	})
}
