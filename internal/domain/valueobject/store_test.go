package valueobject_test

import (
	"testing"

	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
)

func TestNewStore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		region         string
		storeFront     int
		hostPrefix     int
		expectError    bool
		expectedRegion string
	}{
		{"Valid US store", "us", 143441, 36, false, "us"},
		{"Valid RU store", "ru", 143469, 45, false, "ru"},
		{"Region with uppercase", "US", 143441, 36, false, "us"},
		{"Region with mixed case", "Us", 143441, 36, false, "us"},
		{"Region with whitespace", " us ", 143441, 36, false, "us"},
		{"Empty region", "", 143441, 36, true, ""},
		{"Whitespace-only region", "   ", 143441, 36, true, ""},
		{"Zero storefront", "us", 0, 36, true, ""},
		{"Zero hostprefix", "us", 143441, 0, true, ""},
		{"Negative storefront", "us", -1, 36, true, ""},
		{"Negative hostprefix", "us", 143441, -1, true, ""},
		{"Special characters in region", "u$", 143441, 36, false, "u$"},
		{"Very large storeFront", "us", 999999, 36, false, "us"},
		{"Very large hostPrefix", "us", 143441, 999999, false, "us"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store, err := valueobject.NewStore(tt.region, tt.storeFront, tt.hostPrefix)

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

			if store.Region() != tt.expectedRegion {
				t.Errorf("Expected region %s, got %s", tt.expectedRegion, store.Region())
			}

			if store.StoreFront() != tt.storeFront {
				t.Errorf("Expected storeFront %d, got %d", tt.storeFront, store.StoreFront())
			}

			if store.HostPrefix() != tt.hostPrefix {
				t.Errorf("Expected hostPrefix %d, got %d", tt.hostPrefix, store.HostPrefix())
			}
		})
	}
}

func TestStore_XAppleStoreFront(t *testing.T) {
	t.Parallel()

	store, err := valueobject.NewStore("us", 143441, 36)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	expected := "143441,32"

	if store.XAppleStoreFront() != expected {
		t.Errorf("Expected %s, got %s", expected, store.XAppleStoreFront())
	}
}

func TestStore_XAppleStoreFrontWithDevice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		storeFront int
		deviceCode int
		expected   string
	}{
		{"Standard device code", 143441, 29, "143441,29"},
		{"Zero device code", 143441, 0, "143441,0"},
		{"Negative device code", 143441, -1, "143441,-1"},
		{"Large device code", 143441, 999999, "143441,999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store, err := valueobject.NewStore("us", tt.storeFront, 36)
			if err != nil {
				t.Fatalf("Failed to create store: %v", err)
			}

			result := store.XAppleStoreFrontWithDevice(tt.deviceCode)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestStore_Equals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		store1   func() *valueobject.Store
		store2   func() *valueobject.Store
		expected bool
	}{
		{
			name: "identical stores",
			store1: func() *valueobject.Store {
				s, _ := valueobject.NewStore("us", 143441, 36)
				return s
			},
			store2: func() *valueobject.Store {
				s, _ := valueobject.NewStore("us", 143441, 36)
				return s
			},
			expected: true,
		},
		{
			name: "different regions",
			store1: func() *valueobject.Store {
				s, _ := valueobject.NewStore("us", 143441, 36)
				return s
			},
			store2: func() *valueobject.Store {
				s, _ := valueobject.NewStore("ru", 143441, 36)
				return s
			},
			expected: false,
		},
		{
			name: "different storeFront",
			store1: func() *valueobject.Store {
				s, _ := valueobject.NewStore("us", 143441, 36)
				return s
			},
			store2: func() *valueobject.Store {
				s, _ := valueobject.NewStore("us", 143469, 36)
				return s
			},
			expected: false,
		},
		{
			name: "different hostPrefix",
			store1: func() *valueobject.Store {
				s, _ := valueobject.NewStore("us", 143441, 36)
				return s
			},
			store2: func() *valueobject.Store {
				s, _ := valueobject.NewStore("us", 143441, 45)
				return s
			},
			expected: false,
		},
		{
			name: "nil comparison",
			store1: func() *valueobject.Store {
				s, _ := valueobject.NewStore("us", 143441, 36)
				return s
			},
			store2: func() *valueobject.Store {
				return nil
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			store1 := tt.store1()
			store2 := tt.store2()

			result := store1.Equals(store2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
