package entity_test

import (
	"testing"
	"time"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
)

const (
	testRatingValue = 4.5
)

func TestNewApplication(t *testing.T) {
	t.Parallel()

	adamID := "123456"
	bundleID := "com.test.app"
	name := "Test App"

	app := entity.NewApplication(adamID, bundleID, name)

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
	t.Parallel()

	app := entity.NewApplication("123", "com.test", "Test")

	// Test price and currency
	app.SetPrice(9.99, "USD")

	if app.Price() != 9.99 {
		t.Errorf("Expected price 9.99, got %f", app.Price())
	}

	if app.Currency() != "USD" {
		t.Errorf("Expected currency USD, got %s", app.Currency())
	}

	// Test rating
	app.SetRating(testRatingValue, 1000)

	if app.Rating() != testRatingValue {
		t.Errorf("Expected rating %f, got %f", testRatingValue, app.Rating())
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
	t.Parallel()

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
			t.Parallel()

			app := entity.NewApplication("123", "com.test", "Test")
			app.SetPrice(tt.price, "USD")

			if app.IsFree() != tt.expected {
				t.Errorf("Expected IsFree() to be %v for price %f", tt.expected, tt.price)
			}
		})
	}
}

func TestApplication_IsUniversal(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			app := entity.NewApplication("123", "com.test", "Test")
			app.SetDeviceFamilies(tt.families)

			if app.IsUniversal() != tt.expected {
				t.Errorf("Expected IsUniversal() to be %v for families %v", tt.expected, tt.families)
			}
		})
	}
}

func TestApplication_BuilderPattern(t *testing.T) {
	t.Parallel()

	releaseDate := time.Now()

	app := entity.NewApplication("123", "com.test", "Test").
		SetArtistName("Test Artist").
		SetArtistID("456").
		SetPrice(4.99, "USD").
		SetRating(testRatingValue, 100).
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

	if app.Rating() != testRatingValue {
		t.Error("Builder pattern failed for Rating")
	}

	if !app.ReleaseDate().Equal(releaseDate) {
		t.Error("Builder pattern failed for ReleaseDate")
	}
}

func TestNewApplication_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		adamID   string
		bundleID string
		appName  string
	}{
		{"empty strings", "", "", ""},
		{"only adamID", "123", "", ""},
		{"only bundleID", "", "com.test", ""},
		{"only name", "", "", "Test"},
		{"special characters", "!@#", "com.special.app", "Test App 🎮"},
		{"very long values", string(make([]byte, 1000)), string(make([]byte, 1000)), string(make([]byte, 1000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := entity.NewApplication(tt.adamID, tt.bundleID, tt.appName)

			if app == nil {
				t.Fatal("NewApplication should not return nil")
			}

			if app.AdamID() != tt.adamID {
				t.Errorf("Expected adamID %s, got %s", tt.adamID, app.AdamID())
			}

			if app.BundleID() != tt.bundleID {
				t.Errorf("Expected bundleID %s, got %s", tt.bundleID, app.BundleID())
			}

			if app.Name() != tt.appName {
				t.Errorf("Expected name %s, got %s", tt.appName, app.Name())
			}

			// Check default values
			if len(app.DeviceFamilies()) != 0 {
				t.Error("DeviceFamilies should be empty by default")
			}

			if len(app.ScreenshotURLs()) != 0 {
				t.Error("ScreenshotURLs should be empty by default")
			}
		})
	}
}

func TestApplication_SetPrice_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		price    float64
		currency string
	}{
		{"zero price", 0.0, "USD"},
		{"negative price", -1.0, "USD"},
		{"very large price", 999999.99, "USD"},
		{"empty currency", 9.99, ""},
		{"unicode currency", 9.99, "€"},
		{"multiple currencies", 9.99, "$€¥"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := entity.NewApplication("123", "com.test", "Test")
			result := app.SetPrice(tt.price, tt.currency)

			if result != app {
				t.Error("SetPrice should return the same instance")
			}

			if app.Price() != tt.price {
				t.Errorf("Expected price %f, got %f", tt.price, app.Price())
			}

			if app.Currency() != tt.currency {
				t.Errorf("Expected currency %s, got %s", tt.currency, app.Currency())
			}
		})
	}
}

func TestApplication_SetRating_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		rating float64
		count  int
	}{
		{"zero rating zero count", 0.0, 0},
		{"negative rating", -1.0, 100},
		{"rating above 5", 10.0, 100},
		{"negative count", 4.5, -1},
		{"very large count", 4.5, 999999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := entity.NewApplication("123", "com.test", "Test")
			result := app.SetRating(tt.rating, tt.count)

			if result != app {
				t.Error("SetRating should return the same instance")
			}

			if app.Rating() != tt.rating {
				t.Errorf("Expected rating %f, got %f", tt.rating, app.Rating())
			}

			if app.RatingCount() != tt.count {
				t.Errorf("Expected count %d, got %d", tt.count, app.RatingCount())
			}
		})
	}
}

func TestApplication_SetScreenshotURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		urls []string
	}{
		{"valid URLs", []string{"https://example.com/1.png", "https://example.com/2.png"}},
		{"empty slice", []string{}},
		{"nil slice", nil},
		{"single URL", []string{"https://example.com/1.png"}},
		{"many URLs", []string{"url1", "url2", "url3", "url4", "url5"}},
		{"duplicate URLs", []string{"url1", "url1", "url1"}},
		{"empty strings", []string{"", "", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := entity.NewApplication("123", "com.test", "Test")
			result := app.SetScreenshotURLs(tt.urls)

			if result != app {
				t.Error("SetScreenshotURLs should return the same instance")
			}

			urls := app.ScreenshotURLs()
			if tt.urls == nil && urls == nil {
				return // Both nil acceptable
			}

			if tt.urls != nil && urls == nil {
				t.Error("URLs should not be nil when set with non-nil slice")

				return
			}

			if len(urls) != len(tt.urls) {
				t.Errorf("Expected %d URLs, got %d", len(tt.urls), len(urls))
			}
		})
	}
}

func TestApplication_SetDeviceFamilies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		families []string
	}{
		{"valid families", []string{"iphone", "ipad"}},
		{"empty slice", []string{}},
		{"nil slice", nil},
		{"single family", []string{"iphone"}},
		{"with watch", []string{"iphone", "ipad", "watch"}},
		{"duplicates", []string{"iphone", "iphone"}},
		{"uppercase", []string{"iPhone", "iPad"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := entity.NewApplication("123", "com.test", "Test")
			result := app.SetDeviceFamilies(tt.families)

			if result != app {
				t.Error("SetDeviceFamilies should return the same instance")
			}

			families := app.DeviceFamilies()
			if tt.families == nil && families == nil {
				return
			}

			if tt.families != nil && families == nil {
				t.Error("Families should not be nil when set with non-nil slice")

				return
			}

			if len(families) != len(tt.families) {
				t.Errorf("Expected %d families, got %d", len(tt.families), len(families))
			}
		})
	}
}

func TestApplication_IsFree_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		price    float64
		expected bool
	}{
		{"exactly zero", 0.0, true},
		{"negative price", -0.01, false},
		{"very small positive", 0.01, false},
		{"not set (default zero)", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := entity.NewApplication("123", "com.test", "Test")
			if tt.price != 0.0 {
				app.SetPrice(tt.price, "USD")
			}

			if app.IsFree() != tt.expected {
				t.Errorf("Expected IsFree() = %v for price %f", tt.expected, tt.price)
			}
		})
	}
}

func TestApplication_IsUniversal_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		families []string
		expected bool
	}{
		{"case sensitive - lowercase", []string{"iphone", "ipad"}, true},
		{"case sensitive - uppercase", []string{"iPhone", "iPad"}, false},
		{"mixed case", []string{"iPhone", "ipad"}, false},
		{"iphone first", []string{"iphone", "ipad"}, true},
		{"ipad first", []string{"ipad", "iphone"}, true},
		{"with extras before", []string{"watch", "iphone", "ipad"}, true},
		{"with extras after", []string{"iphone", "ipad", "watch"}, true},
		{"only iPhone", []string{"iphone"}, false},
		{"only iPad", []string{"ipad"}, false},
		{"duplicate iphone", []string{"iphone", "iphone", "ipad"}, true},
		{"empty array", []string{}, false},
		{"nil array", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := entity.NewApplication("123", "com.test", "Test")
			app.SetDeviceFamilies(tt.families)

			if app.IsUniversal() != tt.expected {
				t.Errorf("Expected IsUniversal() = %v for families %v", tt.expected, tt.families)
			}
		})
	}
}

func TestApplication_AllSetters(t *testing.T) {
	t.Parallel()

	app := entity.NewApplication("123", "com.test", "Test")

	releaseDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	// Test all setter methods
	app.SetArtistName("Artist")
	app.SetArtistID("456")
	app.SetVersion("2.0", 789)
	app.SetPrice(19.99, "€")
	app.SetRating(4.8, 5000)
	app.SetReleaseDate(releaseDate)
	app.SetGenre("g1", "Games")
	app.SetDeviceFamilies([]string{"iphone", "ipad"})
	app.SetFileSize(100000000)
	app.SetMinimumOSVersion("15.0")
	app.SetDescription("A test application")
	app.SetIconURL("https://example.com/icon.png")
	app.SetScreenshotURLs([]string{"https://example.com/screen1.png"})

	// Verify all values
	if app.ArtistName() != "Artist" {
		t.Error("ArtistName not set")
	}

	if app.ArtistID() != "456" {
		t.Error("ArtistID not set")
	}

	if app.Version() != "2.0" {
		t.Error("Version not set")
	}

	if app.VersionID() != 789 {
		t.Error("VersionID not set")
	}

	if app.Price() != 19.99 {
		t.Error("Price not set")
	}

	if app.Currency() != "€" {
		t.Error("Currency not set")
	}

	if app.Rating() != 4.8 {
		t.Error("Rating not set")
	}

	if app.RatingCount() != 5000 {
		t.Error("RatingCount not set")
	}

	if !app.ReleaseDate().Equal(releaseDate) {
		t.Error("ReleaseDate not set")
	}

	if app.GenreID() != "g1" {
		t.Error("GenreID not set")
	}

	if app.GenreName() != "Games" {
		t.Error("GenreName not set")
	}

	if len(app.DeviceFamilies()) != 2 {
		t.Error("DeviceFamilies not set")
	}

	if app.FileSize() != 100000000 {
		t.Error("FileSize not set")
	}

	if app.MinimumOSVersion() != "15.0" {
		t.Error("MinimumOSVersion not set")
	}

	if app.Description() != "A test application" {
		t.Error("Description not set")
	}

	if app.IconURL() != "https://example.com/icon.png" {
		t.Error("IconURL not set")
	}

	if len(app.ScreenshotURLs()) != 1 {
		t.Error("ScreenshotURLs not set")
	}
}
