package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samuel032khoury/gopherfeed/internal/store"
)

type currUserKey string

const currUserKeyCtx currUserKey = "curr_user"

func (app *application) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allow, retryAfter := app.ratelimiter.Allow(r.RemoteAddr); !allow {
			app.rateLimitExceededError(w, r, strconv.Itoa(int(retryAfter.Seconds())))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) PostParamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(postIDParam, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		ctx := r.Context()
		post, err := app.store.Posts.GetByID(ctx, postID)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		if post == nil {
			app.notFoundError(w, r)
			return
		}
		ctx = context.WithValue(ctx, postKeyCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) UserParamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.getUser(ctx, userID)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		if user == nil {
			app.notFoundError(w, r)
			return
		}
		ctx = context.WithValue(ctx, userKeyCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read the authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, true, fmt.Errorf("missing authorization header"))
			return
		}
		// parse -> get the base64 encoded username:password
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Basic" {
			app.unauthorizedError(w, r, true, fmt.Errorf("malformed authorization header"))
			return
		}
		// decode
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			app.unauthorizedError(w, r, true, fmt.Errorf("invalid base64 encoding"))
			return
		}
		// split into username and password
		username := app.config.auth.basic.username
		password := app.config.auth.basic.password
		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 || creds[0] != username || creds[1] != password {
			app.unauthorizedError(w, r, true, fmt.Errorf("invalid authorization value"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get JWT token from cookie
		cookie, err := r.Cookie("jwt")
		if err != nil {
			app.unauthorizedError(w, r, false, fmt.Errorf("missing authentication cookie"))
			return
		}
		token := cookie.Value
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, false, fmt.Errorf("invalid or expired token"))
			return
		}
		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok || !jwtToken.Valid {
			app.unauthorizedError(w, r, false, fmt.Errorf("invalid token claims"))
			return
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil || userID <= 0 {
			app.unauthorizedError(w, r, false, fmt.Errorf("invalid user ID in token claims"))
			return
		}

		ctx := r.Context()
		user, err := app.getUser(ctx, userID)
		if err != nil {
			app.unauthorizedError(w, r, false, fmt.Errorf("user not found"))
			return
		}
		ctx = context.WithValue(ctx, currUserKeyCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) RBACMiddleware(requiredRole string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := getCurrentUserFromContext(r)
			post := getPostFromContext(r)
			if post.UserID != user.ID {
				if allowed, err := app.checkRolePermissions(r.Context(), user.RoleID, requiredRole); err != nil {
					app.internalServerError(w, r, err)
					return
				} else if !allowed {
					app.forbiddenError(w, r)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) checkRolePermissions(ctx context.Context, userRoleID int64, requiredRoleName string) (bool, error) {
	requiredRole, err := app.store.Roles.GetByName(ctx, requiredRoleName)
	if err != nil {
		return false, err
	}
	if requiredRole == nil {
		return false, fmt.Errorf("required role not found")
	}
	userRole, err := app.store.Roles.GetByID(ctx, userRoleID)
	if err != nil {
		return false, err
	}
	return userRole.Level >= requiredRole.Level, nil
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	var user *store.User
	var err error
	if app.config.cache.enabled {
		user, err = app.cacheStorage.Users.Get(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	if user == nil {
		if user, err = app.store.Users.GetByID(ctx, userID); user == nil {
			return user, err
		}
		if app.config.cache.enabled {
			if err = app.cacheStorage.Users.Set(ctx, user); err != nil {
				return nil, err
			}
		}
	}
	return user, nil
}
