package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/samuel032khoury/gopherfeed/docs" // import docs
	"github.com/samuel032khoury/gopherfeed/internal/auth"
	"github.com/samuel032khoury/gopherfeed/internal/mq/publisher"
	"github.com/samuel032khoury/gopherfeed/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config         config
	store          *store.Storage
	logger         *zap.SugaredLogger
	emailPublisher *publisher.EmailPublisher
	authenticator  auth.Authenticator
}

type config struct {
	addr            string
	frontendBaseURL string
	db              dbConfig
	mq              mqConfig
	auth            authConfig
	env             string
}

type dbConfig struct {
	url          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type mqConfig struct {
	url   string
	names queueNames
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
}

type queueNames struct {
	email string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://"+app.config.addr+"/v1/swagger/doc.json"),
		))
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Post("/comments", app.createCommentHandler)
				r.Delete("/", app.deletePostHandler)
				r.Put("/", app.updatePostHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
		})

		r.Route("/feeds", func(r chi.Router) {
			r.Get("/", app.getFeedHandler)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
			r.Post("/login", app.loginUserHandler)
			r.Post("/activate", app.activateUserHandler)
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
	app.logger.Infow("server has started", "address", app.config.addr, "env", app.config.env)
	return srv.ListenAndServe()
}
