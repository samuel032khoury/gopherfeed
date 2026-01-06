package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		writeJSONError(w, "the server encountered an error", http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	data := errorResponse{
		Error: message,
	}
	writeJSON(w, &data, status)
}

func (app *application) jsonResponse(w http.ResponseWriter, data any, status int) {
	type jsonResponse struct {
		Data any `json:"data"`
	}
	response := &jsonResponse{
		Data: data,
	}
	writeJSON(w, response, status)
}
