package goitunes

import (
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
	infrahttp "github.com/truewebber/goitunes/v2/internal/infrastructure/http"
)

// Option is a functional option for configuring the Client.
type Option func(*Client) error

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient infrahttp.Client) Option {
	return func(c *Client) error {
		c.httpClient = httpClient

		return nil
	}
}

// WithAppleID sets the Apple ID for authentication.
// Use this when you plan to authenticate later via Login().
// Tokens will be set automatically after successful login.
func WithAppleID(appleID string) Option {
	return func(c *Client) error {
		credentials, err := valueobject.NewCredentials(appleID)
		if err != nil {
			return fmt.Errorf("failed to create credentials: %w", err)
		}

		c.credentials = credentials

		return nil
	}
}

// WithCredentials sets authentication credentials with tokens.
// Use this when you already have passwordToken and dsid.
func WithCredentials(appleID, passwordToken, dsid string) Option {
	return func(c *Client) error {
		credentials, err := valueobject.NewCredentialsWithTokens(appleID, passwordToken, dsid)
		if err != nil {
			return fmt.Errorf("failed to create credentials: %w", err)
		}

		c.credentials = credentials

		return nil
	}
}

// WithKbsync sets the kbsync certificate for purchases.
func WithKbsync(kbsync string) Option {
	return func(c *Client) error {
		if c.credentials == nil {
			return ErrInvalidCredentials
		}

		c.credentials.SetKbsync(kbsync)

		return nil
	}
}

// WithDevice sets custom device information.
func WithDevice(guid, machineName, userAgent string) Option {
	return func(c *Client) error {
		device, err := valueobject.NewDevice(guid, machineName, userAgent)
		if err != nil {
			return fmt.Errorf("failed to create device: %w", err)
		}

		c.device = device

		return nil
	}
}
