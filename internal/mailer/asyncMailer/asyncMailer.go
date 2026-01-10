package asyncMailer

import (
	"context"

	"github.com/samuel032khoury/gopherfeed/internal/mq"
)

type AsyncClient struct {
	publisher *mq.EmailPublisher
}

func New(mqConfig mq.Config) (*AsyncClient, error) {
	queue, err := mq.New(mqConfig)
	if err != nil {
		return nil, err
	}
	publisher := mq.NewEmailPublisher(queue)
	return &AsyncClient{
		publisher: publisher,
	}, nil
}

func (c *AsyncClient) Send(templateFile, username, email string, data any) error {
	ctx := context.Background()
	return c.publisher.PublishEmail(ctx, templateFile, username, email, data)
}

func (c *AsyncClient) Close() error {
	return c.publisher.Close()
}
