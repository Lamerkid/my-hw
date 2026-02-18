package rmq

import (
	"context"
	"encoding/json"
	"sync"
)

type MockBroker struct {
	mu       sync.Mutex
	Messages []Notification
}

func NewMockBroker() *MockBroker {
	return &MockBroker{
		Messages: make([]Notification, 0),
	}
}

func (m *MockBroker) Publish(ctx context.Context, msg Notification) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		m.Messages = append(m.Messages, msg)
	}

	return nil
}

func (m *MockBroker) Consume(ctx context.Context, handler MessageHandler) error {
	go func() {
		for i := 0; i < len(m.Messages); i++ {
			select {
			case <-ctx.Done():
				return
			default:
				data, _ := json.Marshal(m.Messages[i])
				handler(ctx, data)
			}
		}
	}()

	return nil
}
