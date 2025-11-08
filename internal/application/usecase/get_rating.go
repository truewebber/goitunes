package usecase

import (
	"context"
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/domain/repository"
)

// GetRating retrieves rating information for an application.
type GetRating struct {
	appRepo repository.ApplicationRepository
}

// NewGetRating creates a new GetRating use case.
func NewGetRating(appRepo repository.ApplicationRepository) *GetRating {
	return &GetRating{
		appRepo: appRepo,
	}
}

// Execute retrieves rating information.
func (uc *GetRating) Execute(ctx context.Context, req dto.GetRatingRequest) (*dto.GetRatingResponse, error) {
	var rating float64

	var count int

	var err error

	if req.Overall {
		rating, count, err = uc.appRepo.GetOverallRating(ctx, req.AdamID)
	} else {
		rating, count, err = uc.appRepo.GetRating(ctx, req.AdamID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get rating: %w", err)
	}

	return &dto.GetRatingResponse{
		Rating:      rating,
		RatingCount: count,
	}, nil
}
