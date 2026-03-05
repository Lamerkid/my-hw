package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	calendarURL := os.Getenv("CALENDAR_URL")

	now := time.Now()
	ctx := context.Background()
	dateStr := now.Format("2006-01-02")

	data := map[string]any{
		"title":        "Test Event Integration",
		"description":  "Test Description",
		"startTime":    now.Add(10 * time.Minute),
		"endTime":      now.Add(30 * time.Minute),
		"userId":       "40ab8025-d98d-41ed-a151-b6a9ee23ccb3",
		"notifyBefore": "15m",
	}

	t.Run("Create event", func(t *testing.T) {
		jsonData, err := json.Marshal(data)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, "POST",
			calendarURL+"/api/v3/event", bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(bodyBytes), "Title")
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Update non existing event", func(t *testing.T) {
		jsonData, err := json.Marshal(data)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(ctx, "PUT",
			calendarURL+"/api/v3/event/fe1a7632-4974-43c3-a391-3d82e211c45", bytes.NewBuffer(jsonData))
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(bodyBytes), "code")
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Select events", func(t *testing.T) {
		selectEvents(t, calendarURL, dateStr, "day")
		selectEvents(t, calendarURL, dateStr, "week")
		selectEvents(t, calendarURL, dateStr, "month")
	})

	t.Run("Check notification", func(t *testing.T) {
		time.Sleep(11 * time.Second)

		// Connect to RabbitMQ
		rabbitConn, err := amqp.Dial(rabbitmqURL)
		require.NoError(t, err)

		rabbitCh, err := rabbitConn.Channel()
		require.NoError(t, err)

		queue, err := rabbitCh.QueueDeclare(
			"test-queue",
			false,
			true,
			false,
			false,
			nil,
		)
		require.NoError(t, err)

		// Publish message
		msg := map[string]string{"test": "data"}
		body, _ := json.Marshal(msg)

		err = rabbitCh.PublishWithContext(
			context.Background(),
			"",
			queue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		require.NoError(t, err)

		// Consume message
		msgs, err := rabbitCh.Consume(
			queue.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		require.NoError(t, err)

		select {
		case delivery := <-msgs:
			var result map[string]string
			err = json.Unmarshal(delivery.Body, &result)
			require.NoError(t, err)
			require.Equal(t, "data", result["test"])
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for message")
		}
	})
}

func selectEvents(t *testing.T, url, dateStr, dateRange string) {
	t.Helper()
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET",
		url+"/api/v3/event/select?date="+dateStr+"&range="+dateRange, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(bodyBytes), "Title")
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
