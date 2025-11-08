package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/truewebber/goitunes/v2/pkg/goitunes"
)

const (
	defaultHTTPTimeout = 30 * time.Second
	maxRetriesDefault  = 3
	topAppsLimit       = 5
	retryDelaySeconds  = 2
)

var (
	// ErrRequestFailedAfterRetries is returned when request fails after all retries.
	ErrRequestFailedAfterRetries = errors.New("request failed after retries")

	// ErrRequestFailedWithStatus is returned when request fails with non-retryable status.
	ErrRequestFailedWithStatus = errors.New("request failed with status")
)

// RetryHTTPClient implements a custom HTTP client with retry logic.
type RetryHTTPClient struct {
	client     *http.Client
	maxRetries int
	retryDelay time.Duration
}

// NewRetryHTTPClient creates a new HTTP client with retry capability.
func NewRetryHTTPClient(maxRetries int, retryDelay time.Duration) *RetryHTTPClient {
	return &RetryHTTPClient{
		client: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// Do executes an HTTP request with retry logic.
func (c *RetryHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response

	var err error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Retry attempt %d/%d...", attempt, c.maxRetries)
			time.Sleep(c.retryDelay)
		}

		resp, err = c.client.Do(req)

		// Success
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// Log the error but continue retrying
		if err != nil {
			log.Printf("Request failed: %v", err)
		} else {
			log.Printf("Server error: %d", resp.StatusCode)

			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Printf("Failed to close response body: %v", closeErr)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %d retries: %w", ErrRequestFailedAfterRetries, c.maxRetries, err)
	}

	return nil, fmt.Errorf("%w: %d after %d retries", ErrRequestFailedWithStatus, resp.StatusCode, c.maxRetries)
}

// LoggingHTTPClient wraps an HTTP client to log all requests.
type LoggingHTTPClient struct {
	client *http.Client
}

// NewLoggingHTTPClient creates a new HTTP client that logs requests.
func NewLoggingHTTPClient() *LoggingHTTPClient {
	return &LoggingHTTPClient{
		client: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
	}
}

// Do executes an HTTP request and logs it.
func (c *LoggingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()

	log.Printf("[HTTP] %s %s", req.Method, req.URL.String())

	resp, err := c.client.Do(req)

	duration := time.Since(start)

	if err != nil {
		log.Printf("[HTTP] Failed after %v: %v", duration, err)

		return nil, fmt.Errorf("http request failed: %w", err)
	}

	log.Printf("[HTTP] %d in %v", resp.StatusCode, duration)

	return resp, nil
}

func main() {
	demonstrateRetryClient()
	demonstrateLoggingClient()
	printBenefits()
}

func demonstrateRetryClient() {
	log.Println("=== Using Custom HTTP Client with Retry Logic ===")

	// Create a custom HTTP client with retry logic
	retryClient := NewRetryHTTPClient(maxRetriesDefault, retryDelaySeconds*time.Second)

	// Create goitunes client with custom HTTP client
	client, err := goitunes.New("us",
		goitunes.WithHTTPClient(retryClient),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Use the client normally - it will retry on failures
	charts, err := client.Charts().GetTop200(
		ctx,
		goitunes.GenreAll,
		goitunes.ChartTypeTopFree,
		goitunes.WithRange(1, topAppsLimit),
	)
	if err != nil {
		log.Fatalf("Failed to get charts: %v", err)
	}

	log.Println("\nTop 5 Free Apps:")

	for i := range charts {
		log.Printf("%d. %s", charts[i].Position, charts[i].App.Name)
	}
}

func demonstrateLoggingClient() {
	log.Println("\n=== Using Custom HTTP Client with Logging ===")

	// Create a logging HTTP client
	loggingClient := NewLoggingHTTPClient()

	// Create another goitunes client with logging
	client2, err := goitunes.New("us",
		goitunes.WithHTTPClient(loggingClient),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// All requests will be logged
	apps, err := client2.Applications().GetByBundleID(ctx, "com.facebook.Facebook")
	if err != nil {
		log.Fatalf("Failed to get app: %v", err)
	}

	if len(apps) > 0 {
		log.Printf("\nApp found: %s", apps[0].Name)
	}
}

func printBenefits() {
	log.Println("=== Custom HTTP Client Benefits ===")
	log.Println("✓ Retry logic for transient failures")
	log.Println("✓ Request/response logging")
	log.Println("✓ Custom timeouts")
	log.Println("✓ Metrics collection")
	log.Println("✓ Request/response modification")
	log.Println("✓ Testing with mock clients")
}
