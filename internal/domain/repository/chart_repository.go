package repository

import (
	"context"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
)

// ChartRepository defines the interface for chart data access
type ChartRepository interface {
	// GetTop200 retrieves the top 200 applications for a genre and chart type
	// genreID: the genre identifier (e.g., "36" for all)
	// chartType: the type of chart (topfree, toppaid, topgrossing)
	// kidPrefix: optional age band filter
	// from: starting position (1-based)
	// limit: number of results to return
	GetTop200(ctx context.Context, genreID string, chartType entity.ChartType, kidPrefix string, from, limit int) ([]*entity.ChartItem, error)

	// GetTop1500 retrieves up to 1500 applications for a genre and chart type
	// genreID: the genre identifier
	// chartType: the type of chart
	// page: page number (0-based)
	// pageSize: number of items per page
	GetTop1500(ctx context.Context, genreID string, chartType entity.ChartType, page, pageSize int) ([]*entity.ChartItem, error)
}
