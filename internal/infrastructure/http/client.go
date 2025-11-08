package http

import (
	"fmt"
	"net/http"
)

// Client is an interface for HTTP client operations.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultClient returns the standard http.Client.
type DefaultClient struct {
	client *http.Client
}

// NewDefaultClient creates a new default HTTP client.
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		client: &http.Client{},
	}
}

// Do executes an HTTP request.
func (c *DefaultClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	return resp, nil
}
