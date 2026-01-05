package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/samuel032khoury/gopherfeed/internal/store"
)

type CreatePostDTO struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags,omitempty"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostDTO
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	userId := 1 // placeholder until we have authentication
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  int64(userId),
	}
	ctx := r.Context()
	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, "failed to create post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := writeJSON(w, post, http.StatusCreated); err != nil {
		writeJSONError(w, "failed to write response", http.StatusInternalServerError)
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	postIDParam := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDParam, 10, 64)
	if err != nil {
		writeJSONError(w, "invalid post ID", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	post, err := app.store.Posts.GetByID(ctx, postID)
	if err != nil {
		writeJSONError(w, "failed to get post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if post == nil {
		writeJSONError(w, "post not found", http.StatusNotFound)
		return
	}
	if err := writeJSON(w, post, http.StatusOK); err != nil {
		writeJSONError(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}
