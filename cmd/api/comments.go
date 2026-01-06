package main

import (
	"net/http"

	"github.com/samuel032khoury/gopherfeed/internal/store"
)

type CreateCommentDTO struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	postID := getPostFromContext(r).ID
	var payload CreateCommentDTO
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	userId := 1 // placeholder until we have authentication
	comment := &store.Comment{
		PostID:  postID,
		UserID:  int64(userId),
		Content: payload.Content,
	}
	ctx := r.Context()
	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, comment, http.StatusCreated)
}
