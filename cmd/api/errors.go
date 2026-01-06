package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, "the server encountered an error", http.StatusInternalServerError)
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, err.Error(), http.StatusBadRequest)
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request) {
	log.Printf("not found: %s, path: %s", r.Method, r.URL.Path)
	writeJSONError(w, "resource not found", http.StatusNotFound)
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: %s, path: %s, error: %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, "edit conflict occurred", http.StatusConflict)
}
