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
	return c.lookupApplications(ctx, "id", strings.Join(adamIDs, ","))
}

// FindByBundleID finds applications by their Bundle IDs.
func (c *ApplicationClient) FindByBundleID(ctx context.Context, bundleIDs []string) ([]*entity.Application, error) {
	return c.lookupApplications(ctx, "bundleId", strings.Join(bundleIDs, ","))
}

// GetFullInfo retrieves detailed information about an application.
func (c *ApplicationClient) GetFullInfo(ctx context.Context, adamID string) (*entity.Application, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(config.NativeAppInfoURL, adamID), http.NoBody)
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
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.FullAppResponse

	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	item, ok := response.StorePlatformData.ProductDv.Results[adamID]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrAdamIDNotFound, adamID)
	}

	return c.mapToEntity(&item), nil
}

// GetRating retrieves rating information for an application.
func (c *ApplicationClient) GetRating(
	ctx context.Context,
	adamID string,
) (*entity.Rating, error) {
	requestURL := fmt.Sprintf(config.NativeAppRatingInfoURL, adamID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
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
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.RatingResponse

	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &entity.Rating{
		Value: response.UserRating.Value,
		Count: response.UserRating.RatingCount,
	}, nil
}

// GetOverallRating retrieves overall rating information.
func (c *ApplicationClient) GetOverallRating(
	ctx context.Context,
	adamID string,
) (*entity.Rating, error) {
	requestURL := fmt.Sprintf(config.OpenAppOverAllRatingInfoURL, adamID, c.store.Region())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
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
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.OverallRatingResponse

	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, fmt.Errorf("%w for adamID: %s", ErrNoRatingFound, adamID)
	}

	result := response.Results[0]

	return &entity.Rating{
		Value: result.AverageUserRating,
		Count: result.UserRatingCount,
	}, nil
}

// mapToEntity maps API response to domain entity.
func (c *ApplicationClient) mapToEntity(item *model.AppItemResponse) *entity.Application {
	app := entity.NewApplication(item.ID, item.BundleID, item.Name)
	c.setBasicInfo(app, item)
	c.setOffersInfo(app, item)
	c.setGenresInfo(app, item)
	c.setArtworkInfo(app, item)
	c.setScreenshots(app, item)
	c.setDescriptionInfo(app, item)
	c.setReleaseDate(app, item)
	c.setFileSizeFromDevice(app, item)

	return app
}

// setBasicInfo sets basic application information.
func (c *ApplicationClient) setBasicInfo(app *entity.Application, item *model.AppItemResponse) {
	app.SetArtistName(item.ArtistName)
	app.SetArtistID(item.ArtistID)
	app.SetRating(item.UserRating.Value, item.UserRating.RatingCount)
	app.SetDeviceFamilies(item.DeviceFamilies)
}

// setOffersInfo sets offer-related information (price, version, file size).
func (c *ApplicationClient) setOffersInfo(app *entity.Application, item *model.AppItemResponse) {
	if len(item.Offers) == 0 {
		return
	}

	offer := item.Offers[0]
	currency := c.currencyService.ExtractCurrency(offer.Price, offer.PriceFormatted)
	app.SetPrice(offer.Price, currency)
	app.SetVersion(offer.Version.Display, int64(offer.Version.ExternalID))

	if len(offer.Assets) > 0 {
		app.SetFileSize(int64(offer.Assets[0].Size))
	}
}

// setGenresInfo sets genre information.
func (c *ApplicationClient) setGenresInfo(app *entity.Application, item *model.AppItemResponse) {
	if len(item.Genres) > 0 {
		app.SetGenre(item.Genres[0].GenreID, item.Genres[0].Name)
	}
}

// setArtworkInfo sets artwork/icon information.
func (c *ApplicationClient) setArtworkInfo(app *entity.Application, item *model.AppItemResponse) {
	if item.Artwork.URL != "" {
		app.SetIconURL(item.Artwork.URL)
	}
}

// setScreenshots sets screenshot URLs.
func (c *ApplicationClient) setScreenshots(app *entity.Application, item *model.AppItemResponse) {
	var screenshots []string

	for _, screenList := range item.ScreenshotsByType {
		for _, screen := range screenList {
			screenshots = append(screenshots, screen.URL)
		}
	}

	if len(screenshots) > 0 {
		app.SetScreenshotURLs(screenshots)
	}
}

// setDescriptionInfo sets description and minimum OS version.
func (c *ApplicationClient) setDescriptionInfo(app *entity.Application, item *model.AppItemResponse) {
	app.SetDescription(item.Description.Standard)
	app.SetMinimumOSVersion(item.MinimumOSVersion)
}

// setReleaseDate parses and sets release date.
func (c *ApplicationClient) setReleaseDate(app *entity.Application, item *model.AppItemResponse) {
	if item.ReleaseDate == "" {
		return
	}

	releaseDate, err := time.Parse(time.RFC3339, item.ReleaseDate)
	if err == nil {
		app.SetReleaseDate(releaseDate)
	}
}

// setFileSizeFromDevice sets file size from device-specific sizes.
func (c *ApplicationClient) setFileSizeFromDevice(app *entity.Application, item *model.AppItemResponse) {
	if len(item.FileSizeByDevice) == 0 {
		return
	}

	for _, size := range item.FileSizeByDevice {
		app.SetFileSize(int64(size))

		break
	}
}

func (c *ApplicationClient) lookupApplications(
	ctx context.Context,
	queryParamName string,
	queryParamValue string,
) ([]*entity.Application, error) {
	query := url.Values{
		"version":      []string{"2"},
		queryParamName: []string{queryParamValue},
		"p":            []string{"mdm-lockup"},
		"caller":       []string{"MDM"},
		"platform":     []string{"itunes"},
		"cc":           []string{c.store.Region()},
		"l":            []string{"en_us"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.AppInfoURL+"?"+query.Encode(), http.NoBody)
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
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response model.LookupResponse

	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, ErrNoResultsFound
	}

	apps := make([]*entity.Application, 0, len(response.Results))

	for key := range response.Results {
		item := response.Results[key]
		app := c.mapToEntity(&item)
		apps = append(apps, app)
	}

	return apps, nil
}
