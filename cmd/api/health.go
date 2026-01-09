package main

import "net/http"

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Check the health of the application
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	DataResponse[map[string]string]
//	@Failure		500	{object}	ErrorResponse
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "OK",
		"env":     app.config.env,
		"version": version,
	}
	app.jsonResponse(w, data, http.StatusOK)
}
