package rmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQP struct {
	url       string
	topicName string
	logger    Logger
	conn      *amqp.Connection
	ch        *amqp.Channel
}

func NewAMQP(logger Logger, url string, topicName string) *AMQP {
	return &AMQP{
		url:       url,
		topicName: topicName,
		logger:    logger,
	}
}

func (a *AMQP) Connect() (err error) {
	a.conn, err = amqp.Dial(a.url)
	if err != nil {
		return err
	}

	a.ch, err = a.conn.Channel()
	if err != nil {
		a.conn.Close()
		return err
	}

	err = a.ch.ExchangeDeclare(
		a.topicName,
		"topic",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		a.conn.Close()
		return err
	}

	queue, err := a.ch.QueueDeclare(
		a.topicName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		a.conn.Close()
		return err
	}

	err = a.ch.QueueBind(
		queue.Name,
		"#",
		a.topicName,
		false,
		nil,
	)
	if err != nil {
		a.conn.Close()
		return err
	}

	return nil
}

func (a *AMQP) Close() error {
	if a.ch != nil {
		a.ch.Close()
	}
	if a.conn != nil {
		return a.conn.Close()
	}
	return nil
}

func (a *AMQP) Publish(ctx context.Context, msg Notification) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return a.ch.PublishWithContext(ctx,
		a.topicName,
		"notification",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

func (a *AMQP) Consume(ctx context.Context, handler MessageHandler) error {
	msgs, err := a.ch.Consume(
		a.topicName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case amqpMsg, ok := <-msgs:
			if !ok {
				return nil
			}

			a.logger.Debug("received message: %s", amqpMsg.Body)

			if err := handler.Handle(ctx, amqpMsg.Body); err != nil {
				a.logger.Error("failed to handle message: %v", err)
				continue
			}
		}
	}
}
