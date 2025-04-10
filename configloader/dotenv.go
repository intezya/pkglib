package configloader

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

// GetEnv returns the value of the environment variable or the fallback if not found.
// If the environment variable is required (i.e., fallback is empty), it returns an error.
func GetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing environment variable: %s", key)
	}
	return value, nil
}

// GetEnvOrFallback returns the value of the environment variable or the fallback if not found.
func GetEnvOrFallback(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// GetEnvInt returns the integer value of the environment variable or the fallback if not found.
// If the environment variable is required and can't be converted to an integer, it returns an error.
func GetEnvInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		if fallback == 0 {
			return 0, fmt.Errorf("missing environment variable: %s", key)
		}
		return fallback, nil
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to convert %s to int: %v", key, err)
	}
	return result, nil
}

// LoadEnv loads environment variables from a specified file or default .env file.
// It handles error logging but does not panic.
func LoadEnv() {
	envFile := flag.String("env-file", ".env", "Path to .env file")
	flag.Parse()

	err := godotenv.Load(*envFile)
	if err != nil {
		// It's a warning because the application can still function without the env file.
		log.Printf(
			"Warning: Error loading .env file from %s: %v. Is it specified correctly? use --env-file=... flag",
			*envFile,
			err,
		)
	} else {
		log.Printf("Environment variables loaded from %s", *envFile)
	}
}
