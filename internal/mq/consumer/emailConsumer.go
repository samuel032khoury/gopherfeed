package consumer

import (
	"context"
	"encoding/json"

	"github.com/samuel032khoury/gopherfeed/internal/email"
	"github.com/samuel032khoury/gopherfeed/internal/logger"
	"github.com/samuel032khoury/gopherfeed/internal/mq"
)

type EmailConsumer struct {
	mq     *mq.RabbitMQ
	sender email.Sender
	logger logger.Logger
}

func NewEmailConsumer(mq *mq.RabbitMQ, sender email.Sender, log logger.Logger) *EmailConsumer {
	if log == nil {
		log = logger.NewNoopLogger()
	}
	return &EmailConsumer{
		mq:     mq,
		sender: sender, logger: log}
}

func (ec *EmailConsumer) Start(ctx context.Context) error {
	msgs, err := ec.mq.Consume()
	if err != nil {
		return err
	}

	ec.logger.Info("Email worker started, waiting for messages...")
	for {
		select {
		case <-ctx.Done():
			ec.logger.Info("Email worker shutting down...")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				ec.logger.Info("Message channel closed")
				return nil
			}
			if err := ec.processMessage(msg.Body); err != nil {
				ec.logger.Errorw("Failed to process message", "error", err)
				msg.Nack(false, false)
				continue
			}
			if err := msg.Ack(false); err != nil {
				ec.logger.Errorw("Failed to acknowledge message", "error", err)
			}
		}
	}
}

func (ec *EmailConsumer) processMessage(body []byte) error {
	emailMsg, err := email.FromBytes(body)
	if err != nil {
		return err
	}
	ec.logger.Infow("Processing email", "template", emailMsg.TemplatePath, "recipient", emailMsg.To)
	var data any
	if err := json.Unmarshal(emailMsg.Data, &data); err != nil {
		return err
	}
	if err := ec.sender.Send(
		emailMsg.To,
		emailMsg.TemplatePath,
		data,
	); err != nil {
		return err
	}
	ec.logger.Infow("Email sent successfully", "template", emailMsg.TemplatePath, "recipient", emailMsg.To)
	return nil
}

func (ec *EmailConsumer) Close() error {
	return ec.mq.Close()
}
