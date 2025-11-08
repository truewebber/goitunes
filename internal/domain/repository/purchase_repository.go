package repository

import (
	"context"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
)

// PurchaseRepository defines the interface for purchase and download operations.
type PurchaseRepository interface {
	// Purchase initiates a purchase for an application
	// Returns download information including URL, keys, and metadata
	Purchase(ctx context.Context, adamID string, versionID int64) (*entity.DownloadInfo, error)

	// ConfirmDownload confirms that a download has been initiated
	ConfirmDownload(ctx context.Context, downloadID string) error
}

// PricingParameter defines the type of purchase operation.
type PricingParameter string

const (
	PricingParameterBuy        PricingParameter = "STDQ"
	PricingParameterReDownload PricingParameter = "STDRDL"
)
