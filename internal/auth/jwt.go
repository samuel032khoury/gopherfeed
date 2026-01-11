package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secretKey           string
	aud                 string
	iss                 string
	tokenExpiryDuration time.Duration
}

func NewJWTAuthenticator(secretKey, aud, iss, tokenExpiryString string) *JWTAuthenticator {
	tokenExpiryDuration, _ := time.ParseDuration(tokenExpiryString)
	return &JWTAuthenticator{
		secretKey:           secretKey,
		aud:                 aud,
		iss:                 iss,
		tokenExpiryDuration: tokenExpiryDuration,
	}
}

func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	// Implementation for generating JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(a.secretKey), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}

func (a *JWTAuthenticator) GetMetadata() (exp time.Duration, iss string, aud string) {
	return a.tokenExpiryDuration, a.iss, a.aud
}
