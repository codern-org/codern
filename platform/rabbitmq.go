package platform

import (
	"context"
	"encoding/json"

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

	_, err = ch.QueueDeclare(
		"grading",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
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
	message interface{},
) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return q.ch.PublishWithContext(ctx, exchange, key, mandatory, immediate, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (q *RabbitMq) Close() {
	q.conn.Close()
	q.ch.Close()
}
