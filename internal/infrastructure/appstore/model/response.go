package model

// Top200Response represents the response from Top200 API.
type Top200Response struct {
	StorePlatformData struct {
		Lockup struct {
			Results map[string]AppItemResponse `json:"results"`
		} `json:"lockup"`
	} `json:"storePlatformData"`
	PageData struct {
		SegmentedControl struct {
			Segments []struct {
				PageData struct {
					SelectedChart struct {
						AdamIDs []string `json:"adamIds"`
					} `json:"selectedChart"`
				} `json:"pageData"`
			} `json:"segments"`
			SelectedIndex int `json:"selectedIndex"`
		} `json:"segmentedControl"`
	} `json:"pageData"`
	Properties struct {
		DI6TopChartsPageNumIDsPerChart int `json:"di6-top-charts-page-num-ids-per-chart"`
	} `json:"properties"`
}

// Top1500Response represents the response from Top1500 API.
type Top1500Response struct {
	ContentData []struct {
		ID         string `json:"id"`
		UserRating string `json:"userRating"`
		ButtonText string `json:"buttonText"`
		BuyData    struct {
			BundleID     string `json:"bundleId"`
			VersionID    string `json:"versionId"`
			ActionParams string `json:"actionParams"`
		} `json:"buyData"`
	} `json:"contentData"`
}

// AppItemResponse represents application information.
type AppItemResponse struct {
	ScreenshotsByType map[string][]struct {
		URL string `json:"url"`
	} `json:"screenshotsByType"`
	FileSizeByDevice map[string]int `json:"fileSizeByDevice"`
	MinimumOSVersion string         `json:"minimumOSVersion"`
	Description      struct {
		Standard string `json:"standard"`
	} `json:"description"`
	ArtistName  string `json:"artistName"`
	ArtistID    string `json:"artistId"`
	ArtistURL   string `json:"artistUrl"`
	BundleID    string `json:"bundleId"`
	ReleaseDate string `json:"releaseDate"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	NameRaw     string `json:"nameRaw"`
	Genres      []struct {
		GenreID   string `json:"genreId"`
		Name      string `json:"name"`
		URL       string `json:"url"`
		MediaType string `json:"mediaType"`
	} `json:"genres"`
	GenreNames     []string `json:"genreNames"`
	DeviceFamilies []string `json:"deviceFamilies"`
	Offers         []struct {
		Type           string `json:"type"`
		PriceFormatted string `json:"priceFormatted"`
		BuyParams      string `json:"buyParams"`
		Version        struct {
			Display    string `json:"display"`
			ExternalID int    `json:"externalId"`
		} `json:"version"`
		Assets []struct {
			Flavor string `json:"flavor"`
			Size   int    `json:"size"`
		} `json:"assets"`
		Price float64 `json:"price"`
	} `json:"offers"`
	Artwork struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"artwork"`
	UserRating struct {
		Value       float64 `json:"value"`
		RatingCount int     `json:"ratingCount"`
	} `json:"userRating"`
}

// LookupResponse represents the response from lookup API.
type LookupResponse struct {
	Results map[string]AppItemResponse `json:"results"`
}

// FullAppResponse represents detailed application information.
type FullAppResponse struct {
	StorePlatformData struct {
		ProductDv struct {
			Results map[string]AppItemResponse `json:"results"`
		} `json:"product-dv"`
	} `json:"storePlatformData"`
}

// RatingResponse represents rating information.
type RatingResponse struct {
	AdamID     int `json:"adamId"`
	UserRating struct {
		Value       float64 `json:"value"`
		RatingCount int     `json:"ratingCount"`
	} `json:"userRating"`
}

// OverallRatingResponse represents overall rating from open API.
type OverallRatingResponse struct {
	Results []struct {
		AverageUserRating                  float64 `json:"averageUserRating"`
		UserRatingCount                    int     `json:"userRatingCount"`
		AverageUserRatingForCurrentVersion float64 `json:"averageUserRatingForCurrentVersion"`
		UserRatingCountForCurrentVersion   int     `json:"userRatingCountForCurrentVersion"`
	} `json:"results"`
}
