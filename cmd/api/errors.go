package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, "the server encountered an error", http.StatusInternalServerError)
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, err.Error(), http.StatusBadRequest)
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("not found", "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, "resource not found", http.StatusNotFound)
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, "edit conflict occurred", http.StatusConflict)
}

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, isBasicAuth bool, err error) {
	app.logger.Warnw("unauthorized", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	if isBasicAuth {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted" charset="UTF-8"`)
	}
	writeJSONError(w, "unauthorized", http.StatusUnauthorized)
}
func (app *application) forbiddenError(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, "forbidden", http.StatusForbidden)
}
