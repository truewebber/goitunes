package usecase

import (
	"context"
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/application/mapper"
	"github.com/truewebber/goitunes/v2/internal/domain/entity"
	"github.com/truewebber/goitunes/v2/internal/domain/repository"
)

// GetTopCharts retrieves top charts based on the request
type GetTopCharts struct {
	chartRepo repository.ChartRepository
	mapper    *mapper.ApplicationMapper
}

// NewGetTopCharts creates a new GetTopCharts use case
func NewGetTopCharts(chartRepo repository.ChartRepository) *GetTopCharts {
	return &GetTopCharts{
		chartRepo: chartRepo,
		mapper:    mapper.NewApplicationMapper(),
	}
}

// Execute retrieves top charts
func (uc *GetTopCharts) Execute(ctx context.Context, req dto.GetTopChartsRequest) (*dto.GetTopChartsResponse, error) {
	chartType := uc.parseChartType(req.ChartType)

	var items []*entity.ChartItem
	var err error

	// Determine which endpoint to use based on request
	if req.MaxResults > 200 || req.Page > 0 {
		// Use Top1500 endpoint
		page := req.Page
		if page < 0 {
			page = 0
		}
		pageSize := req.MaxResults
		if pageSize <= 0 {
			pageSize = 100
		}

		items, err = uc.chartRepo.GetTop1500(ctx, req.GenreID, chartType, page, pageSize)
	} else {
		// Use Top200 endpoint
		from := req.From
		if from < 1 {
			from = 1
		}
		limit := req.Limit
		if limit <= 0 {
			limit = 200
		}

		items, err = uc.chartRepo.GetTop200(ctx, req.GenreID, chartType, req.KidPrefix, from, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get charts: %w", err)
	}

	return &dto.GetTopChartsResponse{
		Items:      uc.mapper.ChartItemsToDTOList(items),
		TotalCount: len(items),
	}, nil
}

// parseChartType converts string to ChartType
func (uc *GetTopCharts) parseChartType(chartType string) entity.ChartType {
	switch chartType {
	case "topfree":
		return entity.ChartTypeTopFree
	case "toppaid":
		return entity.ChartTypeTopPaid
	case "topgrossing":
		return entity.ChartTypeTopGrossing
	default:
		return entity.ChartTypeTopFree
	}
}
