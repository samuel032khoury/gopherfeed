package main

import (
	"log"

	"github.com/samuel032khoury/gopherfeed/internal/db"
	"github.com/samuel032khoury/gopherfeed/internal/env"
	asyncMailer "github.com/samuel032khoury/gopherfeed/internal/mailer/asyncMailer"
	"github.com/samuel032khoury/gopherfeed/internal/mq"
	"github.com/samuel032khoury/gopherfeed/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GopherFeed API
//	@description	This is the GopherFeed API server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	cfg := config{
		addr:            env.GetString("ADDR", ":8080"),
		frontendBaseURL: env.GetString("FRONTEND_BASE_URL", "localhost:3000"),
		db: dbConfig{
			url:          env.GetString("DB_URL", "postgres://user:password@localhost:5432/gopherfeed?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		rabbitmq: rabbitmqConfig{
			url:       env.GetString("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			queueName: env.GetString("RABBITMQ_EMAIL_QUEUE", "email_queue"),
		},
		env: env.GetString("ENV", "development"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.url,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")
	store := store.NewPostgresStorage(db)
	mailClient, err := asyncMailer.New(mq.Config{
		URL:       cfg.rabbitmq.url,
		QueueName: cfg.rabbitmq.queueName,
	})
	if err != nil {
		logger.Fatal("unable to set up async mailer: ", err)
	}
	defer mailClient.Close()
	logger.Info("async mailer configured (RabbitMQ)")

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailClient,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
