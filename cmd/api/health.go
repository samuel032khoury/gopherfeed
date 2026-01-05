package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "OK",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(w, data, http.StatusOK); err != nil {
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
	}
}
