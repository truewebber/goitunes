package model

import (
	"time"
)

type (
	AppGenre struct {
		GenreID   string `json:"genreId"`
		Name      string `json:"name"`
		URL       string `json:"url"`
		MediaType string `json:"mediaType"`
	}

	///////////////////////////////////

	AgeBand struct {
		MinAge int `json:"minAge"`
		MaxAge int `json:"maxAge"`
	}

	///////////////////////////////////

	Artwork struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		URL    string `json:"url"`
	}

	///////////////////////////////////

	Offer struct {
		ActionText     OfferActionText `json:"actionText"`
		PriceFormatted string          `json:"priceFormatted"`
		Price          float64         `json:"price"`
		Type           string          `json:"type"`
		BuyParams      string          `json:"buyParams"`
		Version        OfferVersion    `json:"version"`
		Assets         []OfferAsset    `json:"assets"`
	}

	OfferActionText struct {
		Short       string `json:"short"`
		Medium      string `json:"medium"`
		Long        string `json:"long"`
		Downloaded  string `json:"downloaded"`
		Downloading string `json:"downloading"`
	}

	OfferVersion struct {
		Display    string `json:"display"`
		ExternalID int    `json:"externalId"`
	}

	OfferAsset struct {
		Flavor string `json:"flavor"`
		Size   int    `json:"size"`
	}

	///////////////////////////////////

	ContentRatingsBySystem struct {
		AppsApple ContentRatingsBySystemAppsApple `json:"appsApple"`
	}

	ContentRatingsBySystemAppsApple struct {
		Name       string   `json:"name"`
		Value      int      `json:"value"`
		Rank       int      `json:"rank"`
		Advisories []string `json:"advisories"`
	}

	ContentRating struct {
		System     string   `json:"system"`
		Name       string   `json:"name"`
		Value      int      `json:"value"`
		Rank       int      `json:"rank"`
		Advisories []string `json:"advisories"`
	}

	UserRating struct {
		Value                     float64 `json:"value"`
		RatingCount               int     `json:"ratingCount"`
		ValueCurrentVersion       float64 `json:"valueCurrentVersion"`
		RatingCountCurrentVersion int     `json:"ratingCountCurrentVersion"`
	}

	UserRatingV2 struct {
		Value               float64 `json:"value"`
		RatingCount         int     `json:"ratingCount"`
		RatingCountList     []int   `json:"ratingCountList"`
		AriaLabelForRatings string  `json:"ariaLabelForRatings"`
	}

	///////////////////////////////////

	ApplicationInfoByBundleIdResponse struct {
		ResultCount int                `json:"resultCount"`
		Results     map[string]AppItem `json:"results"`
	}

	///////////////////////////////////

	AppItemFullResponse struct {
		StorePlatformData struct {
			ProductDv struct {
				Results         map[string]*AppItemFull `json:"results"`
				Version         int                     `json:"version"`
				IsAuthenticated bool                    `json:"isAuthenticated"`
				Meta            struct {
					Storefront struct {
						ID string `json:"id"`
						Cc string `json:"cc"`
					} `json:"storefront"`
					Language struct {
						Tag string `json:"tag"`
					} `json:"language"`
				} `json:"meta"`
			} `json:"product-dv"`
		} `json:"storePlatformData"`
		PageData struct {
			ComponentName string `json:"componentName"`
			MetricsBase   struct {
				PageType              string `json:"pageType"`
				PageID                string `json:"pageId"`
				PageDetails           string `json:"pageDetails"`
				Page                  string `json:"page"`
				ServerInstance        string `json:"serverInstance"`
				StoreFrontHeader      string `json:"storeFrontHeader"`
				Language              string `json:"language"`
				PlatformID            string `json:"platformId"`
				PlatformName          string `json:"platformName"`
				StoreFront            string `json:"storeFront"`
				EnvironmentDataCenter string `json:"environmentDataCenter"`
			} `json:"metricsBase"`
			Metrics struct {
				Config struct {
				} `json:"config"`
				Fields struct {
					IsUber bool `json:"isUber"`
				} `json:"fields"`
			} `json:"metrics"`
			TopApps struct {
				Iphone struct {
					Ids   []string `json:"ids"`
					Title string   `json:"title"`
				} `json:"iphone"`
			} `json:"topApps"`
			RatingAndAdvisories struct {
				Advisories []string `json:"advisories"`
				RatingText string   `json:"rating-text"`
			} `json:"rating-and-advisories"`
			VersionHistory []struct {
				ReleaseNotes  string    `json:"releaseNotes"`
				VersionString string    `json:"versionString"`
				ReleaseDate   time.Time `json:"releaseDate"`
			} `json:"versionHistory"`
			AddOns []struct {
				BuyParams string `json:"buyParams"`
				OfferType struct {
					OfferType       string `json:"offerType"`
					AccompaniesNoun bool   `json:"accompaniesNoun"`
				} `json:"offerType"`
				Price string `json:"price"`
				Name  string `json:"name"`
			} `json:"addOns"`
			CustomersAlsoBoughtApps []string `json:"customersAlsoBoughtApps"`
			ID                      string   `json:"id"`
			AppRatingsLearnMoreURL  string   `json:"appRatingsLearnMoreUrl"`
			SellerLabel             string   `json:"sellerLabel"`
			MoreByThisDeveloper     []string `json:"moreByThisDeveloper"`
			Uber                    struct {
				TitleTextColor   string `json:"titleTextColor"`
				BackgroundColor  string `json:"backgroundColor"`
				PrimaryTextColor string `json:"primaryTextColor"`
				MasterArt        []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"masterArt"`
			} `json:"uber"`
			CustomerReviewsURL   string `json:"customerReviewsUrl"`
			Sf6ResourceImagePath string `json:"sf6ResourceImagePath"`
		} `json:"pageData"`
		Properties struct {
			RevNum    string    `json:"revNum"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"properties"`
	}

	//AppRating

	ApplicationRatingResponse struct {
		AdamID                 int    `json:"adamId"`
		ClickToRateURL         string `json:"clickToRateUrl"`
		WriteUserReviewURL     string `json:"writeUserReviewUrl"`
		TotalNumberOfReviews   int    `json:"totalNumberOfReviews"`
		UserReviewsRowURL      string `json:"userReviewsRowUrl"`
		UserReviewsSortOptions []struct {
			SortID int    `json:"sortId"`
			Name   string `json:"name"`
		} `json:"userReviewsSortOptions"`
		KindID              int     `json:"kindId"`
		KindExtID           string  `json:"kindExtId"`
		KindName            string  `json:"kindName"`
		SaveUserReviewURL   string  `json:"saveUserReviewUrl"`
		AriaLabelForRatings string  `json:"ariaLabelForRatings"`
		RatingCount         int     `json:"ratingCount"`
		RatingCountList     []int   `json:"ratingCountList"`
		RatingAverage       float64 `json:"ratingAverage"`
		WasReset            bool    `json:"wasReset"`
	}

	//AppOverAllRating

	ApplicationOverAllRatingResponse struct {
		ResultCount int                             `json:"resultCount"`
		Results     []*ApplicationOverAllRatingItem `json:"results"`
	}

	ApplicationOverAllRatingItem struct {
		ScreenshotUrls                     []string      `json:"screenshotUrls"`
		IpadScreenshotUrls                 []interface{} `json:"ipadScreenshotUrls"`
		AppletvScreenshotUrls              []interface{} `json:"appletvScreenshotUrls"`
		ArtworkURL512                      string        `json:"artworkUrl512"`
		ArtistViewURL                      string        `json:"artistViewUrl"`
		ArtworkURL60                       string        `json:"artworkUrl60"`
		ArtworkURL100                      string        `json:"artworkUrl100"`
		Advisories                         []string      `json:"advisories"`
		IsGameCenterEnabled                bool          `json:"isGameCenterEnabled"`
		SupportedDevices                   []string      `json:"supportedDevices"`
		Kind                               string        `json:"kind"`
		Features                           []interface{} `json:"features"`
		LanguageCodesISO2A                 []string      `json:"languageCodesISO2A"`
		AverageUserRatingForCurrentVersion float64       `json:"averageUserRatingForCurrentVersion"`
		UserRatingCountForCurrentVersion   int           `json:"userRatingCountForCurrentVersion"`
		TrackContentRating                 string        `json:"trackContentRating"`
		FileSizeBytes                      string        `json:"fileSizeBytes"`
		TrackViewURL                       string        `json:"trackViewUrl"`
		ContentAdvisoryRating              string        `json:"contentAdvisoryRating"`
		TrackCensoredName                  string        `json:"trackCensoredName"`
		IsVppDeviceBasedLicensingEnabled   bool          `json:"isVppDeviceBasedLicensingEnabled"`
		TrackID                            int           `json:"trackId"`
		TrackName                          string        `json:"trackName"`
		PrimaryGenreName                   string        `json:"primaryGenreName"`
		ReleaseDate                        time.Time     `json:"releaseDate"`
		FormattedPrice                     string        `json:"formattedPrice"`
		Currency                           string        `json:"currency"`
		WrapperType                        string        `json:"wrapperType"`
		Version                            string        `json:"version"`
		ArtistID                           int           `json:"artistId"`
		ArtistName                         string        `json:"artistName"`
		Genres                             []string      `json:"genres"`
		Price                              float64       `json:"price"`
		Description                        string        `json:"description"`
		BundleID                           string        `json:"bundleId"`
		PrimaryGenreID                     int           `json:"primaryGenreId"`
		SellerName                         string        `json:"sellerName"`
		GenreIds                           []string      `json:"genreIds"`
		MinimumOsVersion                   string        `json:"minimumOsVersion"`
		CurrentVersionReleaseDate          time.Time     `json:"currentVersionReleaseDate"`
		ReleaseNotes                       string        `json:"releaseNotes"`
		AverageUserRating                  float64       `json:"averageUserRating"`
		UserRatingCount                    int           `json:"userRatingCount"`
	}
)
