package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/truewebber/goitunes/pkg/goitunes"
)

// RetryHTTPClient implements a custom HTTP client with retry logic
type RetryHTTPClient struct {
	client     *http.Client
	maxRetries int
	retryDelay time.Duration
}

// NewRetryHTTPClient creates a new HTTP client with retry capability
func NewRetryHTTPClient(maxRetries int, retryDelay time.Duration) *RetryHTTPClient {
	return &RetryHTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// Do executes an HTTP request with retry logic
func (c *RetryHTTPClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			fmt.Printf("Retry attempt %d/%d...\n", attempt, c.maxRetries)
			time.Sleep(c.retryDelay)
		}

		resp, err = c.client.Do(req)

		// Success
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		// Log the error but continue retrying
		if err != nil {
			fmt.Printf("Request failed: %v\n", err)
		} else {
			fmt.Printf("Server error: %d\n", resp.StatusCode)
			resp.Body.Close()
		}
	}

	if err != nil {
		return nil, fmt.Errorf("request failed after %d retries: %w", c.maxRetries, err)
	}

	return nil, fmt.Errorf("request failed with status %d after %d retries", resp.StatusCode, c.maxRetries)
}

// LoggingHTTPClient wraps an HTTP client to log all requests
type LoggingHTTPClient struct {
	client *http.Client
}

// NewLoggingHTTPClient creates a new HTTP client that logs requests
func NewLoggingHTTPClient() *LoggingHTTPClient {
	return &LoggingHTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Do executes an HTTP request and logs it
func (c *LoggingHTTPClient) Do(req *http.Request) (*http.Response, error) {
	start := time.Now()

	fmt.Printf("[HTTP] %s %s\n", req.Method, req.URL.String())

	resp, err := c.client.Do(req)

	duration := time.Since(start)

	if err != nil {
		fmt.Printf("[HTTP] Failed after %v: %v\n", duration, err)
		return nil, err
	}

	fmt.Printf("[HTTP] %d in %v\n", resp.StatusCode, duration)

	return resp, nil
}

func main() {
	fmt.Println("=== Using Custom HTTP Client with Retry Logic ===")

	// Create a custom HTTP client with retry logic
	retryClient := NewRetryHTTPClient(3, 2*time.Second)

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
		goitunes.WithRange(1, 5),
	)
	if err != nil {
		log.Fatalf("Failed to get charts: %v", err)
	}

	fmt.Println("\nTop 5 Free Apps:")
	for _, item := range charts {
		fmt.Printf("%d. %s\n", item.Position, item.App.Name)
	}

	fmt.Println("\n=== Using Custom HTTP Client with Logging ===")

	// Create a logging HTTP client
	loggingClient := NewLoggingHTTPClient()

	// Create another goitunes client with logging
	client2, err := goitunes.New("us",
		goitunes.WithHTTPClient(loggingClient),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// All requests will be logged
	apps, err := client2.Applications().GetByBundleID(ctx, "com.facebook.Facebook")
	if err != nil {
		log.Fatalf("Failed to get app: %v", err)
	}

	if len(apps) > 0 {
		fmt.Printf("\nApp found: %s\n", apps[0].Name)
	}

	fmt.Println("=== Custom HTTP Client Benefits ===")
	fmt.Println("✓ Retry logic for transient failures")
	fmt.Println("✓ Request/response logging")
	fmt.Println("✓ Custom timeouts")
	fmt.Println("✓ Metrics collection")
	fmt.Println("✓ Request/response modification")
	fmt.Println("✓ Testing with mock clients")
}
