package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestHTTPServer(t *testing.T) {
	logg := logger.NewLogger("DEBUG")
	storage := memorystorage.New()
	calendar := app.NewApp(logg, storage)
	handler := NewHandler(logg, calendar)
	httpSrv := NewServer(logg, handler)

	ctx := context.Background()

	host := "127.0.0.1"
	port := "8080"
	timeout := 30 * time.Second

	go func() {
		if err := httpSrv.Start(ctx, host, port, timeout); err != nil {
			t.Errorf("Server error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	t.Run("Hello world", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, "GET", "http://"+net.JoinHostPort(host, port)+"/hello", nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "Hello world!\n", string(bodyBytes))
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Create event", func(t *testing.T) {
		data := map[string]any{
			"title":        "Test Event",
			"description":  "Test Description",
			"startTime":    time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC),
			"endTime":      time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC),
			"userId":       "40ab8025-d98d-41ed-a151-b6a9ee23ccb3",
			"notifyBefore": "15m",
		}

		jsonData, err := json.Marshal(data)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, "POST",
			"http://"+net.JoinHostPort(host, port)+"/api/v3/event", bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(bodyBytes), "Title")
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Select event by day", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, "GET",
			"http://"+net.JoinHostPort(host, port)+"/api/v3/event/select?date=2026-01-15&range=day", nil)
		require.NoError(t, err)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(bodyBytes), "Title")
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
