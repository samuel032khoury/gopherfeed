package main

import "net/http"

func (app *application) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := 1 // TODO: get authenticated user ID
	feed, err := app.store.Posts.GetFeed(ctx, int64(userId))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, feed, http.StatusOK)
}
