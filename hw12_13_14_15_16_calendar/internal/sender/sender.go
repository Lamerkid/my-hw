package sender

import (
	"context"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/rmq"
)

type Sender struct {
	consumer rmq.Consumer
	handler  *MessageHandler
	logger   Logger
}

func NewSender(logger Logger, consumer rmq.Consumer) *Sender {
	return &Sender{
		consumer: consumer,
		logger:   logger,
	}
}

func (s *Sender) SendMessage(ctx context.Context, msg rmq.Notification) error {
	s.logger.Debug("sending notification to: %s", msg.UserID)

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		s.logger.Info("Notification to %s on upcoming event: %s, on %s", msg.UserID, msg.EventTitle, msg.EventTime.String())
	}

	return s.consumer.Consume(ctx, s.handler)
}
