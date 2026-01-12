package main

import (
	"net/http"

	"github.com/samuel032khoury/gopherfeed/internal/store"
)

// GetFeed godoc
//
//	@Summary		Get user feed
//	@Description	Get posts feed with filtering and pagination
//	@Tags			feeds
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Number of items per page (1-100)"	example(20)
//	@Param			offset	query		int		false	"Number of items to skip"			example(0)
//	@Param			sort	query		string	false	"Sort order"						Enums(asc, desc)	example(desc)
//	@Param			tags	query		string	false	"Comma-separated tags filter"		example("golang,api")
//	@Param			search	query		string	false	"Search in title and content"		example("golang")
//	@Param			since	query		string	false	"Posts since this date (RFC3339)"	example("2026-01-01T00:00:00Z")
//	@Param			until	query		string	false	"Posts until this date (RFC3339)"	example("2026-12-31T23:59:59Z")
//	@Success		200		{object}	DataResponse[[]store.FeedablePost]
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse	"Unauthorized - login required"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/feeds [get]
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
	currentUserID := getCurrentUserFromContext(r).ID
	feed, err := app.store.Posts.GetFeed(ctx, int64(currentUserID), params)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, feed, http.StatusOK)
}
