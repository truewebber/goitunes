package valueobject

import (
	"fmt"
	"strings"
)

// Store represents an App Store region configuration
type Store struct {
	region           string
	storeFront       int
	hostPrefix       int
	xAppleStoreFront string
}

// NewStore creates a new Store value object
func NewStore(region string, storeFront, hostPrefix int) (*Store, error) {
	region = strings.ToLower(strings.TrimSpace(region))
	if region == "" {
		return nil, fmt.Errorf("region cannot be empty")
	}
	if storeFront <= 0 {
		return nil, fmt.Errorf("storeFront must be positive")
	}
	if hostPrefix <= 0 {
		return nil, fmt.Errorf("hostPrefix must be positive")
	}

	return &Store{
		region:           region,
		storeFront:       storeFront,
		hostPrefix:       hostPrefix,
		xAppleStoreFront: fmt.Sprintf("%d,32", storeFront),
	}, nil
}

// Getters
func (s *Store) Region() string           { return s.region }
func (s *Store) StoreFront() int          { return s.storeFront }
func (s *Store) HostPrefix() int          { return s.hostPrefix }
func (s *Store) XAppleStoreFront() string { return s.xAppleStoreFront }

// XAppleStoreFrontWithDevice returns the X-Apple-Store-Front header with device code
func (s *Store) XAppleStoreFrontWithDevice(deviceCode int) string {
	return fmt.Sprintf("%d,%d", s.storeFront, deviceCode)
}

// Equals checks if two stores are equal
func (s *Store) Equals(other *Store) bool {
	if other == nil {
		return false
	}
	return s.region == other.region &&
		s.storeFront == other.storeFront &&
		s.hostPrefix == other.hostPrefix
}

