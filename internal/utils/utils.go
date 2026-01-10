package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const InvitationTokenExpiry = 24 * time.Hour // 24 hours in seconds

func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Hash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func GenerateActivationURL(frontendBaseURL, token string, isProdEnv bool) string {
	scheme := "http"
	if isProdEnv {
		scheme = "https"
	}
	return scheme + "://" + frontendBaseURL + "/activate?token=" + token
}
