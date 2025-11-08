package appstore

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/micromdm/plist"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
	"github.com/truewebber/goitunes/v2/internal/domain/repository"
	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/appstore/model"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/config"
	infrahttp "github.com/truewebber/goitunes/v2/internal/infrastructure/http"
)

// PurchaseClient implements PurchaseRepository interface.
type PurchaseClient struct {
	httpClient  infrahttp.Client
	store       *valueobject.Store
	credentials *valueobject.Credentials
	device      *valueobject.Device
}

// NewPurchaseClient creates a new purchase client.
func NewPurchaseClient(
	httpClient infrahttp.Client,
	store *valueobject.Store,
	credentials *valueobject.Credentials,
	device *valueobject.Device,
) *PurchaseClient {
	return &PurchaseClient{
		httpClient:  httpClient,
		store:       store,
		credentials: credentials,
		device:      device,
	}
}

// Purchase initiates a purchase for an application.
func (c *PurchaseClient) Purchase(
	ctx context.Context,
	adamID string,
	versionID int64,
) (*entity.DownloadInfo, error) {
	if !c.credentials.CanPurchase() {
		return nil, ErrCredentialsDoNotSupportPurchasing
	}

	purchaseResp, err := c.buyApplication(ctx, adamID, versionID, repository.PricingParameterBuy)
	if err != nil {
		return nil, err
	}

	// Handle re-download case
	if purchaseResp.Metrics.DialogID == "MZCommerceSoftware.OwnsSupersededMinorSoftwareApplicationForUpdate" {
		return nil, ErrApplicationRequiresRedownload
	}

	if len(purchaseResp.SongList) == 0 {
		return nil, fmt.Errorf("%w for adamID: %s", ErrDownloadURLNotFound, adamID)
	}

	if len(purchaseResp.SongList) > 1 {
		return nil, fmt.Errorf("%w for adamID: %s", ErrMultipleDownloadURLs, adamID)
	}

	song := purchaseResp.SongList[0]

	// Confirm download
	confirmErr := c.ConfirmDownload(ctx, song.DownloadID)
	if confirmErr != nil {
		return nil, fmt.Errorf("failed to confirm download: %w", confirmErr)
	}

	if len(song.Sinfs) == 0 {
		return nil, fmt.Errorf("%w for adamID: %s", ErrNoSINFFound, adamID)
	}

	if len(song.Sinfs) > 1 {
		return nil, fmt.Errorf("%w for adamID: %s", ErrMultipleSINFs, adamID)
	}

	sinf := bytes.TrimSpace(song.Sinfs[0].Data)
	if len(sinf) == 0 {
		return nil, fmt.Errorf("%w for adamID: %s", ErrSINFEmpty, adamID)
	}

	// Get bundle ID (may be in different fields)
	bundleID := song.Metadata.BundleID
	if bundleID == "" {
		bundleID = song.Metadata.Q
	}

	if bundleID == "" {
		return nil, fmt.Errorf("%w for adamID: %s", ErrBundleIDNotFound, adamID)
	}

	// Generate iTunes metadata
	metadataBytes := c.generateMetadata(&song, bundleID)

	// Build download info
	downloadInfo := entity.NewDownloadInfo(bundleID, song.URL, song.DownloadKey)
	downloadInfo.SetSinf(base64.StdEncoding.EncodeToString(sinf))
	downloadInfo.SetMetadata(base64.StdEncoding.EncodeToString(metadataBytes))
	downloadInfo.SetDownloadID(song.DownloadID)
	downloadInfo.SetVersionID(song.Metadata.ExternalVersionID)

	// Set download headers
	downloadInfo.AddHeader(config.HeaderUserAgent, valueobject.UserAgentDownload)
	downloadInfo.AddHeader(config.HeaderCookie, "downloadKey="+song.DownloadKey)
	downloadInfo.AddHeader(config.HeaderXAppleStoreFront, c.store.XAppleStoreFront())
	downloadInfo.AddHeader(config.HeaderXDsid, c.credentials.DSID())

	return downloadInfo, nil
}

// ConfirmDownload confirms that a download has been initiated.
func (c *PurchaseClient) ConfirmDownload(ctx context.Context, downloadID string) error {
	query := url.Values{
		"download-id": []string{downloadID},
		"guid":        []string{c.device.GUID()},
	}

	requestURL := fmt.Sprintf(config.ConfirmDownloadTemplate, c.store.HostPrefix()) + "?" + query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add(config.HeaderUserAgent, valueobject.UserAgentDownload)
	req.Header.Add(config.HeaderXAppleStoreFront, c.store.XAppleStoreFront())
	req.Header.Add(config.HeaderXDsid, c.credentials.DSID())
	req.Header.Add(config.HeaderXToken, c.credentials.PasswordToken())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log error but don't fail the function
			_ = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	return nil
}

// buyApplication performs the buy operation.
func (c *PurchaseClient) buyApplication(
	ctx context.Context,
	adamID string,
	versionID int64,
	pricingParameter repository.PricingParameter,
) (*model.PurchaseResponse, error) {
	query := url.Values{
		"xToken": {c.credentials.PasswordToken()},
	}
	requestURL := fmt.Sprintf(config.BuyProductURLTemplate, c.store.HostPrefix()) + "?" + query.Encode()

	body := c.buildBuyBody(adamID, versionID, pricingParameter)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add(config.HeaderContentType, config.ContentTypePlist)
	req.Header.Add(config.HeaderReferer, fmt.Sprintf("http://itunes.apple.com/app/id%s", adamID))
	req.Header.Add(config.HeaderUserAgent, c.device.UserAgent())
	req.Header.Add(config.HeaderXAppleStoreFront, c.store.XAppleStoreFront())
	req.Header.Add(config.HeaderXAppleTz, config.DefaultTimeZone)
	req.Header.Add(config.HeaderXDsid, c.credentials.DSID())
	req.Header.Add(config.HeaderXToken, c.credentials.PasswordToken())

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

	var purchaseResp model.PurchaseResponse

	unmarshalErr := plist.Unmarshal(data, &purchaseResp)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", unmarshalErr)
	}

	return &purchaseResp, nil
}

// buildBuyBody creates the request body for purchase.
func (c *PurchaseClient) buildBuyBody(
	adamID string,
	versionID int64,
	pricingParameter repository.PricingParameter,
) *strings.Reader {
	const nanosecondsPerMillisecond = 1000000

	unixTime := time.Now().UnixNano() / nanosecondsPerMillisecond

	rebuy := "false"
	if pricingParameter == repository.PricingParameterReDownload {
		rebuy = "true"
	}

	template := `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
<dict>
	<key>appExtVrsId</key><string>%d</string>
	<key>guid</key><string>%s</string>
	<key>kbsync</key><data>%s</data>
	<key>machineName</key><string>%s</string>
	<key>mtApp</key><string>com.apple.iTunes</string>
	<key>mtClientId</key><string>3z30dhYIz29Wz4gvz9AEz1NIUDKelm</string>
	<key>mtEventTime</key><string>%d</string>
	<key>mtPageContext</key><string>App Store</string>
	<key>mtPageId</key><string>1140828062</string>
	<key>mtPageType</key><string>Software</string>
	<key>mtPrevPage</key><string>Genre_134583</string>
	<key>mtRequestId</key><string>3z30dhYIz29Wz4gvz9AEz1NIUDKelmzJ4H6DIUSz1HZC</string>
	<key>mtTopic</key><string>xp_its_main</string>
	<key>needDiv</key><string>0</string>
	<key>pg</key><string>default</string>
	<key>price</key><string>0</string>
	<key>pricingParameters</key><string>%s</string>
	<key>rebuy</key><string>%s</string>
	<key>productType</key><string>C</string>
	<key>salableAdamId</key><string>%s</string>
	<key>uuid</key><string>353F3F00-9D87-5BB1-9055-B7761CCD57AA</string>
</dict>
</plist>`

	body := fmt.Sprintf(
		template,
		versionID,
		c.device.GUID(),
		c.credentials.Kbsync(),
		c.device.MachineName(),
		unixTime,
		string(pricingParameter),
		rebuy,
		adamID,
	)

	return strings.NewReader(body)
}

// generateMetadata generates iTunes metadata plist.
func (c *PurchaseClient) generateMetadata(song *model.SongItem, bundleID string) []byte {
	metadata := map[string]interface{}{
		"softwareVersionBundleId":           bundleID,
		"itemId":                            song.Metadata.ItemID,
		"itemName":                          song.Metadata.ItemName,
		"kind":                              "software",
		"playlistName":                      song.Metadata.PlaylistName,
		"artistName":                        song.Metadata.ArtistName,
		"artistId":                          song.Metadata.ArtistID,
		"softwareIcon57x57URL":              song.Metadata.SoftwareIcon57x57URL,
		"bundleShortVersionString":          song.Metadata.BundleShortVersionString,
		"bundleVersion":                     song.Metadata.BundleVersion,
		"genre":                             song.Metadata.Genre,
		"genreId":                           song.Metadata.GenreID,
		"releaseDate":                       song.Metadata.ReleaseDate,
		"copyright":                         song.Metadata.Copyright,
		"softwareVersionExternalIdentifier": song.Metadata.ExternalVersionID,
		"softwareSupportedDeviceIds":        song.Metadata.SoftwareSupportedDeviceIDs,
		"appleId":                           c.credentials.AppleID(),
		"purchaseDate":                      song.PurchaseDate,
		"storeFront":                        fmt.Sprintf("%d", c.store.StoreFront()),
	}

	data, err := plist.Marshal(metadata)
	if err != nil {
		// Return empty slice on error - caller should handle this
		return []byte{}
	}

	return data
}
