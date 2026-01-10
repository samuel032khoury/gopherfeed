package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/samuel032khoury/gopherfeed/internal/env"
	"github.com/samuel032khoury/gopherfeed/internal/mailer"
	"github.com/samuel032khoury/gopherfeed/internal/mq"
)

func main() {
	log.Println("Starting email worker...")
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
	mailtrap, err := mailer.NewMailtrap(
		mailConfig.fromEmail,
		mailConfig.host,
		mailConfig.username,
		mailConfig.password,
		mailConfig.port,
	)
	if err != nil {
		log.Fatal("Failed to create mailer:", err)
	}
	rabbitmq, err := mq.New(mq.Config{
		URL:       mqConfig.url,
		QueueName: mqConfig.queueName,
	})
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitmq.Close()
	log.Println("Connected to RabbitMQ")
	consumer := mq.NewEmailConsumer(rabbitmq, mailtrap)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v, initiating shutdown...", sig)
		cancel()
	}()
	if err := consumer.Start(ctx); err != nil && err != context.Canceled {
		log.Fatal("Worker error:", err)
	}

	log.Println("Email worker stopped")
}
