package publisher

import (
	"context"
	"fmt"

	"github.com/samuel032khoury/gopherfeed/internal/email"
	"github.com/samuel032khoury/gopherfeed/internal/logger"
	"github.com/samuel032khoury/gopherfeed/internal/mq"
)

// EmailPublisher handles publishing email messages to the queue.
// It wraps the RabbitMQ connection and provides a simple API for sending emails.
type EmailPublisher struct {
	queue *mq.RabbitMQ
}

// New creates a new EmailPublisher with the given RabbitMQ connection.
func NewEmailPublisher(url, queueName string, log logger.Logger) (*EmailPublisher, error) {
	queue, err := mq.New(url, queueName, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ connection: %w", err)
	}

	return &EmailPublisher{
		queue: queue,
	}, nil
}

func (p *EmailPublisher) Publish(to string, templatePath string, data any) error {
	ctx := context.Background()
	msg, err := email.New(to, templatePath, data)
	if err != nil {
		return fmt.Errorf("failed to create email message: %w", err)
	}
	body, err := msg.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to serialize email message: %w", err)
	}
	if err := p.queue.Publish(ctx, body); err != nil {
		return fmt.Errorf("failed to publish email message: %w", err)
	}
	return nil
}

// Close closes the underlying RabbitMQ connection.
func (p *EmailPublisher) Close() error {
	return p.queue.Close()
}
