package model

// PurchaseResponse represents the response from buy/purchase API
type PurchaseResponse struct {
	SongList []SongItem `plist:"songList"`
	Metrics  struct {
		DialogID    string `plist:"dialogId"`
		MtRequestID string `plist:"mtRequestId"`
	} `plist:"metrics"`
}

// SongItem represents a downloadable item
type SongItem struct {
	SongID       int64    `plist:"songId"`
	URL          string   `plist:"URL"`
	DownloadKey  string   `plist:"downloadKey"`
	Sinfs        []Sinf   `plist:"sinfs"`
	PurchaseDate string   `plist:"purchaseDate"`
	DownloadID   string   `plist:"download-id"`
	Metadata     Metadata `plist:"metadata"`
}

// Sinf represents DRM information
type Sinf struct {
	ID   int    `plist:"id"`
	Data []byte `plist:"sinf"`
}

// Metadata represents application metadata
type Metadata struct {
	BundleDisplayName            string   `plist:"bundleDisplayName"`
	BundleID                     string   `plist:"softwareVersionBundleId"`
	Q                            string   `plist:"q"`
	ArtistID                     int64    `plist:"artistId"`
	ArtistName                   string   `plist:"artistName"`
	BundleShortVersionString     string   `plist:"bundleShortVersionString"`
	BundleVersion                string   `plist:"bundleVersion"`
	Copyright                    string   `plist:"copyright"`
	Genre                        string   `plist:"genre"`
	GenreID                      int      `plist:"genreId"`
	ItemID                       int64    `plist:"itemId"`
	ItemName                     string   `plist:"itemName"`
	PlaylistName                 string   `plist:"playlistName"`
	ReleaseDate                  string   `plist:"releaseDate"`
	SoftwareIcon57x57URL         string   `plist:"softwareIcon57x57URL"`
	SoftwareSupportedDeviceIDs   []int    `plist:"softwareSupportedDeviceIds"`
	ExternalVersionID            int64    `plist:"softwareVersionExternalIdentifier"`
	ExternalVersionIDList        []int64  `plist:"softwareVersionExternalIdentifiers"`
	VendorID                     int64    `plist:"vendorId"`
	DRMVersionNumber             int      `plist:"drmVersionNumber"`
	VersionRestrictions          int64    `plist:"versionRestrictions"`
	Rating                       Rating   `plist:"rating"`
}

// Rating represents content rating
type Rating struct {
	Content string `plist:"content"`
	Label   string `plist:"label"`
	Rank    int    `plist:"rank"`
	System  string `plist:"system"`
}

