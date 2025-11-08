package appstore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
	"github.com/truewebber/goitunes/v2/internal/domain/service"
	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/appstore/model"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/config"
	infrahttp "github.com/truewebber/goitunes/v2/internal/infrastructure/http"
)

// ApplicationClient implements ApplicationRepository interface.
type ApplicationClient struct {
	httpClient      infrahttp.Client
	store           *valueobject.Store
	currencyService *service.CurrencyService
}

// NewApplicationClient creates a new application client.
func NewApplicationClient(httpClient infrahttp.Client, store *valueobject.Store) *ApplicationClient {
	return &ApplicationClient{
		httpClient:      httpClient,
		store:           store,
		currencyService: service.NewCurrencyService(),
	}
}

// FindByAdamID finds applications by their Adam IDs.
func (c *ApplicationClient) FindByAdamID(ctx context.Context, adamIDs []string) ([]*entity.Application, error) {
	query := url.Values{
		"version":  []string{"2"},
		"id":       []string{strings.Join(adamIDs, ",")},
		"p":        []string{"mdm-lockup"},
		"caller":   []string{"MDM"},
		"platform": []string{"itunes"},
		"cc":       []string{c.store.Region()},
		"l":        []string{"en_us"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.AppInfoURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.LookupResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no results found for adamIds: %s", strings.Join(adamIDs, ","))
	}

	var apps []*entity.Application
	for _, item := range response.Results {
		app := c.mapToEntity(item)
		apps = append(apps, app)
	}

	return apps, nil
}

// FindByBundleID finds applications by their Bundle IDs.
func (c *ApplicationClient) FindByBundleID(ctx context.Context, bundleIDs []string) ([]*entity.Application, error) {
	query := url.Values{
		"version":  []string{"2"},
		"bundleId": []string{strings.Join(bundleIDs, ",")},
		"p":        []string{"mdm-lockup"},
		"caller":   []string{"MDM"},
		"platform": []string{"itunes"},
		"cc":       []string{c.store.Region()},
		"l":        []string{"en_us"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.AppInfoURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.LookupResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no results found for bundleIds: %s", strings.Join(bundleIDs, ","))
	}

	var apps []*entity.Application
	for _, item := range response.Results {
		app := c.mapToEntity(item)
		apps = append(apps, app)
	}

	return apps, nil
}

// GetFullInfo retrieves detailed information about an application.
func (c *ApplicationClient) GetFullInfo(ctx context.Context, adamID string) (*entity.Application, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(config.NativeAppInfoURL, adamID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add(config.HeaderXAppleStoreFront, c.store.XAppleStoreFrontWithDevice(config.IPhoneDeviceCode))

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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.FullAppResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	item, ok := response.StorePlatformData.ProductDv.Results[adamID]
	if !ok {
		return nil, fmt.Errorf("adamID %s not found in response", adamID)
	}

	return c.mapToEntity(item), nil
}

// GetRating retrieves rating information for an application.
func (c *ApplicationClient) GetRating(ctx context.Context, adamID string) (float64, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(config.NativeAppRatingInfoURL, adamID), nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add(config.HeaderXAppleStoreFront, c.store.XAppleStoreFrontWithDevice(config.IPhoneDeviceCode))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't fail the function
			_ = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.RatingResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.UserRating.Value, response.UserRating.RatingCount, nil
}

// GetOverallRating retrieves overall rating information.
func (c *ApplicationClient) GetOverallRating(ctx context.Context, adamID string) (float64, int, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(config.OpenAppOverAllRatingInfoURL, adamID, c.store.Region()),
		nil,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't fail the function
			_ = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.OverallRatingResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Results) == 0 {
		return 0, 0, fmt.Errorf("no rating found for adamID: %s", adamID)
	}

	result := response.Results[0]
	return result.AverageUserRating, result.UserRatingCount, nil
}

// mapToEntity maps API response to domain entity.
func (c *ApplicationClient) mapToEntity(item model.AppItemResponse) *entity.Application {
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

	// Parse release date
	if item.ReleaseDate != "" {
		if releaseDate, err := time.Parse(time.RFC3339, item.ReleaseDate); err == nil {
			app.SetReleaseDate(releaseDate)
		}
	}

	// Get file size from device-specific sizes
	if len(item.FileSizeByDevice) > 0 {
		for _, size := range item.FileSizeByDevice {
			app.SetFileSize(int64(size))
			break
		}
	}

	return app
}
