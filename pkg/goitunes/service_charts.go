package goitunes

import (
	"context"
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/application/usecase"
)

// ChartType represents the type of chart.
type ChartType string

const (
	// ChartTypeTopFree represents free applications chart.
	ChartTypeTopFree ChartType = "topfree"
	// ChartTypeTopPaid represents paid applications chart.
	ChartTypeTopPaid ChartType = "toppaid"
	// ChartTypeTopGrossing represents top grossing applications chart.
	ChartTypeTopGrossing ChartType = "topgrossing"
)

// ChartService provides methods for retrieving app charts.
type ChartService struct {
	useCase *usecase.GetTopCharts
}

const (
	defaultTop200Limit = 200
)

// GetTop200 retrieves the top 200 applications for a genre and chart type.
func (s *ChartService) GetTop200(
	ctx context.Context,
	genre Genre,
	chartType ChartType,
	options ...Top200Option,
) ([]dto.ChartItemDTO, error) {
	req := dto.GetTopChartsRequest{
		GenreID:   genre.String(),
		ChartType: string(chartType),
		From:      1,
		Limit:     defaultTop200Limit,
	}

	// Apply options
	for _, opt := range options {
		opt(&req)
	}

	resp, err := s.useCase.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get top 200 charts: %w", err)
	}

	return resp.Items, nil
}

// GetTop1500 retrieves up to 1500 applications for a genre and chart type.
func (s *ChartService) GetTop1500(
	ctx context.Context,
	genre Genre,
	chartType ChartType,
	page, pageSize int,
) ([]dto.ChartItemDTO, error) {
	req := dto.GetTopChartsRequest{
		GenreID:    genre.String(),
		ChartType:  string(chartType),
		Page:       page,
		MaxResults: pageSize,
	}

	resp, err := s.useCase.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get top 1500 charts: %w", err)
	}

	return resp.Items, nil
}

// Top200Option is a functional option for Top200 requests.
type Top200Option func(*dto.GetTopChartsRequest)

// WithKidPrefix sets the age band filter for charts.
func WithKidPrefix(kidPrefix string) Top200Option {
	return func(req *dto.GetTopChartsRequest) {
		req.KidPrefix = kidPrefix
	}
}

// WithRange sets the range of results to retrieve.
func WithRange(from, limit int) Top200Option {
	return func(req *dto.GetTopChartsRequest) {
		req.From = from
		req.Limit = limit
	}
}
