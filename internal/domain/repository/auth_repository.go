package repository

import (
	"context"

	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
)

// AuthRepository defines the interface for authentication operations.
type AuthRepository interface {
	// Authenticate performs authentication with Apple ID and password
	// Returns updated credentials with password token and DSID
	Authenticate(ctx context.Context, appleID, password string) (*valueobject.Credentials, error)
}

// AuthResponse represents the response from authentication.
type AuthResponse struct {
	PasswordToken   string
	DSID            string
	CreditBalance   string
	FreeSongBalance string
}
