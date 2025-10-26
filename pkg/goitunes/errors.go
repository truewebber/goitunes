package goitunes

import (
	"errors"
)

var (
	// ErrUnsupportedRegion is returned when the specified region is not supported
	ErrUnsupportedRegion = errors.New("unsupported region")

	// ErrNotAuthenticated is returned when authentication is required but not provided
	ErrNotAuthenticated = errors.New("not authenticated")

	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrApplicationNotFound is returned when the application is not found
	ErrApplicationNotFound = errors.New("application not found")

	// ErrPurchaseFailed is returned when purchase operation fails
	ErrPurchaseFailed = errors.New("purchase failed")

	// ErrInvalidRequest is returned when the request parameters are invalid
	ErrInvalidRequest = errors.New("invalid request")
)

