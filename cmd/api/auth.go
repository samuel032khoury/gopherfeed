package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/samuel032khoury/gopherfeed/internal/store"
	"github.com/samuel032khoury/gopherfeed/internal/utils"
)

// authPayload represents the expected payload for authentication endpoints
//
//	@Description	Authentication payload
type authPayload struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=30" example:"johndoe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=8,max=72" example:"password123"`
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		authPayload					true	"User registration payload"
//	@Success		201		{object}	DataResponse[store.User]	"User registered successfully"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := &authPayload{}
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
	app.jsonResponse(w, user, http.StatusCreated)
}

// LoginUser godoc
//
//	@Summary		User login
//	@Description	Authenticate a user and return a token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		authPayload				true	"User login payload"
//	@Success		200		{object}	DataResponse[string]	"Authentication token"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/auth/login [post]
func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := &authPayload{}
	if err := readJSON(w, r, payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
}
