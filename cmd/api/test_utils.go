package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samuel032khoury/gopherfeed/internal/auth"
	"github.com/samuel032khoury/gopherfeed/internal/ratelimiter"
	"github.com/samuel032khoury/gopherfeed/internal/store"
	"github.com/samuel032khoury/gopherfeed/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()
	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockCacheStorage(nil)
	mockRatelimiter := ratelimiter.NewMockRateLimiter()
	mockAuthenticator := auth.NewJWTAuthenticator(
		"test-secret-key",
		"24h",
		"gopherfeed-api",
		"gopherfeed",
	)
	testConfig := config{
		addr:            ":8080",
		frontendBaseURL: "http://localhost:5173",
		cache: cacheConfig{
			enabled: false,
		},
		env: "test",
	}
	return &application{
		config:        testConfig,
		logger:        logger,
		store:         &mockStore,
		cacheStorage:  mockCache,
		ratelimiter:   mockRatelimiter,
		authenticator: mockAuthenticator,
	}
}

func execRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected status %d; got %d", expected, actual)
	}
}
