package valueobject

import "errors"

var (
	// ErrEmptyAppleID is returned when Apple ID is empty.
	ErrEmptyAppleID = errors.New("appleID cannot be empty")

	// ErrEmptyPasswordToken is returned when password token is empty.
	ErrEmptyPasswordToken = errors.New("passwordToken cannot be empty")

	// ErrEmptyDSID is returned when DSID is empty.
	ErrEmptyDSID = errors.New("dsid cannot be empty")

	// ErrEmptyGUID is returned when GUID is empty.
	ErrEmptyGUID = errors.New("guid cannot be empty")

	// ErrEmptyMachineName is returned when machine name is empty.
	ErrEmptyMachineName = errors.New("machineName cannot be empty")

	// ErrEmptyUserAgent is returned when user agent is empty.
	ErrEmptyUserAgent = errors.New("userAgent cannot be empty")

	// ErrEmptyRegion is returned when region is empty.
	ErrEmptyRegion = errors.New("region cannot be empty")

	// ErrInvalidStoreFront is returned when store front is not positive.
	ErrInvalidStoreFront = errors.New("storeFront must be positive")

	// ErrInvalidHostPrefix is returned when host prefix is not positive.
	ErrInvalidHostPrefix = errors.New("hostPrefix must be positive")
)

