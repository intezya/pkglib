package configloader

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

// GetEnv returns the value of the environment variable.
// If the environment variable is not found, it returns an error.
func GetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("missing environment variable: %s", key)
	}
	return value, nil
}

// GetEnvOrPanic returns the value of the environment variable.
// If the environment variable is not found, it panics.
func GetEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("missing environment variable: %s", key))
	}
	return value
}

// GetEnvOrFallback returns the value of the environment variable.
// If the environment variable is not found, it returns the fallback value.
func GetEnvOrFallback(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// GetEnvInt returns the integer value of the environment variable.
// If the environment variable is not found and fallback is 0, it returns an error.
// If the environment variable is not found and fallback is not 0, it returns the fallback.
// If the environment variable cannot be converted to an integer, it returns an error.
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

// GetEnvIntOrPanic returns the integer value of the environment variable.
// If the environment variable is not found or cannot be converted to an integer, it panics.
func GetEnvIntOrPanic(key string) int {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("missing environment variable: %s", key))
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("failed to convert %s to int: %v", key, err))
	}
	return result
}

// GetEnvIntOrFallback returns the integer value of the environment variable.
// If the environment variable is not found or cannot be converted to an integer,
// it returns the fallback value.
func GetEnvIntOrFallback(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return result
}

type Logger interface {
	Println(v ...any)
}

// LoadEnv loads environment variables from a specified file or default .env file.
// It accepts a command-line flag --env-file to specify the path to the .env file.
// If the flag is not provided, it defaults to ".env" in the current directory.
// It handles error logging but does not panic if the file cannot be loaded.
func LoadEnv(logger ...Logger) {
	logFunc := log.Println

	if len(logger) > 0 {
		logFunc = logger[0].Println
	}

	envFile := flag.String("env-file", ".env", "Path to .env file")
	flag.Parse()

	err := godotenv.Load(*envFile)
	if err != nil {
		// It's a warning because the application can still function without the env file.
		// Environment variables might be set by other means (system environment, docker, etc.)
		logFunc(fmt.Sprintf(
			"Warning: Error loading .env file from %s: %v. Is it specified correctly? use --env-file=... flag",
			*envFile,
			err,
		))
	} else {
		logFunc(fmt.Sprintf("Environment variables loaded from %s", *envFile))
	}
}
