package model

import (
	"time"
)

type (
	AppItem struct {
		Artwork                Artwork                `json:"artwork"`
		ArtistName             string                 `json:"artistName"`
		URL                    string                 `json:"url"`
		ShortURL               string                 `json:"shortUrl"`
		DeviceFamilies         []string               `json:"deviceFamilies"`
		GenreNames             []string               `json:"genreNames"`
		NameSortValue          string                 `json:"nameSortValue"`
		ID                     string                 `json:"id"`
		ReleaseDate            string                 `json:"releaseDate"`
		UserRating             UserRatingV2           `json:"userRating"`
		ContentRatingsBySystem ContentRatingsBySystem `json:"contentRatingsBySystem"`
		Name                   string                 `json:"name"`
		ArtistURL              string                 `json:"artistUrl"`
		NameRaw                string                 `json:"nameRaw"`
		EditorialArtwork       struct{}               `json:"editorialArtwork"`
		BundleID               string                 `json:"bundleId"`
		Kind                   string                 `json:"kind"`
		Copyright              string                 `json:"copyright"`
		ArtistID               string                 `json:"artistId"`
		Genres                 []AppGenre             `json:"genres"`
		Uber                   interface{}            `json:"uber"`
		OvalArtwork            interface{}            `json:"ovalArtwork"`
		AgeBand                AgeBand                `json:"ageBand"`
		HasInAppPurchases      bool                   `json:"hasInAppPurchases"`
		Offers                 []Offer                `json:"offers"`
	}

	//

	TopAppItemResponse struct {
		AdamId        string
		BundleId      string
		Position      int
		Rating        float64
		Price         float64
		CurrencyLabel string
		Version       string
		VersionCode   int64
	}

	//

	AppItemFull struct {
		Artwork struct {
			Width                int    `json:"width"`
			URL                  string `json:"url"`
			Height               int    `json:"height"`
			TextColor3           string `json:"textColor3"`
			TextColor2           string `json:"textColor2"`
			TextColor4           string `json:"textColor4"`
			HasAlpha             bool   `json:"hasAlpha"`
			TextColor1           string `json:"textColor1"`
			BgColor              string `json:"bgColor"`
			HasP3                bool   `json:"hasP3"`
			SupportsLayeredImage bool   `json:"supportsLayeredImage"`
		} `json:"artwork"`
		ArtistName             string    `json:"artistName"`
		FamilyShareEnabledDate time.Time `json:"familyShareEnabledDate"`
		URL                    string    `json:"url"`
		ShortURL               string    `json:"shortUrl"`
		SoftwareInfo           struct {
			Seller                 string      `json:"seller"`
			LanguagesDisplayString string      `json:"languagesDisplayString"`
			RequirementsString     string      `json:"requirementsString"`
			EulaURL                interface{} `json:"eulaUrl"`
			SupportURL             string      `json:"supportUrl"`
			WebsiteURL             interface{} `json:"websiteUrl"`
			PrivacyPolicyURL       string      `json:"privacyPolicyUrl"`
			PrivacyPolicyTextURL   interface{} `json:"privacyPolicyTextUrl"`
			HasInAppPurchases      bool        `json:"hasInAppPurchases"`
		} `json:"softwareInfo"`
		TellAFriendMessageContentsURL string   `json:"tellAFriendMessageContentsUrl"`
		DeviceFamilies                []string `json:"deviceFamilies"`
		GenreNames                    []string `json:"genreNames"`
		AgeBand                       AgeBand  `json:"ageBand"`
		ItunesNotes                   struct {
			Short   string `json:"short"`
			Tagline string `json:"tagline"`
		} `json:"itunesNotes"`
		NameSortValue    string        `json:"nameSortValue"`
		ID               string        `json:"id"`
		AppBundleAdamIds []interface{} `json:"appBundleAdamIds"`
		ReleaseDate      string        `json:"releaseDate"`
		UserRating       struct {
			Value       float64 `json:"value"`
			RatingCount int     `json:"ratingCount"`
		} `json:"userRating"`
		ContentRatingsBySystem struct {
			AppsApple struct {
				Name       string   `json:"name"`
				Value      int      `json:"value"`
				Rank       int      `json:"rank"`
				Advisories []string `json:"advisories"`
			} `json:"appsApple"`
		} `json:"contentRatingsBySystem"`
		Name string `json:"name"`
		Uber struct {
			BackgroundColor  string      `json:"backgroundColor"`
			Name             interface{} `json:"name"`
			TitleTextColor   string      `json:"titleTextColor"`
			PrimaryTextColor string      `json:"primaryTextColor"`
			MasterArt        struct {
				Width                int    `json:"width"`
				URL                  string `json:"url"`
				Height               int    `json:"height"`
				TextColor3           string `json:"textColor3"`
				TextColor2           string `json:"textColor2"`
				TextColor4           string `json:"textColor4"`
				HasAlpha             bool   `json:"hasAlpha"`
				TextColor1           string `json:"textColor1"`
				BgColor              string `json:"bgColor"`
				HasP3                bool   `json:"hasP3"`
				SupportsLayeredImage bool   `json:"supportsLayeredImage"`
			} `json:"masterArt"`
			HeaderTextColor         string      `json:"headerTextColor"`
			PrimaryTextColorOnBlack string      `json:"primaryTextColorOnBlack"`
			TitleTextColorOnBlack   string      `json:"titleTextColorOnBlack"`
			Description             interface{} `json:"description"`
		} `json:"uber"`
		ArtistURL         string                              `json:"artistUrl"`
		ScreenshotsByType map[string][]AppItemFullScreenshots `json:"screenshotsByType"`
		NameRaw           string                              `json:"nameRaw"`
		EditorialArtwork  struct {
			StoreFlowcase struct {
				Width                int    `json:"width"`
				URL                  string `json:"url"`
				Height               int    `json:"height"`
				TextColor3           string `json:"textColor3"`
				TextColor2           string `json:"textColor2"`
				TextColor4           string `json:"textColor4"`
				HasAlpha             bool   `json:"hasAlpha"`
				TextColor1           string `json:"textColor1"`
				BgColor              string `json:"bgColor"`
				HasP3                bool   `json:"hasP3"`
				SupportsLayeredImage bool   `json:"supportsLayeredImage"`
			} `json:"storeFlowcase"`
			SubscriptionHero struct {
				Width                int    `json:"width"`
				URL                  string `json:"url"`
				Height               int    `json:"height"`
				TextColor3           string `json:"textColor3"`
				TextColor2           string `json:"textColor2"`
				TextColor4           string `json:"textColor4"`
				HasAlpha             bool   `json:"hasAlpha"`
				TextColor1           string `json:"textColor1"`
				BgColor              string `json:"bgColor"`
				HasP3                bool   `json:"hasP3"`
				SupportsLayeredImage bool   `json:"supportsLayeredImage"`
			} `json:"subscriptionHero"`
			BannerUber struct {
				Width                int    `json:"width"`
				URL                  string `json:"url"`
				Height               int    `json:"height"`
				TextColor3           string `json:"textColor3"`
				TextColor2           string `json:"textColor2"`
				TextColor4           string `json:"textColor4"`
				HasAlpha             bool   `json:"hasAlpha"`
				TextColor1           string `json:"textColor1"`
				BgColor              string `json:"bgColor"`
				HasP3                bool   `json:"hasP3"`
				SupportsLayeredImage bool   `json:"supportsLayeredImage"`
			} `json:"bannerUber"`
		} `json:"editorialArtwork"`
		Subtitle           string `json:"subtitle"`
		BundleID           string `json:"bundleId"`
		HasInAppPurchases  bool   `json:"hasInAppPurchases"`
		Kind               string `json:"kind"`
		Copyright          string `json:"copyright"`
		VideoPreviewByType struct {
		} `json:"videoPreviewByType"`
		ArtistID         string         `json:"artistId"`
		FileSizeByDevice map[string]int `json:"fileSizeByDevice"`
		Genres           []struct {
			GenreID   string `json:"genreId"`
			Name      string `json:"name"`
			URL       string `json:"url"`
			MediaType string `json:"mediaType"`
		} `json:"genres"`
		MinimumOSVersion    string `json:"minimumOSVersion"`
		MessagesScreenshots struct {
		} `json:"messagesScreenshots"`
		Description struct {
			Standard string `json:"standard"`
		} `json:"description"`
		Offers []struct {
			ActionText struct {
				Short       string `json:"short"`
				Medium      string `json:"medium"`
				Long        string `json:"long"`
				Downloaded  string `json:"downloaded"`
				Downloading string `json:"downloading"`
			} `json:"actionText"`
			Type           string  `json:"type"`
			PriceFormatted string  `json:"priceFormatted"`
			Price          float64 `json:"price"`
			BuyParams      string  `json:"buyParams"`
			Version        struct {
				Display    string `json:"display"`
				ExternalID int    `json:"externalId"`
			} `json:"version"`
			Assets []struct {
				Flavor string `json:"flavor"`
				Size   int    `json:"size"`
			} `json:"assets"`
		} `json:"offers"`
	}

	AppItemFullScreenshots struct {
		Width                int    `json:"width"`
		URL                  string `json:"url"`
		Height               int    `json:"height"`
		TextColor3           string `json:"textColor3"`
		TextColor2           string `json:"textColor2"`
		TextColor4           string `json:"textColor4"`
		HasAlpha             bool   `json:"hasAlpha"`
		TextColor1           string `json:"textColor1"`
		BgColor              string `json:"bgColor"`
		HasP3                bool   `json:"hasP3"`
		SupportsLayeredImage bool   `json:"supportsLayeredImage"`
	}
)
