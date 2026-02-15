package rmq

import "context"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Producer interface {
	Publish(ctx context.Context, msg Notification)
	Close() error
}

type Consumer interface {
	Consume(ctx context.Context)
	Close() error
}
