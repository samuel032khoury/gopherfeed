package main

import (
	"net/http"

	"github.com/samuel032khoury/gopherfeed/internal/store"
)

type userKey string

const userKeyCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Get user profile
//	@Description	Get a user by their unique ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	DataResponse[store.User]
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	app.jsonResponse(w, user, http.StatusOK)
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by their unique ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	nil	"User followed successfully"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse	"Unauthorized - login required"
//	@Failure		404		{object}	ErrorResponse	"User not found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followee := getUserFromContext(r)

	ctx := r.Context()
	currentUserID := getCurrentUserFromContext(r).ID

	if err := app.store.Followers.Follow(ctx, currentUserID, followee.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by their unique ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	nil	"User unfollowed successfully"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse	"Unauthorized - login required"
//	@Failure		404		{object}	ErrorResponse	"User not found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followee := getUserFromContext(r)

	currentUserID := getCurrentUserFromContext(r).ID

	if err := app.store.Followers.Unfollow(r.Context(), currentUserID, followee.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getCurrentUserFromContext(r *http.Request) *store.User {
	user, ok := r.Context().Value(currUserKeyCtx).(*store.User)
	if !ok {
		return nil
	}
	return user
}

func getUserFromContext(r *http.Request) *store.User {
	user, ok := r.Context().Value(userKeyCtx).(*store.User)
	if !ok {
		return nil
	}
	return user
}
