package usecase

import "errors"

var (
	// ErrEmptyAppleID is returned when Apple ID is empty.
	ErrEmptyAppleID = errors.New("appleID cannot be empty")

	// ErrEmptyPassword is returned when password is empty.
	ErrEmptyPassword = errors.New("password cannot be empty")

	// ErrEmptyAdamID is returned when Adam ID is empty.
	ErrEmptyAdamID = errors.New("adamID cannot be empty")

	// ErrInvalidVersionID is returned when version ID is not positive.
	ErrInvalidVersionID = errors.New("versionID must be positive")

	// ErrMissingIdentifiers is returned when neither adamIDs nor bundleIDs are provided.
	ErrMissingIdentifiers = errors.New("either adamIDs or bundleIDs must be provided")
)
