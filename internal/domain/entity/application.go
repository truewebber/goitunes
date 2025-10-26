package entity

import "time"

// Application represents an iOS application in the App Store
type Application struct {
	adamID           string
	bundleID         string
	name             string
	artistName       string
	artistID         string
	version          string
	versionID        int64
	price            float64
	currency         string
	rating           float64
	ratingCount      int
	releaseDate      time.Time
	genreID          string
	genreName        string
	deviceFamilies   []string
	fileSize         int64
	minimumOSVersion string
	description      string
	iconURL          string
	screenshotURLs   []string
}

// NewApplication creates a new Application entity
func NewApplication(adamID, bundleID, name string) *Application {
	return &Application{
		adamID:         adamID,
		bundleID:       bundleID,
		name:           name,
		deviceFamilies: make([]string, 0),
		screenshotURLs: make([]string, 0),
	}
}

// Getters
func (a *Application) AdamID() string           { return a.adamID }
func (a *Application) BundleID() string         { return a.bundleID }
func (a *Application) Name() string             { return a.name }
func (a *Application) ArtistName() string       { return a.artistName }
func (a *Application) ArtistID() string         { return a.artistID }
func (a *Application) Version() string          { return a.version }
func (a *Application) VersionID() int64         { return a.versionID }
func (a *Application) Price() float64           { return a.price }
func (a *Application) Currency() string         { return a.currency }
func (a *Application) Rating() float64          { return a.rating }
func (a *Application) RatingCount() int         { return a.ratingCount }
func (a *Application) ReleaseDate() time.Time   { return a.releaseDate }
func (a *Application) GenreID() string          { return a.genreID }
func (a *Application) GenreName() string        { return a.genreName }
func (a *Application) DeviceFamilies() []string { return a.deviceFamilies }
func (a *Application) FileSize() int64          { return a.fileSize }
func (a *Application) MinimumOSVersion() string { return a.minimumOSVersion }
func (a *Application) Description() string      { return a.description }
func (a *Application) IconURL() string          { return a.iconURL }
func (a *Application) ScreenshotURLs() []string { return a.screenshotURLs }

// Setters (for builder pattern)
func (a *Application) SetArtistName(name string) *Application {
	a.artistName = name
	return a
}

func (a *Application) SetArtistID(id string) *Application {
	a.artistID = id
	return a
}

func (a *Application) SetVersion(version string, versionID int64) *Application {
	a.version = version
	a.versionID = versionID
	return a
}

func (a *Application) SetPrice(price float64, currency string) *Application {
	a.price = price
	a.currency = currency
	return a
}

func (a *Application) SetRating(rating float64, count int) *Application {
	a.rating = rating
	a.ratingCount = count
	return a
}

func (a *Application) SetReleaseDate(date time.Time) *Application {
	a.releaseDate = date
	return a
}

func (a *Application) SetGenre(id, name string) *Application {
	a.genreID = id
	a.genreName = name
	return a
}

func (a *Application) SetDeviceFamilies(families []string) *Application {
	a.deviceFamilies = families
	return a
}

func (a *Application) SetFileSize(size int64) *Application {
	a.fileSize = size
	return a
}

func (a *Application) SetMinimumOSVersion(version string) *Application {
	a.minimumOSVersion = version
	return a
}

func (a *Application) SetDescription(desc string) *Application {
	a.description = desc
	return a
}

func (a *Application) SetIconURL(url string) *Application {
	a.iconURL = url
	return a
}

func (a *Application) SetScreenshotURLs(urls []string) *Application {
	a.screenshotURLs = urls
	return a
}

// IsFree returns true if the application is free
func (a *Application) IsFree() bool {
	return a.price == 0
}

// IsUniversal returns true if the application supports both iPhone and iPad
func (a *Application) IsUniversal() bool {
	hasIPhone := false
	hasIPad := false
	for _, family := range a.deviceFamilies {
		if family == "iphone" {
			hasIPhone = true
		}
		if family == "ipad" {
			hasIPad = true
		}
	}
	return hasIPhone && hasIPad
}
