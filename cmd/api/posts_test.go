package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/samuel032khoury/gopherfeed/internal/auth"
	"github.com/samuel032khoury/gopherfeed/internal/store"
)

func TestGetFeed(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/feeds", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := execRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		// Create a test user
		user := &store.User{
			ID: 1,
		}

		// Generate JWT token
		authenticator := auth.NewJWTAuthenticator(
			"test-secret-key",
			"24h",
			"gopherfeed-api",
			"gopherfeed",
		)

		exp, iss, aud := authenticator.GetMetadata()
		claims := jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(exp).Unix(),
			"iat": time.Now().Unix(),
			"nbf": time.Now().Unix(),
			"iss": iss,
			"aud": aud,
		}

		token, err := authenticator.GenerateToken(claims)
		if err != nil {
			t.Fatal(err)
		}

		// Create request with JWT cookie
		req, err := http.NewRequest(http.MethodGet, "/v1/feeds", nil)
		if err != nil {
			t.Fatal(err)
		}

		cookie := &http.Cookie{
			Name:  "jwt",
			Value: token,
		}
		req.AddCookie(cookie)

		// Execute request
		rr := execRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
