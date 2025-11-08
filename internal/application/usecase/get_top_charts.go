package usecase

import (
	"context"
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/application/mapper"
	"github.com/truewebber/goitunes/v2/internal/domain/entity"
	"github.com/truewebber/goitunes/v2/internal/domain/repository"
)

// GetTopCharts retrieves top charts based on the request.
type GetTopCharts struct {
	chartRepo repository.ChartRepository
	mapper    *mapper.ApplicationMapper
}

// NewGetTopCharts creates a new GetTopCharts use case.
func NewGetTopCharts(chartRepo repository.ChartRepository) *GetTopCharts {
	return &GetTopCharts{
		chartRepo: chartRepo,
		mapper:    mapper.NewApplicationMapper(),
	}
}

// Execute retrieves top charts.
func (uc *GetTopCharts) Execute(ctx context.Context, req *dto.GetTopChartsRequest) (*dto.GetTopChartsResponse, error) {
	chartType := uc.parseChartType(req.ChartType)

	var items []*entity.ChartItem

	var err error

	// Determine which endpoint to use based on request
	const top200Limit = 200

	useTop1500 := req.MaxResults > top200Limit || req.Page > 0

	if useTop1500 {
		items, err = uc.getTop1500(ctx, req, chartType)
	} else {
		items, err = uc.getTop200(ctx, req, chartType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get charts: %w", err)
	}

	return &dto.GetTopChartsResponse{
		Items:      uc.mapper.ChartItemsToDTOList(items),
		TotalCount: len(items),
	}, nil
}

// getTop1500 retrieves top 1500 charts.
func (uc *GetTopCharts) getTop1500(
	ctx context.Context,
	req *dto.GetTopChartsRequest,
	chartType entity.ChartType,
) ([]*entity.ChartItem, error) {
	page := req.Page
	if page < 0 {
		page = 0
	}

	pageSize := req.MaxResults
	if pageSize <= 0 {
		pageSize = 100
	}

	items, err := uc.chartRepo.GetTop1500(ctx, req.GenreID, chartType, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get top 1500 charts: %w", err)
	}

	return items, nil
}

// getTop200 retrieves top 200 charts.
func (uc *GetTopCharts) getTop200(
	ctx context.Context,
	req *dto.GetTopChartsRequest,
	chartType entity.ChartType,
) ([]*entity.ChartItem, error) {
	from := req.From
	if from < 1 {
		from = 1
	}

	limit := req.Limit

	const defaultTop200Limit = 200

	if limit <= 0 {
		limit = defaultTop200Limit
	}

	items, err := uc.chartRepo.GetTop200(ctx, req.GenreID, chartType, req.KidPrefix, from, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top 200 charts: %w", err)
	}

	return items, nil
}

// parseChartType converts string to ChartType.
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
