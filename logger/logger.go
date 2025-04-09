package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Environment types
const (
	EnvDevelopment = "dev"
	EnvProduction  = "prod"
)

// Log is the global logger instance
var Log *Logger

// Logger wraps zap.SugaredLogger to provide additional functionality
type Logger struct {
	*zap.SugaredLogger
}

// Config holds all configuration options for the logger
type Config struct {
	Debug         bool
	TimeZone      string
	Environment   string
	CallerEnabled bool
	LokiConfig    *LokiConfig
}

// Option defines a function that can modify the logger config
type Option func(*Config)

// WithDebug enables debug logging
func WithDebug(debug bool) Option {
	return func(c *Config) {
		c.Debug = debug
	}
}

// WithTimeZone sets the timezone for log timestamps
func WithTimeZone(tz string) Option {
	return func(c *Config) {
		c.TimeZone = tz
	}
}

// WithEnvironment sets the environment type (dev/prod)
func WithEnvironment(env string) Option {
	return func(c *Config) {
		c.Environment = env
	}
}

// WithCaller enables logging the caller's file and line
func WithCaller(enabled bool) Option {
	return func(c *Config) {
		c.CallerEnabled = enabled
	}
}

// WithLoki configures Loki as a log sink
func WithLoki(lokiConfig *LokiConfig) Option {
	return func(c *Config) {
		c.LokiConfig = lokiConfig
	}
}

// New creates a new logger with the given options
func New(opts ...Option) (*Logger, error) {
	// Default configuration
	cfg := &Config{
		Debug:         false,
		TimeZone:      "",
		Environment:   EnvDevelopment,
		CallerEnabled: true,
	}

	// Apply options
	for _, opt := range opts {
		opt(cfg)
	}

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		CallerKey:      "caller",
		EncodeTime:     getTimeEncoder(cfg.TimeZone),
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	// Set log level
	var level zapcore.Level
	if cfg.Debug {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	// Configure encoder based on environment
	var encoder zapcore.Encoder
	if cfg.Environment == EnvProduction {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure sinks
	sinks := []zapcore.WriteSyncer{zapcore.Lock(os.Stdout)}

	// Add Loki sink if configured
	if cfg.LokiConfig != nil && cfg.LokiConfig.URL != "" {
		lokiSink, err := newLokiSink(cfg.LokiConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create Loki sink: %w", err)
		}
		sinks = append(sinks, lokiSink)
	}

	// Create core and logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(sinks...), level),
	)

	zapOpts := []zap.Option{}
	if cfg.CallerEnabled {
		zapOpts = append(zapOpts, zap.AddCaller())
	}
	zapLogger := zap.New(core, zapOpts...)

	logger := &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}

	// Set global logger
	Log = logger

	zap.ReplaceGlobals(zapLogger)

	return logger, nil
}

// getTimeEncoder returns a time encoder function based on the provided timezone
func getTimeEncoder(timezone string) zapcore.TimeEncoder {
	if timezone == "" {
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02 15:04:05"))
		}
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// Fallback to UTC if timezone is invalid
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format("2006-01-02 15:04:05"))
		}
	}

	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(loc).Format("2006-01-02 15:04:05"))
	}
}
