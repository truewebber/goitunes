package config

import "errors"

var (
	// ErrUnsupportedRegion is returned when the specified region is not supported.
	ErrUnsupportedRegion = errors.New("unsupported region")
)
