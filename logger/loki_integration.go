package logger

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// LokiConfig holds configuration for the Loki sink
type LokiConfig struct {
	URL         string            // URL to Loki API (e.g., "http://loki:3100/loki/api/v1/push")
	Labels      map[string]string // Labels for log classification
	BatchSize   int               // Max number of entries to batch before sending
	MaxWait     time.Duration     // Max time to wait before sending a batch
	Timeout     time.Duration     // HTTP request timeout
	Compression bool              // Enable gzip compression
	RetryCount  int               // Number of retries for failed requests
	RetryWait   time.Duration     // Time to wait between retries
}

// NewLokiConfig creates a new Loki configuration with defaults
func NewLokiConfig(url string, labels map[string]string) *LokiConfig {
	return &LokiConfig{
		URL:         url,
		Labels:      labels,
		BatchSize:   100,
		MaxWait:     5 * time.Second,
		Timeout:     5 * time.Second,
		Compression: true,
		RetryCount:  3,
		RetryWait:   1 * time.Second,
	}
}

// lokiSink implements zapcore.WriteSyncer for sending logs to Loki
type lokiSink struct {
	config     *LokiConfig
	buffer     []lokiEntry
	bufferLock sync.Mutex
	client     *http.Client
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

type lokiEntry struct {
	Timestamp time.Time
	Line      string
}

type lokiPayload struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// newLokiSink creates a new Loki sink
func newLokiSink(config *LokiConfig) (*lokiSink, error) {
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}
	if config.MaxWait <= 0 {
		config.MaxWait = 5 * time.Second
	}
	if config.Timeout <= 0 {
		config.Timeout = 5 * time.Second
	}
	if config.RetryCount < 0 {
		config.RetryCount = 0
	}
	if config.RetryWait <= 0 {
		config.RetryWait = 1 * time.Second
	}

	client := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 5,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	sink := &lokiSink{
		config: config,
		buffer: make([]lokiEntry, 0, config.BatchSize),
		client: client,
		stopCh: make(chan struct{}),
	}

	sink.wg.Add(1)
	go sink.periodicFlush()

	return sink, nil
}

// Write implements io.Writer
func (s *lokiSink) Write(p []byte) (n int, err error) {
	s.bufferLock.Lock()
	s.buffer = append(
		s.buffer, lokiEntry{
			Timestamp: time.Now(),
			Line:      string(p),
		},
	)
	shouldFlush := len(s.buffer) >= s.config.BatchSize
	s.bufferLock.Unlock()

	if shouldFlush {
		go func() {
			if err := s.flush(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to flush logs to Loki: %v\n", err)
			}
		}()
	}

	return len(p), nil
}

// Sync implements zapcore.WriteSyncer
func (s *lokiSink) Sync() error {
	return s.flush()
}

// Close stops the periodic flushing and flushes any remaining logs
func (s *lokiSink) Close() error {
	close(s.stopCh)
	s.wg.Wait()
	return s.flush()
}

// periodicFlush periodically flushes logs to Loki
func (s *lokiSink) periodicFlush() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.config.MaxWait)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.bufferLock.Lock()
			hasLogs := len(s.buffer) > 0
			s.bufferLock.Unlock()

			if hasLogs {
				if err := s.flush(); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to flush logs to Loki: %v\n", err)
				}
			}
		case <-s.stopCh:
			return
		}
	}
}

// flush sends buffered logs to Loki
func (s *lokiSink) flush() error {
	s.bufferLock.Lock()
	if len(s.buffer) == 0 {
		s.bufferLock.Unlock()
		return nil
	}

	entries := s.buffer
	s.buffer = make([]lokiEntry, 0, s.config.BatchSize)
	s.bufferLock.Unlock()

	return s.sendToLoki(entries)
}

// sendToLoki sends log entries to Loki with retry logic
func (s *lokiSink) sendToLoki(entries []lokiEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()

	payload, err := s.preparePayload(entries)
	if err != nil {
		return fmt.Errorf("failed to prepare payload: %w", err)
	}

	var attempt int
	for attempt = 0; attempt <= s.config.RetryCount; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(s.config.RetryWait):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		err = s.doSend(ctx, payload)
		if err == nil {
			break
		}

		if attempt < s.config.RetryCount {
			fmt.Fprintf(
				os.Stderr, "Failed to send logs to Loki (attempt %d/%d): %v\n",
				attempt+1, s.config.RetryCount+1, err,
			)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to send logs after %d attempts: %w", attempt, err)
	}

	return nil
}

// preparePayload creates the JSON payload for Loki
func (s *lokiSink) preparePayload(entries []lokiEntry) ([]byte, error) {
	values := make([][]string, 0, len(entries))
	for _, entry := range entries {
		timestampNano := fmt.Sprintf("%d", entry.Timestamp.UnixNano())
		values = append(values, []string{timestampNano, entry.Line})
	}

	payload := lokiPayload{
		Streams: []lokiStream{
			{
				Stream: s.config.Labels,
				Values: values,
			},
		},
	}

	return json.Marshal(payload)
}

// doSend performs the actual HTTP request to Loki
func (s *lokiSink) doSend(ctx context.Context, jsonPayload []byte) error {
	var body io.Reader = bytes.NewBuffer(jsonPayload)

	// Apply compression if enabled
	if s.config.Compression {
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write(jsonPayload); err != nil {
			return fmt.Errorf("compression failed: %w", err)
		}
		if err := gz.Close(); err != nil {
			return fmt.Errorf("compression closure failed: %w", err)
		}
		body = &b
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.URL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.config.Compression {
		req.Header.Set("Content-Encoding", "gzip")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Loki API error (status %d): %s", resp.StatusCode, body)
	}

	return nil
}
