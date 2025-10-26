package dto

import (
	"time"
)

// ApplicationDTO represents application data transfer object
type ApplicationDTO struct {
	AdamID           string    `json:"adamId"`
	BundleID         string    `json:"bundleId"`
	Name             string    `json:"name"`
	ArtistName       string    `json:"artistName"`
	ArtistID         string    `json:"artistId"`
	Version          string    `json:"version"`
	VersionID        int64     `json:"versionId"`
	Price            float64   `json:"price"`
	Currency         string    `json:"currency"`
	Rating           float64   `json:"rating"`
	RatingCount      int       `json:"ratingCount"`
	ReleaseDate      time.Time `json:"releaseDate"`
	GenreID          string    `json:"genreId"`
	GenreName        string    `json:"genreName"`
	DeviceFamilies   []string  `json:"deviceFamilies"`
	FileSize         int64     `json:"fileSize"`
	MinimumOSVersion string    `json:"minimumOsVersion"`
	Description      string    `json:"description"`
	IconURL          string    `json:"iconUrl"`
	ScreenshotURLs   []string  `json:"screenshotUrls"`
	IsFree           bool      `json:"isFree"`
	IsUniversal      bool      `json:"isUniversal"`
}

// ChartItemDTO represents a chart item data transfer object
type ChartItemDTO struct {
	Position int            `json:"position"`
	App      ApplicationDTO `json:"app"`
}

// DownloadInfoDTO represents download information data transfer object
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

