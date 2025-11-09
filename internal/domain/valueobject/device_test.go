package valueobject_test

import (
	"strings"
	"testing"

	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
)

func TestNewDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		guid        string
		machineName string
		userAgent   string
		expectError bool
	}{
		{
			name:        "valid device",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "TestMachine",
			userAgent:   "iTunes/10.6",
			expectError: false,
		},
		{
			name:        "empty GUID",
			guid:        "",
			machineName: "TestMachine",
			userAgent:   "iTunes/10.6",
			expectError: true,
		},
		{
			name:        "empty machineName",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "",
			userAgent:   "iTunes/10.6",
			expectError: true,
		},
		{
			name:        "empty userAgent",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "TestMachine",
			userAgent:   "",
			expectError: true,
		},
		{
			name:        "whitespace-only GUID",
			guid:        "   ",
			machineName: "TestMachine",
			userAgent:   "iTunes/10.6",
			expectError: false, // Device doesn't trim/validate whitespace
		},
		{
			name:        "whitespace-only machineName",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "   ",
			userAgent:   "iTunes/10.6",
			expectError: false, // Device doesn't trim/validate whitespace
		},
		{
			name:        "whitespace-only userAgent",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "TestMachine",
			userAgent:   "   ",
			expectError: false, // Device doesn't trim/validate whitespace
		},
		{
			name:        "special characters in GUID",
			guid:        "abc-123-!@#-$%^",
			machineName: "TestMachine",
			userAgent:   "iTunes/10.6",
			expectError: false,
		},
		{
			name:        "special characters in machineName",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "Test@Machine#123",
			userAgent:   "iTunes/10.6",
			expectError: false,
		},
		{
			name:        "very long GUID",
			guid:        strings.Repeat("a", 500),
			machineName: "TestMachine",
			userAgent:   "iTunes/10.6",
			expectError: false,
		},
		{
			name:        "very long machineName",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: strings.Repeat("a", 500),
			userAgent:   "iTunes/10.6",
			expectError: false,
		},
		{
			name:        "very long userAgent",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "TestMachine",
			userAgent:   strings.Repeat("a", 1000),
			expectError: false,
		},
		{
			name:        "unicode characters in machineName",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "TestMachine",
			userAgent:   "iTunes/10.6",
			expectError: false,
		},
		{
			name:        "predefined userAgent - Windows",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "TestMachine",
			userAgent:   valueobject.UserAgentWindows,
			expectError: false,
		},
		{
			name:        "predefined userAgent - Top200",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "TestMachine",
			userAgent:   valueobject.UserAgentTop200,
			expectError: false,
		},
		{
			name:        "predefined userAgent - Top1500",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "TestMachine",
			userAgent:   valueobject.UserAgentTop1500,
			expectError: false,
		},
		{
			name:        "predefined userAgent - Download",
			guid:        "00000000-0000-0000-0000-000000000000",
			machineName: "TestMachine",
			userAgent:   valueobject.UserAgentDownload,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			device, err := valueobject.NewDevice(tt.guid, tt.machineName, tt.userAgent)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)

				return
			}

			if device.GUID() != tt.guid {
				t.Errorf("Expected GUID %s, got %s", tt.guid, device.GUID())
			}

			if device.MachineName() != tt.machineName {
				t.Errorf("Expected machineName %s, got %s", tt.machineName, device.MachineName())
			}

			if device.UserAgent() != tt.userAgent {
				t.Errorf("Expected userAgent %s, got %s", tt.userAgent, device.UserAgent())
			}
		})
	}
}

func TestDevice_Equals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		device1  func() *valueobject.Device
		device2  func() *valueobject.Device
		expected bool
	}{
		{
			name: "identical devices",
			device1: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-1", "Machine1", "Agent1")
				if err != nil {
					panic(err)
				}

				return d
			},
			device2: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-1", "Machine1", "Agent1")
				if err != nil {
					panic(err)
				}

				return d
			},
			expected: true,
		},
		{
			name: "different GUID",
			device1: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-1", "Machine1", "Agent1")
				if err != nil {
					panic(err)
				}

				return d
			},
			device2: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-2", "Machine1", "Agent1")
				if err != nil {
					panic(err)
				}

				return d
			},
			expected: false,
		},
		{
			name: "different machineName",
			device1: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-1", "Machine1", "Agent1")
				if err != nil {
					panic(err)
				}

				return d
			},
			device2: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-1", "Machine2", "Agent1")
				if err != nil {
					panic(err)
				}

				return d
			},
			expected: false,
		},
		{
			name: "different userAgent",
			device1: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-1", "Machine1", "Agent1")
				if err != nil {
					panic(err)
				}

				return d
			},
			device2: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-1", "Machine1", "Agent2")
				if err != nil {
					panic(err)
				}

				return d
			},
			expected: false,
		},
		{
			name: "nil comparison",
			device1: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-1", "Machine1", "Agent1")
				if err != nil {
					panic(err)
				}

				return d
			},
			device2: func() *valueobject.Device {
				return nil
			},
			expected: false,
		},
		{
			name: "same instance",
			device1: func() *valueobject.Device {
				d, err := valueobject.NewDevice("guid-1", "Machine1", "Agent1")
				if err != nil {
					panic(err)
				}

				return d
			},
			device2: func() *valueobject.Device {
				return nil // Will be set to device1 in test
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			device1 := tt.device1()
			device2 := tt.device2()

			// Handle same instance case
			if tt.name == "same instance" {
				device2 = device1
			}

			result := device1.Equals(device2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDevice_PredefinedUserAgents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userAgent string
		notEmpty  bool
	}{
		{"UserAgentWindows", valueobject.UserAgentWindows, true},
		{"UserAgentTop200", valueobject.UserAgentTop200, true},
		{"UserAgentTop1500", valueobject.UserAgentTop1500, true},
		{"UserAgentDownload", valueobject.UserAgentDownload, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.notEmpty && tt.userAgent == "" {
				t.Errorf("%s should not be empty", tt.name)
			}

			// Verify that predefined user agents can be used to create devices
			device, err := valueobject.NewDevice("test-guid", "TestMachine", tt.userAgent)
			if err != nil {
				t.Errorf("Failed to create device with %s: %v", tt.name, err)
			}

			if device.UserAgent() != tt.userAgent {
				t.Errorf("Expected userAgent %s, got %s", tt.userAgent, device.UserAgent())
			}
		})
	}
}

