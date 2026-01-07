package main

import (
	"net/http"

	"github.com/samuel032khoury/gopherfeed/internal/store"
)

func (app *application) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	params := &store.PaginationParams{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Search: "",
		Tags:   []string{},
		Since:  "",
		Until:  "",
	}
	params, err := params.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(params); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	userId := 1 // TODO: get authenticated user ID
	feed, err := app.store.Posts.GetFeed(ctx, int64(userId), params)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, feed, http.StatusOK)
}
