package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read the authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedError(w, r, true, fmt.Errorf("missing authorization header"))
				return
			}
			// parse -> get the base64 encoded username:password
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedError(w, r, true, fmt.Errorf("malformed authorization header"))
				return
			}
			// decode
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedError(w, r, true, fmt.Errorf("invalid base64 encoding"))
				return
			}
			// split into username and password
			username := app.config.auth.basic.username
			password := app.config.auth.basic.password
			creds := strings.SplitN(string(decoded), ":", 2)
			fmt.Println(creds, username, password)
			if len(creds) != 2 || creds[0] != username || creds[1] != password {
				app.unauthorizedError(w, r, true, fmt.Errorf("invalid authorization value"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
