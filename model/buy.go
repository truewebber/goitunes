package model

type (
	SongSinf struct {
		Id   int    `plist:"id"`
		Sinf []byte `plist:"sinf"`
	}

	SongRating struct {
		Content string `plist:"content"`
		Label   string `plist:"label"`
		Rank    int    `plist:"rank"`
		System  string `plist:"system"`
	}

	SongMetaData struct {
		DisplayBundleId            string     `plist:"bundleDisplayName"`
		BundleId                   string     `plist:"softwareVersionBundleId"`
		Q                          string     `plist:"q"`
		ArtistId                   int64      `plist:"artistId"`
		ArtistName                 string     `plist:"artistName"`
		BundleShortVersionString   string     `plist:"bundleShortVersionString"`
		BundleVersion              string     `plist:"bundleVersion"`
		Copyright                  string     `plist:"copyright"`
		Genre                      string     `plist:"genre"`
		GenreId                    int        `plist:"genreId"`
		ItemId                     int64      `plist:"itemId"`
		ItemName                   string     `plist:"itemName"`
		PlaylistName               string     `plist:"playlistName"`
		Rating                     SongRating `plist:"rating"`
		ReleaseDate                string     `plist:"releaseDate"`
		SoftwareIcon57x57URL       string     `plist:"softwareIcon57x57URL"`
		SoftwareSupportedDeviceIds []int      `plist:"softwareSupportedDeviceIds"`
		ExternalVersionId          int64      `plist:"softwareVersionExternalIdentifier"`
		ExternalVersionIdList      []int64    `plist:"softwareVersionExternalIdentifiers"`
		VendorId                   int64      `plist:"vendorId"`
		DrmVersionNumber           int        `plist:"drmVersionNumber"`
		VersionRestrictions        int64      `plist:"versionRestrictions"`
	}

	SongPlistSlice struct {
		SongId       int64        `plist:"songId"`
		Url          string       `plist:"URL"`
		DownloadKey  string       `plist:"downloadKey"`
		Sinfs        []SongSinf   `plist:"sinfs"`
		PurchaseDate string       `plist:"purchaseDate"`
		DownloadId   string       `plist:"download-id"`
		MetaData     SongMetaData `plist:"metadata"`
	}

	Metrics struct {
		DialogId    string `plist:"dialogId"`
		MtRequestId string `plist:"mtRequestId"`
	}

	BuyProductResponse struct {
		SongPlist []SongPlistSlice `plist:"songList"`
		Metrics   Metrics          `plist:"metrics"`
	}
)
