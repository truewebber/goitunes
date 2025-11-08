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

	response, err := c.fetchTop200Response(ctx, genreID, chartType, kidPrefix)
	if err != nil {
		return nil, err
	}

	adamIDs, limit := c.extractAdamIDsAndLimit(response, limit)
	fromToChunk := c.calculateChunkRange(from, limit, len(adamIDs))

	topResults := response.StorePlatformData.Lockup.Results

	infoResults, err := c.fetchMissingAppInfo(ctx, adamIDs, from, fromToChunk, topResults)
	if err != nil {
		return nil, err
	}

	return c.buildChartItems(adamIDs, from, fromToChunk, topResults, infoResults, chartType), nil
}

// GetTop1500 retrieves up to 1500 applications.
func (c *ChartClient) GetTop1500(
	ctx context.Context,
	genreID string,
	chartType entity.ChartType,
	page, pageSize int,
) ([]*entity.ChartItem, error) {
	response, err := c.fetchTop1500Response(ctx, genreID, chartType, page, pageSize)
	if err != nil {
		return nil, err
	}

	if len(response) != 1 {
		return nil, ErrUnexpectedResponseStructure
	}

	return c.buildTop1500ChartItems(&response[0], page, pageSize, chartType), nil
}

// fetchTop200Response fetches the top 200 response from API.
func (c *ChartClient) fetchTop200Response(
	ctx context.Context,
	genreID string,
	chartType entity.ChartType,
	kidPrefix string,
) (*model.Top200Response, error) {
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
	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// fetchTop1500Response fetches the top 1500 response from API.
func (c *ChartClient) fetchTop1500Response(
	ctx context.Context,
	genreID string,
	chartType entity.ChartType,
	page, pageSize int,
) ([]model.Top1500Response, error) {
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
	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

// extractAdamIDsAndLimit extracts Adam IDs and adjusts limit if needed.
func (c *ChartClient) extractAdamIDsAndLimit(
	response *model.Top200Response,
	limit int,
) ([]string, int) {
	index := response.PageData.SegmentedControl.SelectedIndex
	adamIDs := response.PageData.SegmentedControl.Segments[index].PageData.SelectedChart.AdamIDs

	if limit <= 0 {
		limit = response.Properties.DI6TopChartsPageNumIDsPerChart
	}

	return adamIDs, limit
}

// calculateChunkRange calculates the range for processing.
func (c *ChartClient) calculateChunkRange(from, limit, totalLen int) int {
	fromToChunk := from + limit - 1
	if totalLen < fromToChunk {
		fromToChunk = totalLen
	}

	return fromToChunk
}

// fetchMissingAppInfo fetches missing application info in batches.
func (c *ChartClient) fetchMissingAppInfo(
	ctx context.Context,
	adamIDs []string,
	from, fromToChunk int,
	topResults map[string]model.AppItemResponse,
) (map[string]*entity.Application, error) {
	infoResults := make(map[string]*entity.Application)

	var needGetInfo []string

	for i := from - 1; i < fromToChunk; i++ {
		if _, ok := topResults[adamIDs[i]]; !ok {
			needGetInfo = append(needGetInfo, adamIDs[i])
		}
	}

	const batchSize = 50

	for len(needGetInfo) > 0 {
		currentBatchSize := batchSize
		if len(needGetInfo) < batchSize {
			currentBatchSize = len(needGetInfo)
		}

		batch := needGetInfo[:currentBatchSize]
		needGetInfo = needGetInfo[currentBatchSize:]

		apps, err := c.appRepo.FindByAdamID(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("failed to get application info: %w", err)
		}

		for _, app := range apps {
			infoResults[app.AdamID()] = app
		}
	}

	return infoResults, nil
}

// buildChartItems builds chart items from available data.
func (c *ChartClient) buildChartItems(
	adamIDs []string,
	from, fromToChunk int,
	topResults map[string]model.AppItemResponse,
	infoResults map[string]*entity.Application,
	chartType entity.ChartType,
) []*entity.ChartItem {
	chartItems := make([]*entity.ChartItem, 0, fromToChunk-from+1)

	for i := from - 1; i < fromToChunk; i++ {
		position := i + 1
		adamID := adamIDs[i]

		app := c.getApplicationForChartItem(adamID, topResults, infoResults)
		if app == nil {
			continue
		}

		chartItem := entity.NewChartItem(app, position, chartType)
		chartItems = append(chartItems, chartItem)
	}

	return chartItems
}

// getApplicationForChartItem gets application for chart item from available sources.
func (c *ChartClient) getApplicationForChartItem(
	adamID string,
	topResults map[string]model.AppItemResponse,
	infoResults map[string]*entity.Application,
) *entity.Application {
	if appInfo, ok := topResults[adamID]; ok {
		return c.mapAppItemToEntity(&appInfo)
	}

	if appInfo, found := infoResults[adamID]; found {
		return appInfo
	}

	return nil
}

// buildTop1500ChartItems builds chart items from top 1500 response data.
func (c *ChartClient) buildTop1500ChartItems(
	response *model.Top1500Response,
	page, pageSize int,
	chartType entity.ChartType,
) []*entity.ChartItem {
	chartItems := make([]*entity.ChartItem, 0, len(response.ContentData))

	for i := range response.ContentData {
		item := &response.ContentData[i]
		position := page*pageSize + i + 1

		app := c.buildAppFromTop1500Item(item)
		chartItem := entity.NewChartItem(app, position, chartType)
		chartItems = append(chartItems, chartItem)
	}

	return chartItems
}

// buildAppFromTop1500Item builds application entity from top 1500 item.
func (c *ChartClient) buildAppFromTop1500Item(item *struct {
	ID         string `json:"id"`
	UserRating string `json:"userRating"`
	ButtonText string `json:"buttonText"`
	BuyData    struct {
		BundleID     string `json:"bundleId"`
		VersionID    string `json:"versionId"`
		ActionParams string `json:"actionParams"`
	} `json:"buyData"`
}) *entity.Application {
	rating := c.parseRating(item.UserRating)
	price, currency := c.parsePrice(item.BuyData.ActionParams, item.ButtonText)
	versionID := c.parseVersionID(item.BuyData.VersionID)

	app := entity.NewApplication(item.ID, item.BuyData.BundleID, "")
	app.SetPrice(price, currency)
	app.SetRating(rating, 0)
	app.SetVersion("", versionID)

	return app
}

// parseRating parses rating from string.
func (c *ChartClient) parseRating(ratingStr string) float64 {
	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		return 0
	}

	return rating
}

// parsePrice parses price and currency from buy params and button text.
func (c *ChartClient) parsePrice(actionParams, buttonText string) (float64, string) {
	buyParams, err := url.ParseQuery(actionParams)
	if err != nil {
		return 0, ""
	}

	paramPrice, err := strconv.ParseInt(buyParams.Get("price"), 10, 64)
	if err != nil {
		return 0, ""
	}

	currency := c.currencyService.ExtractCurrency(float64(paramPrice), buttonText)

	const priceDivisor = 1000

	price := float64(paramPrice) / priceDivisor

	return price, currency
}

// parseVersionID parses version ID from string.
func (c *ChartClient) parseVersionID(versionIDStr string) int64 {
	versionID, err := strconv.ParseInt(versionIDStr, 10, 64)
	if err != nil {
		return 0
	}

	return versionID
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
func (c *ChartClient) mapAppItemToEntity(item *model.AppItemResponse) *entity.Application {
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
