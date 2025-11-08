package appstore

import "errors"

var (
	// ErrCredentialsDoNotSupportPurchasing is returned when credentials do not support purchasing.
	ErrCredentialsDoNotSupportPurchasing = errors.New("credentials do not support purchasing")

	// ErrApplicationRequiresRedownload is returned when application requires re-download.
	ErrApplicationRequiresRedownload = errors.New(
		"application requires re-download (STDRDL), which requires different kbsync certificate",
	)

	// ErrDownloadURLNotFound is returned when download URL is not found in response.
	ErrDownloadURLNotFound = errors.New("download URL not found in response")

	// ErrMultipleDownloadURLs is returned when multiple download URLs are found in response.
	ErrMultipleDownloadURLs = errors.New("unexpected: multiple download URLs in response")

	// ErrNoSINFFound is returned when no SINF is found in response.
	ErrNoSINFFound = errors.New("no SINF found")

	// ErrMultipleSINFs is returned when multiple SINFs are found in response.
	ErrMultipleSINFs = errors.New("unexpected: multiple SINFs in response")

	// ErrSINFEmpty is returned when SINF is empty.
	ErrSINFEmpty = errors.New("SINF is empty")

	// ErrBundleIDNotFound is returned when bundle ID is not found in response.
	ErrBundleIDNotFound = errors.New("bundle ID not found in response")

	// ErrEmptyPassword is returned when password is empty.
	ErrEmptyPassword = errors.New("password cannot be empty")

	// ErrUnexpectedStatusCode is returned when HTTP response has unexpected status code.
	ErrUnexpectedStatusCode = errors.New("unexpected status code")

	// ErrPasswordTokenNotFound is returned when password token is not found in response.
	ErrPasswordTokenNotFound = errors.New("password token not found in response")

	// ErrDSIDNotFound is returned when DSID is not found in response.
	ErrDSIDNotFound = errors.New("DSID not found in response")

	// ErrNoResultsFound is returned when no results are found.
	ErrNoResultsFound = errors.New("no results found")

	// ErrAdamIDNotFound is returned when Adam ID is not found in response.
	ErrAdamIDNotFound = errors.New("adamID not found in response")

	// ErrNoRatingFound is returned when no rating is found.
	ErrNoRatingFound = errors.New("no rating found")

	// ErrUnexpectedResponseStructure is returned when response structure is unexpected.
	ErrUnexpectedResponseStructure = errors.New("unexpected response structure")
)
