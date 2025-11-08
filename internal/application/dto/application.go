package dto

import "time"

// ApplicationDTO represents application data transfer object.
type ApplicationDTO struct {
	ReleaseDate      time.Time `json:"releaseDate"`
	Description      string    `json:"description"`
	Currency         string    `json:"currency"`
	ArtistName       string    `json:"artistName"`
	ArtistID         string    `json:"artistId"`
	Version          string    `json:"version"`
	GenreName        string    `json:"genreName"`
	IconURL          string    `json:"iconUrl"`
	Name             string    `json:"name"`
	MinimumOSVersion string    `json:"minimumOsVersion"`
	AdamID           string    `json:"adamId"`
	BundleID         string    `json:"bundleId"`
	GenreID          string    `json:"genreId"`
	ScreenshotURLs   []string  `json:"screenshotUrls"`
	DeviceFamilies   []string  `json:"deviceFamilies"`
	RatingCount      int       `json:"ratingCount"`
	FileSize         int64     `json:"fileSize"`
	Rating           float64   `json:"rating"`
	Price            float64   `json:"price"`
	VersionID        int64     `json:"versionId"`
	IsFree           bool      `json:"isFree"`
	IsUniversal      bool      `json:"isUniversal"`
}

// ChartItemDTO represents a chart item data transfer object.
type ChartItemDTO struct {
	App      ApplicationDTO `json:"app"`
	Position int            `json:"position"`
}

// DownloadInfoDTO represents download information data transfer object.
type DownloadInfoDTO struct {
	BundleID    string            `json:"bundleId"`
	URL         string            `json:"url"`
	DownloadKey string            `json:"downloadKey"`
	Sinf        string            `json:"sinf"`
	Metadata    string            `json:"metadata"`
	Headers     map[string]string `json:"headers"`
	DownloadID  string            `json:"downloadId"`
	VersionID   int64             `json:"versionId"`
	FileSize    int64             `json:"fileSize"`
}
