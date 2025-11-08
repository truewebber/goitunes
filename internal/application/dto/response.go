package dto

// GetTopChartsResponse represents the response for getting top charts.
type GetTopChartsResponse struct {
	Items      []ChartItemDTO `json:"items"`
	TotalCount int            `json:"totalCount"`
}

// GetApplicationInfoResponse represents the response for getting application info.
type GetApplicationInfoResponse struct {
	Applications []ApplicationDTO `json:"applications"`
}

// GetRatingResponse represents the response for getting rating info.
type GetRatingResponse struct {
	Rating      float64 `json:"rating"`
	RatingCount int     `json:"ratingCount"`
}

// AuthenticateResponse represents the authentication response.
type AuthenticateResponse struct {
	AppleID       string `json:"appleId"`
	PasswordToken string `json:"passwordToken"`
	DSID          string `json:"dsid"`
	Authenticated bool   `json:"authenticated"`
}

// PurchaseResponse represents the purchase response.
type PurchaseResponse struct {
	DownloadInfo DownloadInfoDTO `json:"downloadInfo"`
}
