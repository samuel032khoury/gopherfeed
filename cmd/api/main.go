package main

import (
	"expvar"
	"log"
	"runtime"

	"github.com/samuel032khoury/gopherfeed/internal/auth"
	"github.com/samuel032khoury/gopherfeed/internal/db"
	"github.com/samuel032khoury/gopherfeed/internal/env"
	"github.com/samuel032khoury/gopherfeed/internal/mq/publisher"
	"github.com/samuel032khoury/gopherfeed/internal/ratelimiter"
	"github.com/samuel032khoury/gopherfeed/internal/store"
	"github.com/samuel032khoury/gopherfeed/internal/store/cache"
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
func main() {
	// =========================================================================
	// Configuration
	// =========================================================================
	cfg := loadConfig()

	// =========================================================================
	// Logger
	// =========================================================================
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// =========================================================================
	// Database & Storage
	// =========================================================================
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

	// =========================================================================
	// Cache (Optional)
	// =========================================================================
	var cacheStorage *cache.CacheStorage
	if cfg.cache.enabled {
		redisClient := cache.NewRedisClient(
			cfg.cache.redisAddr,
			cfg.cache.redisPassword,
			cfg.cache.redisDB,
		)
		cacheStorage = cache.NewRedisStorage(redisClient)
		logger.Info("redis cache client created")
	} else {
		logger.Info("redis cache is disabled")
	}

	// =========================================================================
	// Message Queue
	// =========================================================================
	emailPublisher, err := publisher.NewEmailPublisher(
		cfg.mq.url,
		cfg.mq.names.email,
		logger,
	)
	if err != nil {
		log.Fatal("failed to create email publisher:", err)
	}
	defer emailPublisher.Close()
	logger.Info("email publisher created")

	// =========================================================================
	// Authentication
	// =========================================================================
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.jwt.secretKey,
		cfg.auth.jwt.tokenDuration,
		cfg.auth.jwt.iss,
		cfg.auth.jwt.aud,
	)

	// =========================================================================
	// Rate Limiter
	// =========================================================================
	limiter, err := ratelimiter.NewFixedWindowLimiter(
		cfg.ratelimiter.quota,
		cfg.ratelimiter.interval,
	)
	if err != nil {
		logger.Fatal("failed to create rate limiter:", err)
	}

	// =========================================================================
	// Application
	// =========================================================================
	app := &application{
		config:         cfg,
		store:          store,
		cacheStorage:   cacheStorage,
		logger:         logger,
		emailPublisher: emailPublisher,
		authenticator:  jwtAuthenticator,
		ratelimiter:    limiter,
	}

	// =========================================================================
	// Stats
	// =========================================================================
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// =========================================================================
	// Server
	// =========================================================================
	mux := app.mount()
	app.run(mux)
}

// loadConfig loads application configuration from environment variables
func loadConfig() config {
	return config{
		addr:            env.GetString("ADDR", ":8080"),
		frontendBaseURL: env.GetString("FRONTEND_BASE_URL", "localhost:5173"),
		db: dbConfig{
			url:          env.GetString("DB_URL", "postgres://user:password@localhost:5432/gopherfeed?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		mq: mqConfig{
			url: env.GetString("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			names: queueNames{
				email: env.GetString("RABBITMQ_EMAIL_QUEUE", "email_queue"),
			},
		},
		auth: authConfig{
			basic: basicAuthConfig{
				username: env.GetString("BASIC_AUTH_USERNAME", "admin"),
				password: env.GetString("BASIC_AUTH_PASSWORD", "password"),
			},
			jwt: jwtConfig{
				secretKey:     env.GetString("JWT_SECRET_KEY", "your-secret-key"),
				tokenDuration: env.GetString("JWT_TOKEN_DURATION", "24h"),
				iss:           env.GetString("JWT_ISSUER", "gopherfeed-api"),
				aud:           env.GetString("JWT_AUDIENCE", "gopherfeed"),
			},
		},
		cache: cacheConfig{
			redisAddr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			redisPassword: env.GetString("REDIS_PASSWORD", ""),
			redisDB:       env.GetInt("REDIS_DB", 0),
			enabled:       env.GetBool("REDIS_ENABLED", false),
		},
		ratelimiter: ratelimiterConfig{
			quota:    env.GetInt("RATE_LIMITER_QUOTA", 100),
			interval: env.GetString("RATE_LIMITER_INTERVAL", "5s"),
		},
		env: env.GetString("ENV", "development"),
	}
}
