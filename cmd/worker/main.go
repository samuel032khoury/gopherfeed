package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/samuel032khoury/gopherfeed/internal/email"
	"github.com/samuel032khoury/gopherfeed/internal/env"
	"github.com/samuel032khoury/gopherfeed/internal/mq"
	"github.com/samuel032khoury/gopherfeed/internal/mq/consumer"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	logger.Info("Starting email worker...")
	mqConfig := rabbitmqConfig{
		url:       env.GetString("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		queueName: env.GetString("RABBITMQ_EMAIL_QUEUE", "email_queue"),
	}
	mailConfig := mailConfig{
		fromEmail: env.GetString("MAIL_FROM_EMAIL", "comm@gopherfeed.io"),
		host:      env.GetString("MAIL_HOST", "sandbox.smtp.mailtrap.io"),
		port:      env.GetInt("MAIL_PORT", 587),
		username:  env.GetString("MAIL_USERNAME", ""),
		password:  env.GetString("MAIL_PASSWORD", ""),
	}
	sender, err := email.NewMailtrap(
		mailConfig.fromEmail,
		mailConfig.host,
		mailConfig.username,
		mailConfig.password,
		mailConfig.port,
		logger,
	)
	if err != nil {
		logger.Fatal("Failed to create email sender:", err)
	}
	rabbitmq, err := mq.New(mqConfig.url, mqConfig.queueName, logger)
	if err != nil {
		logger.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitmq.Close()
	logger.Info("Connected to RabbitMQ")
	consumer := consumer.NewEmailConsumer(rabbitmq, sender, logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Infow("Received signal, initiating shutdown...", "signal", sig)
		cancel()
	}()
	if err := consumer.Start(ctx); err != nil && err != context.Canceled {
		logger.Fatal("Worker error:", err)
	}

	logger.Info("Email worker stopped")
}
