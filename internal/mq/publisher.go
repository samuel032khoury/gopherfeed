package mq

import (
	"context"
	"fmt"
	"log"
)

// EmailPublisher handles publishing email messages to the queue.
// It wraps the RabbitMQ connection and provides a simple API for sending emails.
type EmailPublisher struct {
	mq *RabbitMQ
}

// NewEmailPublisher creates a new EmailPublisher with the given RabbitMQ connection.
func NewEmailPublisher(mq *RabbitMQ) *EmailPublisher {
	return &EmailPublisher{
		mq: mq,
	}
}

func (p *EmailPublisher) PublishEmail(ctx context.Context, templateFile, username, email string, data any) error {
	msg, err := NewEmailMessage(templateFile, username, email, data)
	if err != nil {
		return fmt.Errorf("failed to create email message: %w", err)
	}
	body, err := msg.ToBytes()
	if err != nil {
		return fmt.Errorf("failed to serialize email message: %w", err)
	}
	if err := p.mq.Publish(ctx, body); err != nil {
		return fmt.Errorf("failed to publish email message: %w", err)
	}
	log.Printf("Email task published: %s", msg)
	return nil
}

// Close closes the underlying RabbitMQ connection.
func (p *EmailPublisher) Close() error {
	return p.mq.Close()
}
