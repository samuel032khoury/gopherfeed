package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samuel032khoury/gopherfeed/internal/email"
	"github.com/samuel032khoury/gopherfeed/internal/store"
	"github.com/samuel032khoury/gopherfeed/internal/utils"
)

// registerPayload represents the expected payload for authentication endpoints
//
//	@Description	Registration payload
type registerPayload struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=30" example:"newuser"`
	Email    string `json:"email" validate:"required,email" example:"newuser@example.com"`
	Password string `json:"password" validate:"required,min=8,max=72" example:"password"`
}

// loginPayload represents the expected payload for authentication endpoints
//
//	@Description	Login payload
type loginPayload struct {
	Email    string `json:"email" validate:"required,email" example:"user1@example.com"`
	Password string `json:"password" validate:"required,min=8,max=72" example:"password"`
}

// tokenDTO represents the payload for token-based requests
//
//	@Description	Token payload
type tokenDTO struct {
	Token string `json:"token" validate:"required,uuid4" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// activateResponse represents the response payload for account activation
//
//	@Description	Account activation response
type activateResponse struct {
	Message string `json:"message" example:"Account activated successfully"`
}

// loginResponse represents the response payload for user login
//
//	@Description	User login response
type loginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXApJ9..."`
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		registerPayload					true	"User registration payload"
//	@Success		201		{object}	DataResponse[store.User]	"User registered successfully"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := &registerPayload{}
	if err := readJSON(w, r, payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	encryptedPassword, err := utils.EncryptPassword(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: encryptedPassword,
	}
	ctx := r.Context()
	token := uuid.New().String()
	tokenHash := utils.Hash(token)
	exp := utils.InvitationTokenExpiry
	if err := app.store.Users.Register(ctx, user, tokenHash, exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: utils.GenerateActivationURL(app.config.frontendBaseURL, token, isProdEnv),
	}
	err = app.emailPublisher.Publish(user.Email, email.UserInviteTemplate, vars)
	if err != nil {
		app.logger.Errorw("failed to send activation email", "email", user.Email, "error", err)
		// saga
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("failed to rollback user after email send failure", "userID", user.ID, "error", err)
		}
		app.internalServerError(w, r, err)
		return
	}
	app.logger.Infow("email event published", "username", user.Username, "email", user.Email, "token", token)
	app.jsonResponse(w, user, http.StatusCreated)
}

// LoginUser godoc
//
//	@Summary		User login
//	@Description	Authenticate a user and return a token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		loginPayload					true	"User login payload"
//	@Success		200		{object}	DataResponse[loginResponse]	"Authentication token"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/login [post]
func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := &loginPayload{}
	if err := readJSON(w, r, payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	token, err := app.store.Users.Authenticate(ctx, payload.Email, payload.Password, app.authenticator)
	if err != nil {
		switch err {
		case store.ErrInvalidCredentials:
			app.unauthorizedError(w, r, false, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	response := &loginResponse{}
	response.Token = token
	app.jsonResponse(w, response, http.StatusOK)
}

// ActivateUser godoc
//
//	@Summary		Activate a user account
//	@Description	Activate a user account using the provided token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			token	body		tokenDTO						true	"Activation token"
//	@Success		200		{object}	DataResponse[activateResponse]	"Account activated successfully"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/activate [post]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := &tokenDTO{}
	if err := readJSON(w, r, payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.store.Users.Activate(ctx, payload.Token); err != nil {
		switch err {
		case store.ErrInvalidToken:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	response := activateResponse{
		Message: "Account activated successfully",
	}
	app.jsonResponse(w, response, http.StatusOK)
}
