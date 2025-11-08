package usecase

import (
	"context"
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/application/mapper"
	"github.com/truewebber/goitunes/v2/internal/domain/repository"
)

// PurchaseApplication handles application purchase.
type PurchaseApplication struct {
	purchaseRepo repository.PurchaseRepository
	mapper       *mapper.ApplicationMapper
}

// NewPurchaseApplication creates a new PurchaseApplication use case.
func NewPurchaseApplication(purchaseRepo repository.PurchaseRepository) *PurchaseApplication {
	return &PurchaseApplication{
		purchaseRepo: purchaseRepo,
		mapper:       mapper.NewApplicationMapper(),
	}
}

// Execute performs the purchase.
func (uc *PurchaseApplication) Execute(ctx context.Context, req dto.PurchaseRequest) (*dto.PurchaseResponse, error) {
	if req.AdamID == "" {
		return nil, ErrEmptyAdamID
	}

	if req.VersionID <= 0 {
		return nil, ErrInvalidVersionID
	}

	downloadInfo, err := uc.purchaseRepo.Purchase(ctx, req.AdamID, req.VersionID)
	if err != nil {
		return nil, fmt.Errorf("purchase failed: %w", err)
	}

	return &dto.PurchaseResponse{
		DownloadInfo: uc.mapper.DownloadInfoToDTO(downloadInfo),
	}, nil
}
