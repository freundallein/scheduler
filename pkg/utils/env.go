package utils

import (
	"os"
	"strconv"
)

// GetEnv allows to extract environment variables.
func GetEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// GetIntEnv allows extracting environment variables of int type.
// Supports default value for fallback.
func GetIntEnv(key string, fallback int) (int, error) {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fallback, err
		}
		return int(i), nil
	}
	return fallback, nil
}
