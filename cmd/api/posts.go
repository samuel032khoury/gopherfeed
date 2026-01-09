package main

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samuel032khoury/gopherfeed/internal/store"
)

type postKey string

const postKeyCtx postKey = "post"

// PostDTO represents the payload for creating or updating a post
//
//	@Description	Post creation/update payload
type PostDTO struct {
	Title   string   `json:"title" validate:"required,max=100" example:"My First Post"`
	Content string   `json:"content" validate:"required,max=2000" example:"This is the content of my post"`
	Tags    []string `json:"tags" example:"golang,api"`
}

// CreatePost godoc
//
//	@Summary		Create a post
//	@Description	Create a new post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		PostDTO						true	"Post payload"
//	@Success		201		{object}	DataResponse[store.Post]	"Post created successfully"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload PostDTO
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	userId := 1 // TODO: get authenticated user ID
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  int64(userId),
	}
	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	app.jsonResponse(w, post, http.StatusCreated)

}

// GetPost godoc
//
//	@Summary		Get a post
//	@Description	Get a post by its unique ID with comments
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int	true	"Post ID"
//	@Success		200		{object}	DataResponse[store.Post]
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments
	app.jsonResponse(w, post, http.StatusOK)
}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete a post by its unique ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int		true	"Post ID"
//	@Success		204		{string}	string	"Post deleted successfully"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := getPostFromContext(r).ID
	if err := app.store.Posts.Delete(r.Context(), postID); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Update a post by its unique ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int		true	"Post ID"
//	@Param			post	body		PostDTO	true	"Post payload"
//	@Success		200		{string}	string	"Post updated successfully"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse	"Edit conflict"
//	@Failure		500		{object}	ErrorResponse
//	@Security		ApiKeyAuth
//	@Router			/posts/{postID} [put]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	var payload PostDTO
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	post.Title = payload.Title
	post.Content = payload.Content
	post.Tags = payload.Tags
	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		if err == sql.ErrNoRows {
			app.conflictError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
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

func getPostFromContext(r *http.Request) *store.Post {
	post, ok := r.Context().Value(postKeyCtx).(*store.Post)
	if !ok {
		return nil
	}
	return post
}
