package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secretKey           string
	Aud                 string
	Iss                 string
	TokenExpiryDuration time.Duration
}

func NewJWTAuthenticator(secretKey, aud, iss, tokenExpiryString string) *JWTAuthenticator {
	tokenExpiryDuration, _ := time.ParseDuration(tokenExpiryString)
	return &JWTAuthenticator{
		secretKey:           secretKey,
		Aud:                 aud,
		Iss:                 iss,
		TokenExpiryDuration: tokenExpiryDuration,
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
	// Implementation for validating JWT token
	return nil, nil
}

func (a *JWTAuthenticator) GetMetadata() (exp time.Duration, iss string, aud string) {
	return a.TokenExpiryDuration, a.Iss, a.Aud
}
