package entity

// ChartItem represents an application in a chart position.
type ChartItem struct {
	application *Application
	chartType   ChartType
	position    int
}

// ChartType represents the type of chart.
type ChartType string

const (
	ChartTypeTopFree     ChartType = "topfree"
	ChartTypeTopPaid     ChartType = "toppaid"
	ChartTypeTopGrossing ChartType = "topgrossing"
)

// NewChartItem creates a new ChartItem.
func NewChartItem(app *Application, position int, chartType ChartType) *ChartItem {
	return &ChartItem{
		application: app,
		position:    position,
		chartType:   chartType,
	}
}

// Getters.
func (c *ChartItem) Application() *Application { return c.application }
func (c *ChartItem) Position() int             { return c.position }
func (c *ChartItem) ChartType() ChartType      { return c.chartType }
