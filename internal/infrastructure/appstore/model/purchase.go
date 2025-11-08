package model

// PurchaseResponse represents the response from buy/purchase API.
type PurchaseResponse struct {
	Metrics struct {
		DialogID    string `plist:"dialogId"`
		MtRequestID string `plist:"mtRequestId"`
	} `plist:"metrics"`
	SongList []SongItem `plist:"songList"`
}

// SongItem represents a downloadable item.
type SongItem struct {
	URL          string   `plist:"URL"`
	DownloadKey  string   `plist:"downloadKey"`
	PurchaseDate string   `plist:"purchaseDate"`
	DownloadID   string   `plist:"download-id"`
	Sinfs        []Sinf   `plist:"sinfs"`
	Metadata     Metadata `plist:"metadata"`
	SongID       int64    `plist:"songId"`
}

// Sinf represents DRM information.
type Sinf struct {
	Data []byte `plist:"sinf"`
	ID   int    `plist:"id"`
}

// Metadata represents application metadata.
type Metadata struct {
	Rating                     Rating  `plist:"rating"`
	BundleVersion              string  `plist:"bundleVersion"`
	Q                          string  `plist:"q"`
	BundleID                   string  `plist:"softwareVersionBundleId"`
	ArtistName                 string  `plist:"artistName"`
	BundleShortVersionString   string  `plist:"bundleShortVersionString"`
	ReleaseDate                string  `plist:"releaseDate"`
	Copyright                  string  `plist:"copyright"`
	Genre                      string  `plist:"genre"`
	BundleDisplayName          string  `plist:"bundleDisplayName"`
	SoftwareIcon57x57URL       string  `plist:"softwareIcon57x57URL"`
	ItemName                   string  `plist:"itemName"`
	PlaylistName               string  `plist:"playlistName"`
	SoftwareSupportedDeviceIDs []int   `plist:"softwareSupportedDeviceIds"`
	ExternalVersionIDList      []int64 `plist:"softwareVersionExternalIdentifiers"`
	ItemID                     int64   `plist:"itemId"`
	GenreID                    int     `plist:"genreId"`
	ExternalVersionID          int64   `plist:"softwareVersionExternalIdentifier"`
	VendorID                   int64   `plist:"vendorId"`
	DRMVersionNumber           int     `plist:"drmVersionNumber"`
	VersionRestrictions        int64   `plist:"versionRestrictions"`
	ArtistID                   int64   `plist:"artistId"`
}

// Rating represents content rating.
type Rating struct {
	Content string `plist:"content"`
	Label   string `plist:"label"`
	System  string `plist:"system"`
	Rank    int    `plist:"rank"`
}
