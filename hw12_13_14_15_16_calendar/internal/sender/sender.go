package sender

import (
	"context"

	"github.com/lamerkid/my-hw/hw12_13_14_15_calendar/internal/rmq"
)

type Sender struct {
	consumer rmq.Consumer
	logger   Logger
}

func NewSender(consumer rmq.Consumer, logger Logger) *Sender {
	return &Sender{
		consumer: consumer,
		logger:   logger,
	}
}

func (s *Sender) Send(ctx context.Context, msg rmq.Notification) error {
	return nil
}
