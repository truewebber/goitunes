package goitunes

import (
	"context"

	"github.com/truewebber/goitunes/internal/application/dto"
	"github.com/truewebber/goitunes/internal/application/usecase"
)

// PurchaseService provides purchase and download methods
type PurchaseService struct {
	useCase *usecase.PurchaseApplication
}

// Buy purchases an application and returns download information
func (s *PurchaseService) Buy(
	ctx context.Context,
	adamID string,
	versionID int64,
) (*dto.DownloadInfoDTO, error) {
	req := dto.PurchaseRequest{
		AdamID:    adamID,
		VersionID: versionID,
	}

	resp, err := s.useCase.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	return &resp.DownloadInfo, nil
}

