package usecase

import (
	"context"
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/domain/repository"
)

// Authenticate performs user authentication.
type Authenticate struct {
	authRepo repository.AuthRepository
}

// NewAuthenticate creates a new Authenticate use case.
func NewAuthenticate(authRepo repository.AuthRepository) *Authenticate {
	return &Authenticate{
		authRepo: authRepo,
	}
}

// Execute performs authentication.
func (uc *Authenticate) Execute(ctx context.Context, req dto.AuthenticateRequest) (*dto.AuthenticateResponse, error) {
	if req.AppleID == "" {
		return nil, ErrEmptyAppleID
	}

	if req.Password == "" {
		return nil, ErrEmptyPassword
	}

	credentials, err := uc.authRepo.Authenticate(ctx, req.AppleID, req.Password)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return &dto.AuthenticateResponse{
		AppleID:       credentials.AppleID(),
		PasswordToken: credentials.PasswordToken(),
		DSID:          credentials.DSID(),
		Authenticated: credentials.IsAuthenticated(),
	}, nil
}
