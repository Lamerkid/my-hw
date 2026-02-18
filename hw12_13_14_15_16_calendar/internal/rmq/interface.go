package rmq

import "context"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type MessageHandler interface {
	Handle(ctx context.Context, msg []byte) error
}

type Producer interface {
	Publish(ctx context.Context, msg Notification) error
}

type Consumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
}
