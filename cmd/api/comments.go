package main

import (
	"net/http"

	"github.com/samuel032khoury/gopherfeed/internal/store"
)

// CreateCommentDTO represents the payload for creating a comment
//
//	@Description	Comment creation payload
type CreateCommentDTO struct {
	Content string `json:"content" validate:"required,max=1000" example:"Great post!"`
}

// CreateComment godoc
//
//	@Summary		Create a comment
//	@Description	Create a new comment on a post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int					true	"Post ID"
//	@Param			comment	body		CreateCommentDTO	true	"Comment payload"
//	@Success		201		{object}	DataResponse[store.Comment]
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID}/comments [post]
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
