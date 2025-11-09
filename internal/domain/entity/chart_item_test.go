package entity_test

import (
	"testing"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
)

func TestNewChartItem(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		app       *entity.Application
		position  int
		chartType entity.ChartType
	}{
		{
			name:      "valid chart item - topfree",
			app:       entity.NewApplication("123", "com.test.app", "Test App"),
			position:  1,
			chartType: entity.ChartTypeTopFree,
		},
		{
			name:      "valid chart item - toppaid",
			app:       entity.NewApplication("456", "com.test.paid", "Paid App"),
			position:  10,
			chartType: entity.ChartTypeTopPaid,
		},
		{
			name:      "valid chart item - topgrossing",
			app:       entity.NewApplication("789", "com.test.gross", "Grossing App"),
			position:  100,
			chartType: entity.ChartTypeTopGrossing,
		},
		{
			name:      "position zero",
			app:       entity.NewApplication("123", "com.test.app", "Test App"),
			position:  0,
			chartType: entity.ChartTypeTopFree,
		},
		{
			name:      "negative position",
			app:       entity.NewApplication("123", "com.test.app", "Test App"),
			position:  -1,
			chartType: entity.ChartTypeTopFree,
		},
		{
			name:      "large position",
			app:       entity.NewApplication("123", "com.test.app", "Test App"),
			position:  999999,
			chartType: entity.ChartTypeTopFree,
		},
		{
			name:      "nil application",
			app:       nil,
			position:  1,
			chartType: entity.ChartTypeTopFree,
		},
		{
			name:      "empty chart type",
			app:       entity.NewApplication("123", "com.test.app", "Test App"),
			position:  1,
			chartType: "",
		},
		{
			name:      "invalid chart type",
			app:       entity.NewApplication("123", "com.test.app", "Test App"),
			position:  1,
			chartType: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			item := entity.NewChartItem(tt.app, tt.position, tt.chartType)

			if item == nil {
				t.Fatal("NewChartItem should not return nil")
			}

			if item.Application() != tt.app {
				t.Error("Application mismatch")
			}

			if item.Position() != tt.position {
				t.Errorf("Expected position %d, got %d", tt.position, item.Position())
			}

			if item.ChartType() != tt.chartType {
				t.Errorf("Expected chartType %s, got %s", tt.chartType, item.ChartType())
			}
		})
	}
}

func TestChartItem_Accessors(t *testing.T) {
	t.Parallel()

	app := entity.NewApplication("123", "com.test.app", "Test App")
	app.SetPrice(9.99, "$")
	app.SetRating(4.5, 1000)

	item := entity.NewChartItem(app, 42, entity.ChartTypeTopPaid)

	t.Run("Application", func(t *testing.T) {
		t.Parallel()

		if item.Application() != app {
			t.Error("Application should match")
		}

		if item.Application().AdamID() != "123" {
			t.Error("Should be able to access application properties")
		}
	})

	t.Run("Position", func(t *testing.T) {
		t.Parallel()

		if item.Position() != 42 {
			t.Errorf("Expected position 42, got %d", item.Position())
		}
	})

	t.Run("ChartType", func(t *testing.T) {
		t.Parallel()

		if item.ChartType() != entity.ChartTypeTopPaid {
			t.Errorf("Expected ChartTypeTopPaid, got %s", item.ChartType())
		}
	})
}

func TestChartType_Constants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		chartType entity.ChartType
		expected  string
	}{
		{"ChartTypeTopFree", entity.ChartTypeTopFree, "topfree"},
		{"ChartTypeTopPaid", entity.ChartTypeTopPaid, "toppaid"},
		{"ChartTypeTopGrossing", entity.ChartTypeTopGrossing, "topgrossing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if string(tt.chartType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.chartType))
			}
		})
	}
}

func TestChartItem_WithNilApplication(t *testing.T) {
	t.Parallel()

	item := entity.NewChartItem(nil, 1, entity.ChartTypeTopFree)

	if item == nil {
		t.Fatal("NewChartItem should not return nil even with nil app")
	}

	if item.Application() != nil {
		t.Error("Application should be nil")
	}

	// Should not panic
	_ = item.Position()
	_ = item.ChartType()
}

