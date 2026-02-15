package rmq

import (
	"context"
	"encoding/json"

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

func (a *AMQP) Connect() error {
	var err error
	a.conn, err = amqp.Dial(a.url)
	if err != nil {
		return err
	}

	a.ch, err = a.conn.Channel()
	if err != nil {
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
		"",
		a.topicName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (a *AMQP) Consume(ctx context.Context) error {
	msgs, err := a.ch.Consume(
		"test",
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

			var msg string
			if err := json.Unmarshal(amqpMsg.Body, &msg); err != nil {
				a.logger.Error("error unmarshaling: %s", amqpMsg.Body)
				return err
			}
		}
	}
}
