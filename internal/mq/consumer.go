package mq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/samuel032khoury/gopherfeed/internal/mailer"
)

type EmailConsumer struct {
	mq     *RabbitMQ
	mailer mailer.Client
}

func NewEmailConsumer(mq *RabbitMQ, mailer mailer.Client) *EmailConsumer {
	return &EmailConsumer{
		mq:     mq,
		mailer: mailer,
	}
}

func (ec *EmailConsumer) Start(ctx context.Context) error {
	msgs, err := ec.mq.Consume()
	if err != nil {
		return err
	}

	log.Println("Email worker started, waiting for messages...")
	for {
		select {
		case <-ctx.Done():
			log.Println("Email worker shutting down...")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				log.Println("Message channel closed")
				return nil
			}
			if err := ec.processMessage(msg.Body); err != nil {
				log.Printf("Failed to process message: %v", err)
				msg.Nack(false, false)
				continue
			}
			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			}
		}
	}
}

func (ec *EmailConsumer) processMessage(body []byte) error {
	emailMsg, err := FromBytes(body)
	if err != nil {
		return err
	}
	log.Printf("Processing email: template=%s, recipient=%s", emailMsg.TemplateFile, emailMsg.Email)
	var data any
	if err := json.Unmarshal(emailMsg.Data, &data); err != nil {
		return err
	}
	if err := ec.mailer.Send(
		emailMsg.TemplateFile,
		emailMsg.Username,
		emailMsg.Email,
		data,
	); err != nil {
		return err
	}
	log.Printf("Email sent successfully: template=%s, recipient=%s", emailMsg.TemplateFile, emailMsg.Email)
	return nil
}

func (ec *EmailConsumer) Close() error {
	return ec.mq.Close()
}
