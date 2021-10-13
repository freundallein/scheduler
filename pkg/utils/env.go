package utils

import (
	"os"
	"strconv"
)

func GetEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

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
