package goitunes

import (
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/usecase"
	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/appstore"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/config"
	infrahttp "github.com/truewebber/goitunes/v2/internal/infrastructure/http"
)

// Client is the main entry point for the goitunes library.
type Client struct {
	store         *valueobject.Store
	httpClient    infrahttp.Client
	credentials   *valueobject.Credentials
	device        *valueobject.Device
	storeRegistry *config.StoreRegistry

	// Repository implementations
	appRepo      *appstore.ApplicationClient
	chartRepo    *appstore.ChartClient
	authRepo     *appstore.AuthClient
	purchaseRepo *appstore.PurchaseClient

	// Services
	chartService       *ChartService
	applicationService *ApplicationService
	authService        *AuthService
	purchaseService    *PurchaseService
}

// New creates a new goitunes client for the specified region.
func New(region string, opts ...Option) (*Client, error) {
	storeRegistry := config.NewStoreRegistry()

	store, storeErr := storeRegistry.GetStore(region)
	if storeErr != nil {
		return nil, fmt.Errorf("get store: %w", storeErr)
	}

	client := &Client{
		store:         store,
		httpClient:    infrahttp.NewDefaultClient(),
		storeRegistry: storeRegistry,
	}

	if err := client.applyOptions(opts); err != nil {
		return nil, fmt.Errorf("apply options: %w", err)
	}

	client.setDefaultDevice()
	client.initializeRepositories()
	client.initializeServices()

	return client, nil
}

// Charts returns the chart service.
func (c *Client) Charts() *ChartService {
	return c.chartService
}

// Applications returns the application service.
func (c *Client) Applications() *ApplicationService {
	return c.applicationService
}

// Auth returns the authentication service.
func (c *Client) Auth() *AuthService {
	if c.authService == nil {
		panic("authentication not available: credentials not provided")
	}

	return c.authService
}

// Purchase returns the purchase service.
func (c *Client) Purchase() *PurchaseService {
	if c.purchaseService == nil {
		panic("purchase not available: credentials not provided")
	}

	return c.purchaseService
}

// Region returns the current region.
func (c *Client) Region() string {
	return c.store.Region()
}

// SupportedRegions returns all supported regions.
func (c *Client) SupportedRegions() []string {
	return c.storeRegistry.GetAllRegions()
}

// IsAuthenticated returns true if the client has valid credentials.
func (c *Client) IsAuthenticated() bool {
	return c.credentials != nil && c.credentials.IsAuthenticated()
}

// CanPurchase returns true if the client can make purchases.
func (c *Client) CanPurchase() bool {
	return c.credentials != nil && c.credentials.CanPurchase()
}

// applyOptions applies options to the client.
func (c *Client) applyOptions(opts []Option) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}

	return nil
}

// setDefaultDevice sets default device if not provided.
func (c *Client) setDefaultDevice() {
	if c.device != nil {
		return
	}

	device, err := valueobject.NewDevice(
		"00000000-0000-0000-0000-000000000000",
		"DefaultMachine",
		valueobject.UserAgentWindows,
	)
	if err != nil {
		// This should never happen with valid defaults
		panic(fmt.Errorf("failed to create default device: %w", err))
	}

	c.device = device
}

// initializeRepositories initializes repository implementations.
func (c *Client) initializeRepositories() {
	c.appRepo = appstore.NewApplicationClient(c.httpClient, c.store)
	c.chartRepo = appstore.NewChartClient(c.httpClient, c.store, c.appRepo)

	if c.credentials != nil {
		c.authRepo = appstore.NewAuthClient(c.httpClient, c.store, c.device)
		c.purchaseRepo = appstore.NewPurchaseClient(
			c.httpClient,
			c.store,
			c.credentials,
			c.device,
		)
	}
}

// initializeServices initializes service implementations.
func (c *Client) initializeServices() {
	c.chartService = &ChartService{
		useCase: usecase.NewGetTopCharts(c.chartRepo),
	}
	c.applicationService = &ApplicationService{
		getInfoUseCase:   usecase.NewGetApplicationInfo(c.appRepo),
		getRatingUseCase: usecase.NewGetRating(c.appRepo),
	}

	if c.authRepo != nil {
		c.authService = &AuthService{
			useCase: usecase.NewAuthenticate(c.authRepo),
			client:  c,
		}
	}

	if c.purchaseRepo != nil {
		c.purchaseService = &PurchaseService{
			useCase: usecase.NewPurchaseApplication(c.purchaseRepo),
		}
	}
}
