package config

// URL constants for App Store API endpoints
const (
	// Public API endpoints
	Top200AppsURL               = "https://itunes.apple.com/WebObjects/MZStore.woa/wa/viewTop"
	TopAppsURL                  = "https://itunes.apple.com/WebObjects/MZStore.woa/wa/topChartFragmentData"
	AppInfoURL                  = "https://uclient-api.itunes.apple.com/WebObjects/MZStorePlatform.woa/wa/lookup"
	NativeAppInfoURL            = "https://itunes.apple.com/app/id%s?mt=8"
	NativeAppRatingInfoURL      = "https://itunes.apple.com/customer-reviews/id%s?dataOnly=true&displayable-kind=11"
	OpenAppOverAllRatingInfoURL = "https://itunes.apple.com/lookup?id=%s&entity=software&country=%s"

	// Authenticated API endpoints (require login)
	LoginURLTemplate        = "https://p%d-buy.itunes.apple.com/WebObjects/MZFinance.woa/wa/authenticate"
	BuyProductURLTemplate   = "https://p%d-buy.itunes.apple.com/WebObjects/MZBuy.woa/wa/buyProduct"
	ConfirmDownloadTemplate = "https://p%d-buy.itunes.apple.com/WebObjects/MZFastFinance.woa/wa/songDownloadDone"
)

// Device codes for X-Apple-Store-Front header
const (
	IPhoneDeviceCode = 29
	IPadDeviceCode   = 32
)

// Chart type identifiers
const (
	// iPhone charts
	PopIDTopFree     = "27"
	PopIDTopPaid     = "30"
	PopIDTopGrossing = "38"

	// iPad charts
	PopIDIPadTopFree     = "44"
	PopIDIPadTopPaid     = "47"
	PopIDIPadTopGrossing = "46"
)

// Genre identifiers
const (
	GenreIDAll = "36"
)

