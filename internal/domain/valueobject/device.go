package valueobject

import (
	"fmt"
)

// Device represents a device configuration for App Store requests
type Device struct {
	guid        string
	machineName string
	userAgent   string
}

// NewDevice creates a new Device value object
func NewDevice(guid, machineName, userAgent string) (*Device, error) {
	if guid == "" {
		return nil, fmt.Errorf("guid cannot be empty")
	}
	if machineName == "" {
		return nil, fmt.Errorf("machineName cannot be empty")
	}
	if userAgent == "" {
		return nil, fmt.Errorf("userAgent cannot be empty")
	}

	return &Device{
		guid:        guid,
		machineName: machineName,
		userAgent:   userAgent,
	}, nil
}

// Getters
func (d *Device) GUID() string        { return d.guid }
func (d *Device) MachineName() string { return d.machineName }
func (d *Device) UserAgent() string   { return d.userAgent }

// Equals checks if two devices are equal
func (d *Device) Equals(other *Device) bool {
	if other == nil {
		return false
	}
	return d.guid == other.guid &&
		d.machineName == other.machineName &&
		d.userAgent == other.userAgent
}

// DefaultUserAgents contains commonly used user agents
var (
	UserAgentTop200   = "AppStore/2.0 iOS/9.0 model/iPhone6,1 hwp/s5l8960x build/13A344 (6; dt:89)"
	UserAgentTop1500  = "iTunes-iPad/5.1.1 (64GB; dt:28)"
	UserAgentDownload = "itunesstored/1.0 iOS/9.0 model/iPhone6,1 hwp/s5l8960x build/13A344 (6; dt:89)"
	UserAgentWindows  = "iTunes/10.6 (Windows; Microsoft Windows 7 x64 Ultimate Edition Service Pack 1 (Build 7601)) AppleWebKit/534.54.16"
)

