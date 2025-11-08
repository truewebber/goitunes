package usecase

import (
	"context"
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/application/mapper"
	"github.com/truewebber/goitunes/v2/internal/domain/entity"
	"github.com/truewebber/goitunes/v2/internal/domain/repository"
)

// GetApplicationInfo retrieves application information.
type GetApplicationInfo struct {
	appRepo repository.ApplicationRepository
	mapper  *mapper.ApplicationMapper
}

// NewGetApplicationInfo creates a new GetApplicationInfo use case.
func NewGetApplicationInfo(appRepo repository.ApplicationRepository) *GetApplicationInfo {
	return &GetApplicationInfo{
		appRepo: appRepo,
		mapper:  mapper.NewApplicationMapper(),
	}
}

// Execute retrieves application information.
func (uc *GetApplicationInfo) Execute(
	ctx context.Context,
	req dto.GetApplicationInfoRequest,
) (*dto.GetApplicationInfoResponse, error) {
	var apps []*entity.Application

	var err error

	if len(req.AdamIDs) > 0 {
		apps, err = uc.appRepo.FindByAdamID(ctx, req.AdamIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to find apps by adamID: %w", err)
		}
	} else if len(req.BundleIDs) > 0 {
		apps, err = uc.appRepo.FindByBundleID(ctx, req.BundleIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to find apps by bundleID: %w", err)
		}
	} else {
		return nil, ErrMissingIdentifiers
	}

	return &dto.GetApplicationInfoResponse{
		Applications: uc.mapper.ToDTOList(apps),
	}, nil
}
