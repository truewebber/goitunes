package entity

import (
	"testing"
	"time"
)

func TestNewApplication(t *testing.T) {
	adamID := "123456"
	bundleID := "com.test.app"
	name := "Test App"

	app := NewApplication(adamID, bundleID, name)

	if app.AdamID() != adamID {
		t.Errorf("Expected adamID %s, got %s", adamID, app.AdamID())
	}
	if app.BundleID() != bundleID {
		t.Errorf("Expected bundleID %s, got %s", bundleID, app.BundleID())
	}
	if app.Name() != name {
		t.Errorf("Expected name %s, got %s", name, app.Name())
	}
}

func TestApplication_SettersAndGetters(t *testing.T) {
	app := NewApplication("123", "com.test", "Test")

	// Test price and currency
	app.SetPrice(9.99, "USD")
	if app.Price() != 9.99 {
		t.Errorf("Expected price 9.99, got %f", app.Price())
	}
	if app.Currency() != "USD" {
		t.Errorf("Expected currency USD, got %s", app.Currency())
	}

	// Test rating
	app.SetRating(4.5, 1000)
	if app.Rating() != 4.5 {
		t.Errorf("Expected rating 4.5, got %f", app.Rating())
	}
	if app.RatingCount() != 1000 {
		t.Errorf("Expected rating count 1000, got %d", app.RatingCount())
	}

	// Test version
	app.SetVersion("1.0.0", 123456)
	if app.Version() != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", app.Version())
	}
	if app.VersionID() != 123456 {
		t.Errorf("Expected versionID 123456, got %d", app.VersionID())
	}
}

func TestApplication_IsFree(t *testing.T) {
	tests := []struct {
		name     string
		price    float64
		expected bool
	}{
		{"Free app", 0.0, true},
		{"Paid app", 9.99, false},
		{"Expensive app", 99.99, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApplication("123", "com.test", "Test")
			app.SetPrice(tt.price, "USD")

			if app.IsFree() != tt.expected {
				t.Errorf("Expected IsFree() to be %v for price %f", tt.expected, tt.price)
			}
		})
	}
}

func TestApplication_IsUniversal(t *testing.T) {
	tests := []struct {
		name     string
		families []string
		expected bool
	}{
		{"iPhone only", []string{"iphone"}, false},
		{"iPad only", []string{"ipad"}, false},
		{"Universal", []string{"iphone", "ipad"}, true},
		{"Universal with watch", []string{"iphone", "ipad", "watch"}, true},
		{"Empty", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApplication("123", "com.test", "Test")
			app.SetDeviceFamilies(tt.families)

			if app.IsUniversal() != tt.expected {
				t.Errorf("Expected IsUniversal() to be %v for families %v", tt.expected, tt.families)
			}
		})
	}
}

func TestApplication_BuilderPattern(t *testing.T) {
	releaseDate := time.Now()

	app := NewApplication("123", "com.test", "Test").
		SetArtistName("Test Artist").
		SetArtistID("456").
		SetPrice(4.99, "USD").
		SetRating(4.5, 100).
		SetVersion("1.0", 1).
		SetReleaseDate(releaseDate).
		SetGenre("1", "Games").
		SetFileSize(1024 * 1024 * 50). // 50 MB
		SetMinimumOSVersion("14.0").
		SetDescription("Test description").
		SetIconURL("https://example.com/icon.png")

	if app.ArtistName() != "Test Artist" {
		t.Error("Builder pattern failed for ArtistName")
	}
	if app.Price() != 4.99 {
		t.Error("Builder pattern failed for Price")
	}
	if app.Rating() != 4.5 {
		t.Error("Builder pattern failed for Rating")
	}
	if !app.ReleaseDate().Equal(releaseDate) {
		t.Error("Builder pattern failed for ReleaseDate")
	}
}

