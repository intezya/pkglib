package main

import (
	"context"
	"fmt"
	"github.com/intezya/pkglib/logger"
	"os"
	"time"
)

// Simple application startup logging
func exampleBasicLogging() {
	log, err := logger.New(
		logger.WithDebug(true),
		logger.WithEnvironment(logger.EnvDevelopment),
	)
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info(
		"Application started",
		"version", "1.0.0",
		"pid", os.Getpid(),
	)

	// Do something...

	log.Info("Application stopped")
}

// API service with structured logging
func exampleAPIService() {
	// Initialize logger
	log, err := logger.New(
		logger.WithEnvironment(logger.EnvProduction),
		logger.WithLoki(
			logger.NewLokiConfig(
				"http://loki:3100/loki/api/v1/push",
				map[string]string{"service": "api-gateway"},
			),
		),
	)
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Simulate API request handling
	ctx := context.Background()
	requestIDs := []string{"req-123", "req-456", "req-789"}

	for _, reqID := range requestIDs {
		// Create a request-scoped logger
		reqCtx := context.WithValue(ctx, "request_id", reqID)

		log.Info(
			"Request received",
			"method", "GET",
			"path", "/api/users",
			"client_ip", "192.168.1.100",
			"request_id", reqID,
		)

		// Simulate processing
		processAPIRequest(reqCtx, log)

		log.Info(
			"Request completed",
			"status", 200,
			"response_time_ms", 25,
			"request_id", reqID,
		)
	}
}

// Process an API request with context and logger
func processAPIRequest(ctx context.Context, log *logger.Logger) {
	// Get request ID from context
	reqID := ctx.Value("request_id").(string)

	// Simulate database query
	log.Debug(
		"Executing database query",
		"query", "SELECT * FROM users LIMIT 10",
		"query_id", fmt.Sprintf("%s-query1", reqID),
	)

	// Simulate request processing error
	if reqID == "req-456" {
		log.Error(
			"Error processing request",
			"error", "database connection failed",
			"component", "database",
			"retry_count", 3,
		)
	}
}

// Background job with progress logging
func exampleBackgroundJob() {
	log, err := logger.New(
		logger.WithEnvironment(logger.EnvProduction),
	)
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Start a background job
	jobID := "job-abc-123"
	totalItems := 1000

	log.Info(
		"Background job started",
		"job_id", jobID,
		"total_items", totalItems,
	)

	// Process in batches
	batchSize := 100
	batches := totalItems / batchSize

	for i := 0; i < batches; i++ {
		start := time.Now()

		// Process batch
		processedItems := batchSize
		if i == 2 {
			// Simulate partial batch failure
			processedItems = 80
			log.Warn(
				"Partial batch processing failure",
				"job_id", jobID,
				"batch", i+1,
				"failed_items", batchSize-processedItems,
				"error", "timeout",
			)
		}

		duration := time.Since(start)

		// Log progress
		log.Info(
			"Batch processed",
			"job_id", jobID,
			"batch", i+1,
			"total_batches", batches,
			"items_processed", processedItems,
			"duration_ms", duration.Milliseconds(),
			"progress_pct", float64((i+1)*batchSize)/float64(totalItems)*100,
		)
	}

	log.Info(
		"Background job completed",
		"job_id", jobID,
		"total_processed", totalItems,
		"status", "success",
	)
}

// Application with environment-based configuration
func exampleEnvironmentBasedConfig() {
	// Get environment settings
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = logger.EnvDevelopment
	}

	debug := os.Getenv("DEBUG") == "true"

	// Configure Loki if we're in a non-dev environment
	var lokiConfig *logger.LokiConfig
	if env != logger.EnvDevelopment {
		lokiURL := os.Getenv("LOKI_URL")
		if lokiURL != "" {
			lokiConfig = logger.NewLokiConfig(
				lokiURL,
				map[string]string{
					"app": "inventory-service",
					"env": env,
				},
			)
		}
	}

	// Create logger with environment-specific configuration
	log, err := logger.New(
		logger.WithDebug(debug),
		logger.WithEnvironment(env),
		logger.WithTimeZone(os.Getenv("TZ")),
	)
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}

	if lokiConfig != nil {
		// Add Loki sink if configured
		log, err = logger.New(
			logger.WithDebug(debug),
			logger.WithEnvironment(env),
			logger.WithTimeZone(os.Getenv("TZ")),
			logger.WithLoki(lokiConfig),
		)
		if err != nil {
			fmt.Println("Failed to initialize logger with Loki:", err)
			os.Exit(1)
		}
	}

	defer log.Sync()

	log.Info(
		"Application initialized",
		"env", env,
		"debug", debug,
		"loki_enabled", lokiConfig != nil,
	)
}

// Performance monitoring with logger
func examplePerformanceMonitoring() {
	log, err := logger.New(
		logger.WithEnvironment(logger.EnvProduction),
	)
	if err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Track function execution time
	operationLog := func(name string, fn func()) {
		start := time.Now()
		log.Debug("Operation started", "operation", name)

		fn() // Execute the function

		duration := time.Since(start)
		log.Info(
			"Operation completed",
			"operation", name,
			"duration_ms", duration.Milliseconds(),
		)
	}

	// Use the timing function
	operationLog(
		"data_import", func() {
			// Simulate work
			time.Sleep(150 * time.Millisecond)
		},
	)

	operationLog(
		"data_processing", func() {
			// Simulate work
			time.Sleep(75 * time.Millisecond)
		},
	)

	operationLog(
		"data_export", func() {
			// Simulate work
			time.Sleep(100 * time.Millisecond)
		},
	)
}

func main() {
	fmt.Println("Running logger examples...")

	// Run examples
	exampleBasicLogging()
	exampleAPIService()
	exampleBackgroundJob()
	exampleEnvironmentBasedConfig()
	examplePerformanceMonitoring()

	fmt.Println("All examples completed")
}
