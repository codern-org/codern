package platform

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	conn *amqp.Connection
	Ch   *amqp.Channel
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
		Ch:   ch,
	}, nil
}

func (q *RabbitMq) Close() {
	q.conn.Close()
	q.Ch.Close()
}
