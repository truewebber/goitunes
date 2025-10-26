package goitunes

import (
	"github.com/truewebber/goitunes/internal/domain/valueobject"
	infrahttp "github.com/truewebber/goitunes/internal/infrastructure/http"
)

// Option is a functional option for configuring the Client
type Option func(*Client) error

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient infrahttp.Client) Option {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

// WithCredentials sets authentication credentials
func WithCredentials(appleID, passwordToken, dsid string) Option {
	return func(c *Client) error {
		credentials, err := valueobject.NewCredentialsWithTokens(appleID, passwordToken, dsid)
		if err != nil {
			return err
		}
		c.credentials = credentials
		return nil
	}
}

// WithKbsync sets the kbsync certificate for purchases
func WithKbsync(kbsync string) Option {
	return func(c *Client) error {
		if c.credentials == nil {
			return ErrInvalidCredentials
		}
		c.credentials.SetKbsync(kbsync)
		return nil
	}
}

// WithDevice sets custom device information
func WithDevice(guid, machineName, userAgent string) Option {
	return func(c *Client) error {
		device, err := valueobject.NewDevice(guid, machineName, userAgent)
		if err != nil {
			return err
		}
		c.device = device
		return nil
	}
}

