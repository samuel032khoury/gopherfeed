package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
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
	app.jsonResponse(w, user, http.StatusOK)
}
