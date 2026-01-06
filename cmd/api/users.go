package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samuel032khoury/gopherfeed/internal/store"
)

type userKey string

const userKeyCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	app.jsonResponse(w, user, http.StatusOK)
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	var userId int64 = 1 // TODO: get authenticated user ID

	if err := app.store.Followers.Follow(r.Context(), userId, followerUser.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	var userId int64 = 1 // TODO: get authenticated user ID

	if err := app.store.Followers.Unfollow(r.Context(), userId, followerUser.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)
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

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userKeyCtx).(*store.User)
	return user
}
