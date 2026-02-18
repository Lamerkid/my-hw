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

func (s *Sender) Start(ctx context.Context) error {
	s.logger.Info("sender servise starting")

	consumeHandler := rmq.MessageHandler(func(ctx context.Context, data []byte) error {
		return s.handler.Handle(ctx, data)
	})

	if err := s.consumer.Consume(ctx, consumeHandler); err != nil {
		s.logger.Error("consumer stopped: %v", err)
	}

	return nil
}
