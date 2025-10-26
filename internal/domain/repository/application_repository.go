package repository

import (
	"context"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
)

// ApplicationRepository defines the interface for application data access
type ApplicationRepository interface {
	// FindByAdamID finds applications by their Adam IDs
	FindByAdamID(ctx context.Context, adamIDs []string) ([]*entity.Application, error)

	// FindByBundleID finds applications by their Bundle IDs
	FindByBundleID(ctx context.Context, bundleIDs []string) ([]*entity.Application, error)

	// GetFullInfo retrieves detailed information about an application
	GetFullInfo(ctx context.Context, adamID string) (*entity.Application, error)

	// GetRating retrieves rating information for an application
	GetRating(ctx context.Context, adamID string) (rating float64, count int, err error)

	// GetOverallRating retrieves overall rating information
	GetOverallRating(ctx context.Context, adamID string) (rating float64, count int, err error)
}
