package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/samuel032khoury/gopherfeed/docs" // import docs
	"github.com/samuel032khoury/gopherfeed/internal/auth"
	"github.com/samuel032khoury/gopherfeed/internal/mq/publisher"
	"github.com/samuel032khoury/gopherfeed/internal/ratelimiter"
	"github.com/samuel032khoury/gopherfeed/internal/store"
	"github.com/samuel032khoury/gopherfeed/internal/store/cache"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config         config
	logger         *zap.SugaredLogger
	store          *store.Storage
	cacheStorage   *cache.CacheStorage
	emailPublisher *publisher.EmailPublisher
	authenticator  auth.Authenticator
	ratelimiter    ratelimiter.Limiter
}

type config struct {
	addr            string
	frontendBaseURL string
	db              dbConfig
	cache           cacheConfig
	mq              mqConfig
	auth            authConfig
	ratelimiter     ratelimiterConfig
	env             string
}

type dbConfig struct {
	url          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type cacheConfig struct {
	redisAddr     string
	redisPassword string
	redisDB       int
	enabled       bool
}

type mqConfig struct {
	url   string
	names queueNames
}

type queueNames struct {
	email string
}

type authConfig struct {
	basic basicAuthConfig
	jwt   jwtConfig
}

type basicAuthConfig struct {
	username string
	password string
}

type jwtConfig struct {
	secretKey     string
	tokenDuration string
	iss           string
	aud           string
}

type ratelimiterConfig struct {
	quota    int
	interval string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.config.frontendBaseURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // Allow cookies to be sent
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(app.RateLimitMiddleware)

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware).Get("/health", app.healthCheckHandler)

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://"+app.config.addr+"/v1/swagger/doc.json"),
		))
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.TokenAuthMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.PostParamMiddleware)
				r.Get("/", app.getPostHandler)
				r.Post("/comments", app.createCommentHandler)
				r.With(app.RBACMiddleware("moderator")).Put("/", app.updatePostHandler)
				r.With(app.RBACMiddleware("admin")).Delete("/", app.deletePostHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.UserParamMiddleware)
				r.Get("/", app.getUserHandler)
				r.Group(func(r chi.Router) {
					r.Use(app.TokenAuthMiddleware)
					r.Put("/follow", app.followUserHandler)
					r.Put("/unfollow", app.unfollowUserHandler)
				})
			})
		})

		r.Route("/feeds", func(r chi.Router) {
			r.Use(app.TokenAuthMiddleware)
			r.Get("/", app.getFeedHandler)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
			r.Post("/login", app.loginUserHandler)
			r.Post("/activate", app.activateUserHandler)
			r.Post("/logout", app.logoutUserHandler)
		})
	})
	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.addr
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	shutdown := make(chan error, 1)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		app.logger.Infow("Received signal, initiating shutdown...", "signal", sig)
		shutdown <- srv.Shutdown(ctx)
	}()
	app.logger.Infow("server has started", "address", app.config.addr, "env", app.config.env)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return <-shutdown
}
