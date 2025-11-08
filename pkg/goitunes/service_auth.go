package goitunes

import (
	"context"
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/application/usecase"
	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/appstore"
)

// AuthService provides authentication methods.
type AuthService struct {
	useCase *usecase.Authenticate
	client  *Client
}

// Login performs authentication with Apple ID and password.
// After successful login, the client will be authenticated and can use purchase methods.
func (s *AuthService) Login(ctx context.Context, password string) (*dto.AuthenticateResponse, error) {
	if s.client.credentials == nil {
		return nil, ErrInvalidCredentials
	}

	req := dto.AuthenticateRequest{
		AppleID:  s.client.credentials.AppleID(),
		Password: password,
	}

	resp, err := s.useCase.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	// Update client credentials with tokens
	credentials, err := valueobject.NewCredentialsWithTokens(
		resp.AppleID,
		resp.PasswordToken,
		resp.DSID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials with tokens: %w", err)
	}

	// Copy kbsync if it exists
	if s.client.credentials.Kbsync() != "" {
		credentials.SetKbsync(s.client.credentials.Kbsync())
	}

	s.client.credentials = credentials

	// Reinitialize auth and purchase repositories with new credentials
	s.client.authRepo = appstore.NewAuthClient(
		s.client.httpClient,
		s.client.store,
		s.client.device,
	)
	s.client.purchaseRepo = appstore.NewPurchaseClient(
		s.client.httpClient,
		s.client.store,
		s.client.credentials,
		s.client.device,
	)

	// Reinitialize purchase service
	if s.client.purchaseRepo != nil {
		s.client.purchaseService = &PurchaseService{
			useCase: usecase.NewPurchaseApplication(s.client.purchaseRepo),
		}
	}

	return resp, nil
}

// IsAuthenticated returns true if the client is authenticated.
func (s *AuthService) IsAuthenticated() bool {
	return s.client.IsAuthenticated()
}
