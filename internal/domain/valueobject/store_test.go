package valueobject

import "testing"

func TestNewStore(t *testing.T) {
	tests := []struct {
		name        string
		region      string
		storeFront  int
		hostPrefix  int
		expectError bool
	}{
		{"Valid US store", "us", 143441, 36, false},
		{"Valid RU store", "ru", 143469, 45, false},
		{"Empty region", "", 143441, 36, true},
		{"Zero storefront", "us", 0, 36, true},
		{"Zero hostprefix", "us", 143441, 0, true},
		{"Negative storefront", "us", -1, 36, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, err := NewStore(tt.region, tt.storeFront, tt.hostPrefix)

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

			if store.Region() != tt.region {
				t.Errorf("Expected region %s, got %s", tt.region, store.Region())
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
	store, _ := NewStore("us", 143441, 36)
	expected := "143441,32"

	if store.XAppleStoreFront() != expected {
		t.Errorf("Expected %s, got %s", expected, store.XAppleStoreFront())
	}
}

func TestStore_XAppleStoreFrontWithDevice(t *testing.T) {
	store, _ := NewStore("us", 143441, 36)
	deviceCode := 29
	expected := "143441,29"

	result := store.XAppleStoreFrontWithDevice(deviceCode)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestStore_Equals(t *testing.T) {
	store1, _ := NewStore("us", 143441, 36)
	store2, _ := NewStore("us", 143441, 36)
	store3, _ := NewStore("ru", 143469, 45)

	if !store1.Equals(store2) {
		t.Error("Same stores should be equal")
	}

	if store1.Equals(store3) {
		t.Error("Different stores should not be equal")
	}

	if store1.Equals(nil) {
		t.Error("Store should not equal nil")
	}
}
