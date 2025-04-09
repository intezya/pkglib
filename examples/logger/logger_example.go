package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/intezya/pkglib/logger"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	//Basic logger setup
	log, err := logger.New(
		logger.WithDebug(true),
		logger.WithEnvironment(logger.EnvDevelopment),
	)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Don't forget to flush logs at the end
	defer log.Sync()

	// Basic logging at different levels
	log.Debug("This is a debug message")
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")

	//Structured logging with key-value pairs
	user := struct {
		ID        string
		Name      string
		Email     string
		CreatedAt time.Time
	}{
		ID:        "user-123",
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}

	log.Info(
		"User logged in",
		"user_id", user.ID,
		"name", user.Name,
		"email", user.Email,
		"account_age_hours", time.Since(user.CreatedAt).Hours(),
	)

	//Logging errors with context
	if err := performDatabaseOperation(); err != nil {
		log.Error(
			"Database operation failed",
			"operation", "update_user",
			"error", err,
			"user_id", user.ID,
			"retry", false,
		)
	}

	//Setup with Loki for production
	prodLog, err := logger.New(
		logger.WithEnvironment(logger.EnvProduction),
		logger.WithTimeZone("UTC"),
		logger.WithLoki(
			logger.NewLokiConfig(
				"http://loki:3100/loki/api/v1/push",
				map[string]string{
					"app":     "user-service",
					"env":     "production",
					"version": "1.2.3",
				},
			),
		),
	)
	if err != nil {
		fmt.Printf("Failed to initialize production logger: %v\n", err)
		os.Exit(1)
	}
	defer prodLog.Sync()

	prodLog.Info("Production logger initialized", "timestamp", time.Now().Unix())

	//Using context with request ID
	ctx := context.WithValue(context.Background(), "request_id", "req-abc-123")
	requestLog := log.With("request_id", ctx.Value("request_id"))

	requestLog.Info(
		"Processing API request",
		"method", "GET",
		"path", "/api/users",
		"params", map[string]string{"limit": "10", "offset": "0"},
	)

	//Logging HTTP request/response details
	logHTTPRequest()

	//Handling different environments based on config
	envType := os.Getenv("ENV")
	if envType == "" {
		envType = logger.EnvDevelopment
	}

	appLog, err := logger.New(
		logger.WithDebug(os.Getenv("DEBUG") == "true"),
		logger.WithEnvironment(envType),
		logger.WithTimeZone(getTimezone()),
	)
	if err != nil {
		fmt.Printf("Failed to initialize app logger: %v\n", err)
		os.Exit(1)
	}

	appLog.Info(
		"Application started",
		"env", envType,
		"debug", os.Getenv("DEBUG") == "true",
		"timezone", getTimezone(),
	)

	//Custom Loki configuration
	lokiConfig := logger.NewLokiConfig(
		"http://loki:3100/loki/api/v1/push",
		map[string]string{"app": "payment-service"},
	)

	// Configure Loki behavior
	lokiConfig.BatchSize = 200            // Send logs in batches of 200
	lokiConfig.MaxWait = 10 * time.Second // Flush at least every 10 seconds
	lokiConfig.Compression = true         // Enable gzip compression
	lokiConfig.RetryCount = 3             // Retry failed requests 3 times
	lokiConfig.RetryWait = time.Second    // Wait 1 second between retries

	// Create logger with custom Loki config
	paymentLog, err := logger.New(
		logger.WithEnvironment(logger.EnvProduction),
		logger.WithLoki(lokiConfig),
	)
	if err != nil {
		fmt.Printf("Failed to initialize payment logger: %v\n", err)
		os.Exit(1)
	}

	paymentLog.Info(
		"Payment processed",
		"payment_id", "pmt-456",
		"amount", 99.95,
		"currency", "USD",
		"status", "completed",
	)

	//Using logger in HTTP middleware
	http.HandleFunc("/api/users", loggingMiddleware(handleUsers, log))

	//Logging benchmark results
	start := time.Now()
	// Simulate work
	time.Sleep(50 * time.Millisecond)
	duration := time.Since(start)

	log.Info(
		"Operation completed",
		"operation", "data_processing",
		"records_processed", 1000,
		"duration_ms", duration.Milliseconds(),
		"throughput", float64(1000)/duration.Seconds(),
	)
}

// Simulates a database operation that might fail
func performDatabaseOperation() error {
	// Simulate database error
	return errors.New("connection timeout: could not connect to database")
}

// Log HTTP request details
func logHTTPRequest() {
	// Create a new logger
	log, _ := logger.New(logger.WithEnvironment(logger.EnvDevelopment))

	// Mock HTTP request
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/api/orders",
			RawQuery: "user_id=123",
		},
		Header: http.Header{
			"Content-Type":     []string{"application/json"},
			"X-Correlation-Id": []string{"corr-xyz-789"},
			"User-Agent":       []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
		},
		RemoteAddr: "192.168.1.1:54321",
	}

	// Log request details
	log.Info(
		"Received HTTP request",
		"method", req.Method,
		"path", req.URL.Path,
		"query", req.URL.RawQuery,
		"remote_addr", req.RemoteAddr,
		"content_type", req.Header.Get("Content-Type"),
		"correlation_id", req.Header.Get("X-Correlation-Id"),
		"user_agent", req.Header.Get("User-Agent"),
	)

	// Log response details
	log.Info(
		"Sending HTTP response",
		"status", 201,
		"content_type", "application/json",
		"body_size_bytes", 256,
		"duration_ms", 15,
	)
}

// HTTP handler with logging middleware
func loggingMiddleware(next http.HandlerFunc, log *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Add request ID to context
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		ctx := context.WithValue(r.Context(), "request_id", requestID)
		r = r.WithContext(ctx)

		// Create a logger with request context
		reqLog := log.With(
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		reqLog.Info("Request started")

		// Create response wrapper to capture status code
		lrw := newLoggingResponseWriter(w)

		// Call the next handler
		next(lrw, r)

		// Log completed request
		duration := time.Since(start)
		reqLog.Info(
			"Request completed",
			"status", lrw.statusCode,
			"duration_ms", duration.Milliseconds(),
			"bytes_written", lrw.bytesWritten,
		)
	}
}

// Handle users endpoint
func handleUsers(w http.ResponseWriter, r *http.Request) {
	// Handler implementation
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"users": [{"id": 1, "name": "John"}, {"id": 2, "name": "Jane"}]}`))
}

// Helper to capture response data for logging
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.bytesWritten += size
	return size, err
}

// Get timezone from environment or use default
func getTimezone() string {
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "UTC"
	}
	return tz
}

// Generate a request ID
func generateRequestID() string {
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}
