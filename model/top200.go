package model

import (
	"time"
)

type (
	Top200AppItem struct {
		Artwork struct {
			Width                int    `json:"width"`
			URL                  string `json:"url"`
			Height               int    `json:"height"`
			TextColor3           string `json:"textColor3"`
			TextColor2           string `json:"textColor2"`
			TextColor4           string `json:"textColor4"`
			TextColor1           string `json:"textColor1"`
			BgColor              string `json:"bgColor"`
			HasP3                bool   `json:"hasP3"`
			SupportsLayeredImage bool   `json:"supportsLayeredImage"`
		} `json:"artwork"`
		ArtistName             string                 `json:"artistName"`
		URL                    string                 `json:"url"`
		ShortURL               string                 `json:"shortUrl"`
		DeviceFamilies         []string               `json:"deviceFamilies"`
		GenreNames             []string               `json:"genreNames"`
		ID                     string                 `json:"id"`
		ReleaseDate            string                 `json:"releaseDate"`
		UserRating             UserRating             `json:"userRating"`
		ContentRatingsBySystem ContentRatingsBySystem `json:"contentRatingsBySystem"`
		Name                   string                 `json:"name"`
		ArtistURL              string                 `json:"artistUrl"`
		NameRaw                string                 `json:"nameRaw"`
		EditorialArtwork       struct {
			OriginalFlowcaseBrick struct {
				Width                int    `json:"width"`
				URL                  string `json:"url"`
				Height               int    `json:"height"`
				TextColor3           string `json:"textColor3"`
				TextColor2           string `json:"textColor2"`
				TextColor4           string `json:"textColor4"`
				TextColor1           string `json:"textColor1"`
				BgColor              string `json:"bgColor"`
				HasP3                bool   `json:"hasP3"`
				SupportsLayeredImage bool   `json:"supportsLayeredImage"`
			} `json:"originalFlowcaseBrick"`
			StoreFlowcase struct {
				Width                int    `json:"width"`
				URL                  string `json:"url"`
				Height               int    `json:"height"`
				TextColor3           string `json:"textColor3"`
				TextColor2           string `json:"textColor2"`
				TextColor4           string `json:"textColor4"`
				TextColor1           string `json:"textColor1"`
				BgColor              string `json:"bgColor"`
				HasP3                bool   `json:"hasP3"`
				SupportsLayeredImage bool   `json:"supportsLayeredImage"`
			} `json:"storeFlowcase"`
			BannerUber struct {
				Width                int    `json:"width"`
				URL                  string `json:"url"`
				Height               int    `json:"height"`
				TextColor3           string `json:"textColor3"`
				TextColor2           string `json:"textColor2"`
				TextColor4           string `json:"textColor4"`
				TextColor1           string `json:"textColor1"`
				BgColor              string `json:"bgColor"`
				HasP3                bool   `json:"hasP3"`
				SupportsLayeredImage bool   `json:"supportsLayeredImage"`
			} `json:"bannerUber"`
			FullscreenBackground struct {
				Width                int    `json:"width"`
				URL                  string `json:"url"`
				Height               int    `json:"height"`
				TextColor3           string `json:"textColor3"`
				TextColor2           string `json:"textColor2"`
				TextColor4           string `json:"textColor4"`
				TextColor1           string `json:"textColor1"`
				BgColor              string `json:"bgColor"`
				HasP3                bool   `json:"hasP3"`
				SupportsLayeredImage bool   `json:"supportsLayeredImage"`
			} `json:"fullscreenBackground"`
			ContentLogo struct {
				Width                int    `json:"width"`
				URL                  string `json:"url"`
				Height               int    `json:"height"`
				TextColor3           string `json:"textColor3"`
				TextColor2           string `json:"textColor2"`
				TextColor4           string `json:"textColor4"`
				TextColor1           string `json:"textColor1"`
				BgColor              string `json:"bgColor"`
				HasP3                bool   `json:"hasP3"`
				SupportsLayeredImage bool   `json:"supportsLayeredImage"`
			} `json:"contentLogo"`
		} `json:"editorialArtwork"`
		BundleID        string     `json:"bundleId"`
		Kind            string     `json:"kind"`
		Copyright       string     `json:"copyright"`
		ArtistID        string     `json:"artistId"`
		AgeBand         AgeBand    `json:"ageBand"`
		Genres          []AppGenre `json:"genres"`
		IsSiriSupported bool       `json:"isSiriSupported"`
		Offers          []Offer    `json:"offers"`
	}

	Top200ApplicationsResponse struct {
		StorePlatformData struct {
			Lockup struct {
				Results map[string]Top200AppItem `json:"results"`
				Meta    struct {
					Storefront struct {
						ID string `json:"id"`
						Cc string `json:"cc"`
					} `json:"storefront"`
					Language struct {
						Tag string `json:"tag"`
					} `json:"language"`
				} `json:"meta"`
			} `json:"lockup"`
		} `json:"storePlatformData"`
		PageData struct {
			ComponentName    string `json:"componentName"`
			PageTitle        string `json:"pageTitle"`
			SegmentedControl struct {
				Segments      []Top200PageSegment `json:"segments"`
				SelectedIndex int                 `json:"selectedIndex"`
			} `json:"segmentedControl"`
		} `json:"pageData"`
		Properties struct {
			DI6TopChartsPageNumIdsPerChart int       `json:"DI6.TopChartsPage.NumIdsPerChart"`
			RevNum                         string    `json:"revNum"`
			Timestamp                      time.Time `json:"timestamp"`
		} `json:"properties"`
	}

	Top200Chart struct {
		AdamIds   []string `json:"adamIds"`
		SeeAllURL string   `json:"seeAllUrl"`
		Kinds     struct {
			IosSoftware bool `json:"iosSoftware"`
		} `json:"kinds"`
		ShortTitle          string `json:"shortTitle"`
		ID                  string `json:"id"`
		Title               string `json:"title"`
		IsTracklistChart    bool   `json:"isTracklistChart"`
		IneligibleGratisIds []int  `json:"ineligibleGratisIds,omitempty"`
	}

	Top200PageSegment struct {
		PageData struct {
			MetricsBase struct {
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
			TopCharts     []Top200Chart `json:"topCharts"`
			SelectedChart Top200Chart   `json:"selectedChart"`
			PageTitle     string        `json:"pageTitle"`
			Genre         struct {
				Name string `json:"name"`
				ID   string `json:"id"`
				URL  string `json:"url"`
			} `json:"genre"`
			CategoryList struct {
				URL        string `json:"url"`
				ChartURL   string `json:"chartUrl"`
				Name       string `json:"name"`
				ShortName  string `json:"shortName"`
				ButtonName string `json:"buttonName"`
				GenreID    string `json:"genreId"`
				Art        struct {
					Token string `json:"token"`
					Pft   string `json:"pft"`
				} `json:"art"`
				Artwork []struct {
					Width  int    `json:"width"`
					URL    string `json:"url"`
					Height int    `json:"height"`
				} `json:"artwork"`
				Kind                string `json:"kind"`
				ParentCategoryLabel string `json:"parentCategoryLabel"`
				Children            []struct {
					URL       string `json:"url"`
					ChartURL  string `json:"chartUrl"`
					Name      string `json:"name"`
					ShortName string `json:"shortName"`
					GenreID   string `json:"genreId,omitempty"`
					Art       struct {
						Token string `json:"token"`
						Pft   string `json:"pft"`
					} `json:"art"`
					Artwork []struct {
						Width  int    `json:"width"`
						URL    string `json:"url"`
						Height int    `json:"height"`
					} `json:"artwork"`
					Kind                string `json:"kind"`
					GroupingID          string `json:"groupingId,omitempty"`
					ParentCategoryLabel string `json:"parentCategoryLabel,omitempty"`
					Children            []struct {
						URL        string `json:"url"`
						ChartURL   string `json:"chartUrl"`
						Name       string `json:"name"`
						ShortName  string `json:"shortName"`
						GenreID    string `json:"genreId"`
						GroupingID string `json:"groupingId"`
						Art        struct {
							Token string `json:"token"`
							Pft   string `json:"pft"`
						} `json:"art"`
						Artwork []struct {
							Width  int    `json:"width"`
							URL    string `json:"url"`
							Height int    `json:"height"`
						} `json:"artwork"`
						Kind                string `json:"kind"`
						ParentCategoryLabel string `json:"parentCategoryLabel"`
					} `json:"children,omitempty"`
				} `json:"children"`
			} `json:"categoryList"`
			ComponentName string `json:"componentName"`
			Metrics       struct {
				Config struct {
				} `json:"config"`
				Fields struct {
				} `json:"fields"`
			} `json:"metrics"`
		} `json:"pageData"`
		Title string `json:"title"`
		URL   string `json:"url"`
	}
)
