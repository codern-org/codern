package platform

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMq(url string) (*RabbitMq, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMq{
		conn: conn,
		ch:   ch,
	}, nil
}

func (q *RabbitMq) Publish(
	ctx context.Context,
	exchange string,
	key string,
	mandatory bool,
	immediate bool,
	msg amqp.Publishing,
) error {
	return q.ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
}

func (q *RabbitMq) Close() {
	q.conn.Close()
	q.ch.Close()
}
