package http

import (
	"fmt"
	"net/http"
)

//go:generate mockgen -source=client.go -destination=mocks/mock_client.go -package=mocks

// Client is an interface for HTTP client operations.
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultClient returns the standard http.Client.
type DefaultClient struct {
	client *http.Client
}

// NewDefaultClient creates a new default HTTP client.
// Redirects are handled manually, so CheckRedirect returns an error to prevent automatic following.
func NewDefaultClient() *DefaultClient {
	return &DefaultClient{
		client: &http.Client{
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				// Return error to prevent automatic redirect following
				// We handle redirects manually in auth_client.go
				return http.ErrUseLastResponse
			},
		},
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
