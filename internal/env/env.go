package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}

func GetBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return boolValue
}
