package valueobject

import (
	"fmt"
	"strings"
)

// Store represents an App Store region configuration.
type Store struct {
	region           string
	xAppleStoreFront string
	storeFront       int
	hostPrefix       int
}

// NewStore creates a new Store value object.
func NewStore(region string, storeFront, hostPrefix int) (*Store, error) {
	region = strings.ToLower(strings.TrimSpace(region))

	if region == "" {
		return nil, ErrEmptyRegion
	}

	if storeFront <= 0 {
		return nil, ErrInvalidStoreFront
	}

	if hostPrefix <= 0 {
		return nil, ErrInvalidHostPrefix
	}

	return &Store{
		region:           region,
		storeFront:       storeFront,
		hostPrefix:       hostPrefix,
		xAppleStoreFront: fmt.Sprintf("%d,32", storeFront),
	}, nil
}

// Region returns the region code.
func (s *Store) Region() string { return s.region }

// StoreFront returns the store front ID.
func (s *Store) StoreFront() int { return s.storeFront }

// HostPrefix returns the host prefix.
func (s *Store) HostPrefix() int { return s.hostPrefix }

// XAppleStoreFront returns the X-Apple-Store-Front header value.
func (s *Store) XAppleStoreFront() string { return s.xAppleStoreFront }

// XAppleStoreFrontWithDevice returns the X-Apple-Store-Front header with device code.
func (s *Store) XAppleStoreFrontWithDevice(deviceCode int) string {
	return fmt.Sprintf("%d,%d", s.storeFront, deviceCode)
}

// Equals checks if two stores are equal.
func (s *Store) Equals(other *Store) bool {
	if other == nil {
		return false
	}

	return s.region == other.region &&
		s.storeFront == other.storeFront &&
		s.hostPrefix == other.hostPrefix
}
