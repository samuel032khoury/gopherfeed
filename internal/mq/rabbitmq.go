package mq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueryTimeoutDuration = 5 * time.Second
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type Config struct {
	URL string
	QueueName string
}

// New creates a new RabbitMQ instance.
// It returns a configured RabbitMQ instance ready for publishing or consuming.
func New(cfg Config) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	queue, err := ch.QueueDeclare(cfg.QueueName, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	log.Printf("Connected to RabbitMQ, queue '%s' ready", cfg.QueueName)
	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		queue:   queue,
	}, nil
}

func (r *RabbitMQ) Publish(ctx context.Context, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := r.channel.PublishWithContext(ctx, "", r.queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func (r *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	err := r.channel.Qos(1,0,false)
	if err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}
	msgs, err := r.channel.Consume(r.queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}
	return msgs, nil
}

// Close gracefully closes the RabbitMQ connection.
func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	log.Println("RabbitMQ connection closed")
	return nil
}