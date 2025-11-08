package dto

// GetTopChartsRequest represents a request to get top charts.
type GetTopChartsRequest struct {
	GenreID    string
	ChartType  string // "topfree", "toppaid", "topgrossing"
	KidPrefix  string // Optional age band filter
	From       int    // Starting position (1-based)
	Limit      int    // Number of results
	MaxResults int    // For Top1500: page size
	Page       int    // For Top1500: page number (0-based)
}

// GetApplicationInfoRequest represents a request to get application info.
type GetApplicationInfoRequest struct {
	AdamIDs   []string
	BundleIDs []string
}

// GetRatingRequest represents a request to get rating information.
type GetRatingRequest struct {
	AdamID  string
	Overall bool // If true, get overall rating from open API
}

// AuthenticateRequest represents an authentication request.
type AuthenticateRequest struct {
	AppleID  string
	Password string
}

// PurchaseRequest represents a purchase request.
type PurchaseRequest struct {
	AdamID    string
	VersionID int64
}
