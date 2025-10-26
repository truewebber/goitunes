package goitunes

import (
	"context"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/application/usecase"
)

// ApplicationService provides methods for retrieving application information
type ApplicationService struct {
	getInfoUseCase   *usecase.GetApplicationInfo
	getRatingUseCase *usecase.GetRating
}

// GetByAdamID retrieves application information by Adam IDs
func (s *ApplicationService) GetByAdamID(ctx context.Context, adamIDs ...string) ([]dto.ApplicationDTO, error) {
	if len(adamIDs) == 0 {
		return nil, ErrInvalidRequest
	}

	req := dto.GetApplicationInfoRequest{
		AdamIDs: adamIDs,
	}

	resp, err := s.getInfoUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Applications, nil
}

// GetByBundleID retrieves application information by Bundle IDs
func (s *ApplicationService) GetByBundleID(ctx context.Context, bundleIDs ...string) ([]dto.ApplicationDTO, error) {
	if len(bundleIDs) == 0 {
		return nil, ErrInvalidRequest
	}

	req := dto.GetApplicationInfoRequest{
		BundleIDs: bundleIDs,
	}

	resp, err := s.getInfoUseCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Applications, nil
}

// GetRating retrieves rating information for an application
func (s *ApplicationService) GetRating(ctx context.Context, adamID string) (*dto.GetRatingResponse, error) {
	if adamID == "" {
		return nil, ErrInvalidRequest
	}

	req := dto.GetRatingRequest{
		AdamID:  adamID,
		Overall: false,
	}

	return s.getRatingUseCase.Execute(ctx, req)
}

// GetOverallRating retrieves overall rating information for an application
func (s *ApplicationService) GetOverallRating(ctx context.Context, adamID string) (*dto.GetRatingResponse, error) {
	if adamID == "" {
		return nil, ErrInvalidRequest
	}

	req := dto.GetRatingRequest{
		AdamID:  adamID,
		Overall: true,
	}

	return s.getRatingUseCase.Execute(ctx, req)
}
