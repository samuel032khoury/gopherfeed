package main

import "net/http"

type healthResponse struct {
	Status  string `json:"status" example:"OK"`
	Env     string `json:"env" example:"development"`
	Version string `json:"version" example:"1.0.0"`
}

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Check the health of the application
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	DataResponse[healthResponse]	"Application is healthy"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := healthResponse{
		Status:  "OK",
		Env:     app.config.env,
		Version: version,
	}
	app.jsonResponse(w, data, http.StatusOK)
}
