package appstore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
	"github.com/truewebber/goitunes/v2/internal/domain/repository"
	"github.com/truewebber/goitunes/v2/internal/domain/service"
	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/appstore/model"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/config"
	infrahttp "github.com/truewebber/goitunes/v2/internal/infrastructure/http"
)

// ChartClient implements ChartRepository interface.
type ChartClient struct {
	httpClient      infrahttp.Client
	store           *valueobject.Store
	appRepo         repository.ApplicationRepository
	currencyService *service.CurrencyService
}

// NewChartClient creates a new chart client.
func NewChartClient(
	httpClient infrahttp.Client,
	store *valueobject.Store,
	appRepo repository.ApplicationRepository,
) *ChartClient {
	return &ChartClient{
		httpClient:      httpClient,
		store:           store,
		appRepo:         appRepo,
		currencyService: service.NewCurrencyService(),
	}
}

// GetTop200 retrieves the top 200 applications.
func (c *ChartClient) GetTop200(
	ctx context.Context,
	genreID string,
	chartType entity.ChartType,
	kidPrefix string,
	from, limit int,
) ([]*entity.ChartItem, error) {
	if from < 1 {
		from = 1
	}

	popID := c.chartTypeToPopID(chartType)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.Top200AppsURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()

	if kidPrefix != "" {
		q.Add("ageBandId", kidPrefix)
	}

	q.Add("genreId", genreID)
	q.Add("popId", popID)
	q.Add("cc", c.store.Region())
	q.Add("l", "en")
	req.URL.RawQuery = q.Encode()

	req.Header.Add(config.HeaderUserAgent, valueobject.UserAgentTop200)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't fail the function
			_ = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.Top200Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	index := response.PageData.SegmentedControl.SelectedIndex
	adamIDs := response.PageData.SegmentedControl.Segments[index].PageData.SelectedChart.AdamIds

	if limit <= 0 {
		limit = response.Properties.DI6TopChartsPageNumIdsPerChart
	}

	fromToChunk := from + limit - 1
	if len(adamIDs) < fromToChunk {
		fromToChunk = len(adamIDs)
	}

	topResults := response.StorePlatformData.Lockup.Results
	infoResults := make(map[string]*entity.Application)

	// Collect missing adamIDs
	var needGetInfo []string

	for i := from - 1; i < fromToChunk; i++ {
		if _, ok := topResults[adamIDs[i]]; !ok {
			needGetInfo = append(needGetInfo, adamIDs[i])
		}
	}

	// Fetch missing application info in batches
	for len(needGetInfo) > 0 {
		batchSize := 50
		if len(needGetInfo) < batchSize {
			batchSize = len(needGetInfo)
		}

		batch := needGetInfo[:batchSize]
		needGetInfo = needGetInfo[batchSize:]

		apps, err := c.appRepo.FindByAdamID(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("failed to get application info: %w", err)
		}

		for _, app := range apps {
			infoResults[app.AdamID()] = app
		}
	}

	// Build chart items
	var chartItems []*entity.ChartItem

	for i := from - 1; i < fromToChunk; i++ {
		position := i + 1
		adamID := adamIDs[i]

		var app *entity.Application
		if appInfo, ok := topResults[adamID]; ok {
			app = c.mapAppItemToEntity(appInfo)
		} else if appInfo, ok := infoResults[adamID]; ok {
			app = appInfo
		} else {
			continue
		}

		chartItem := entity.NewChartItem(app, position, chartType)
		chartItems = append(chartItems, chartItem)
	}

	return chartItems, nil
}

// GetTop1500 retrieves up to 1500 applications.
func (c *ChartClient) GetTop1500(
	ctx context.Context,
	genreID string,
	chartType entity.ChartType,
	page, pageSize int,
) ([]*entity.ChartItem, error) {
	popID := c.chartTypeToPopID(chartType)

	q := url.Values{}
	q.Add("genreId", genreID)
	q.Add("popId", popID)
	q.Add("pageNumbers", fmt.Sprintf("%d", page))
	q.Add("pageSize", fmt.Sprintf("%d", pageSize))
	q.Add("cc", c.store.Region())

	requestURL := fmt.Sprintf("%s?%s", config.TopAppsURL, q.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add(config.HeaderUserAgent, valueobject.UserAgentTop1500)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't fail the function
			_ = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response []model.Top1500Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response) != 1 {
		return nil, ErrUnexpectedResponseStructure
	}

	var chartItems []*entity.ChartItem

	for i, item := range response[0].ContentData {
		position := page*pageSize + i + 1

		rating, err := strconv.ParseFloat(item.UserRating, 64)
		if err != nil {
			rating = 0
		}

		buyParams, err := url.ParseQuery(item.BuyData.ActionParams)
		if err != nil {
			buyParams = url.Values{}
		}

		paramPrice, err := strconv.ParseInt(buyParams.Get("price"), 10, 64)
		if err != nil {
			paramPrice = 0
		}

		currency := c.currencyService.ExtractCurrency(float64(paramPrice), item.ButtonText)
		price := float64(paramPrice) / 1000

		versionID, err := strconv.ParseInt(item.BuyData.VersionID, 10, 64)
		if err != nil {
			versionID = 0
		}

		app := entity.NewApplication(item.ID, item.BuyData.BundleID, "")
		app.SetPrice(price, currency)
		app.SetRating(rating, 0)
		app.SetVersion("", versionID)

		chartItem := entity.NewChartItem(app, position, chartType)
		chartItems = append(chartItems, chartItem)
	}

	return chartItems, nil
}

// chartTypeToPopID converts chart type to popID.
func (c *ChartClient) chartTypeToPopID(chartType entity.ChartType) string {
	switch chartType {
	case entity.ChartTypeTopFree:
		return config.PopIDTopFree
	case entity.ChartTypeTopPaid:
		return config.PopIDTopPaid
	case entity.ChartTypeTopGrossing:
		return config.PopIDTopGrossing
	default:
		return config.PopIDTopFree
	}
}

// mapAppItemToEntity maps API response to entity.
func (c *ChartClient) mapAppItemToEntity(item model.AppItemResponse) *entity.Application {
	app := entity.NewApplication(item.ID, item.BundleID, item.Name)
	app.SetArtistName(item.ArtistName)
	app.SetArtistID(item.ArtistID)
	app.SetRating(item.UserRating.Value, item.UserRating.RatingCount)
	app.SetDeviceFamilies(item.DeviceFamilies)

	if len(item.Offers) > 0 {
		offer := item.Offers[0]
		currency := c.currencyService.ExtractCurrency(offer.Price, offer.PriceFormatted)
		app.SetPrice(offer.Price, currency)
		app.SetVersion(offer.Version.Display, int64(offer.Version.ExternalID))

		if len(offer.Assets) > 0 {
			app.SetFileSize(int64(offer.Assets[0].Size))
		}
	}

	if len(item.Genres) > 0 {
		app.SetGenre(item.Genres[0].GenreID, item.Genres[0].Name)
	}

	if item.Artwork.URL != "" {
		app.SetIconURL(item.Artwork.URL)
	}

	var screenshots []string

	for _, screenList := range item.ScreenshotsByType {
		for _, screen := range screenList {
			screenshots = append(screenshots, screen.URL)
		}
	}

	if len(screenshots) > 0 {
		app.SetScreenshotURLs(screenshots)
	}

	app.SetDescription(item.Description.Standard)
	app.SetMinimumOSVersion(item.MinimumOSVersion)

	return app
}
