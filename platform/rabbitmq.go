package platform

import (
	"context"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	consumerWg sync.WaitGroup
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

func (q *RabbitMq) Close() {
	q.ch.Cancel("codern", false)
	q.conn.Close()
	q.consumerWg.Wait()
}

func (q *RabbitMq) Publish(exchange string, key string, body []byte) error {
	return q.ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	})
}

func (q *RabbitMq) Consume(queue string, fn func(amqp.Delivery)) error {
	q.consumerWg.Add(1)

	messages, err := q.ch.Consume(queue, "codern", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range messages {
			fn(delivery)
		}
		q.consumerWg.Done()
	}()

	return nil
}
