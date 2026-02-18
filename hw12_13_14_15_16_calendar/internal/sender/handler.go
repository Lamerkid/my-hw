package sender

import (
	"context"
	"encoding/json"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/rmq"
)

type NotificationSender interface {
	Send(ctx context.Context, msg rmq.Notification) error
}

type MessageHandler struct {
	sender NotificationSender
	logger Logger
}

func NewMessageHandler(sender NotificationSender, logger Logger) *MessageHandler {
	return &MessageHandler{
		sender: sender,
		logger: logger,
	}
}

func (h *MessageHandler) Handle(ctx context.Context, body []byte) error {
	var notification rmq.Notification

	if err := json.Unmarshal(body, &notification); err != nil {
		return err
	}

	h.logger.Info("processing notification",
		"event_id", notification.EventID,
		"user_id", notification.UserID)

	if err := h.sender.Send(ctx, notification); err != nil {
		return err
	}

	h.logger.Info("notification %s sent successfully to %s", notification.EventID, notification.UserID)

	return nil
}
