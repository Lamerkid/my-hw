package sender

import (
	"context"
	"encoding/json"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/rmq"
)

type MessageHandler struct {
	logger Logger
}

func NewMessageHandler(logger Logger) *MessageHandler {
	return &MessageHandler{
		logger: logger,
	}
}

func (h *MessageHandler) Handle(ctx context.Context, body []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var notification rmq.Notification

		if err := json.Unmarshal(body, &notification); err != nil {
			return err
		}

		h.logger.Debug("processing notification",
			"event_id", notification.EventID,
			"user_id", notification.UserID)

		h.logger.Info("notification %s sent successfully to %s", notification.EventID, notification.UserID)

		return nil
	}
}
