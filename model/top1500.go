package model

type (
	Top1500ApplicationsResponse struct {
		PageNumber     string                                   `json:"pageNumber"`
		ContentData    []Top1500ApplicationsResponseContentData `json:"contentData"`
		HasNextPage    bool                                     `json:"hasNextPage"`
		NextStartIndex int                                      `json:"nextStartIndex"`
	}

	Top1500ApplicationsResponseContentData struct {
		ID                        string `json:"id"`
		Kind                      string `json:"kind"`
		Genre                     string `json:"genre"`
		URL                       string `json:"url"`
		ArtworkURL                string `json:"artwork_url"`
		ArtworkHeight             int    `json:"artwork_height"`
		ArtworkWidth              int    `json:"artwork_width"`
		Artwork2XURL              string `json:"artwork_2x_url"`
		Artwork2XHeight           int    `json:"artwork_2x_height"`
		Artwork2XWidth            int    `json:"artwork_2x_width"`
		Name                      string `json:"name"`
		ReleaseDate               string `json:"release_date"`
		ParentalControlAttributes struct {
			RatingSoftware string `json:"rating-software"`
			ParentalRating string `json:"parental-rating"`
		} `json:"parental_control_attributes"`
		UserRating      string `json:"user_rating"`
		UserRatingCount string `json:"user_rating_count"`
		BuyData         struct {
			ConfirmText          string `json:"confirm-text"`
			ActionParams         string `json:"action-params"`
			ItemTitle            string `json:"item-title"`
			ArtistName           string `json:"artist-name"`
			ItemID               string `json:"item-id"`
			ItemType             string `json:"item-type"`
			ArtworkURL           string `json:"artwork-url"`
			BundleID             string `json:"bundle-id"`
			IconIsPrerendered    string `json:"icon-is-prerendered"`
			RequiredCapabilities string `json:"required-capabilities"`
			MinimumOsVersion     string `json:"minimum-os-version"`
			VersionID            string `json:"versionID"`
			FileSize             string `json:"file-size"`
		} `json:"buyData"`
		ButtonText     string `json:"button_text"`
		IsUniversalApp string `json:"is_universal_app,omitempty"`
	}
)
