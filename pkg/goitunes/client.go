package goitunes

import (
	"fmt"

	"github.com/truewebber/goitunes/v2/internal/application/usecase"
	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/appstore"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/config"
	infrahttp "github.com/truewebber/goitunes/v2/internal/infrastructure/http"
)

// Client is the main entry point for the goitunes library
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

// New creates a new goitunes client for the specified region
func New(region string, opts ...Option) (*Client, error) {
	storeRegistry := config.NewStoreRegistry()

	store, err := storeRegistry.GetStore(region)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedRegion, region)
	}

	client := &Client{
		store:         store,
		httpClient:    infrahttp.NewDefaultClient(),
		storeRegistry: storeRegistry,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	// Set default device if not provided
	if client.device == nil {
		client.device, _ = valueobject.NewDevice(
			"00000000-0000-0000-0000-000000000000",
			"DefaultMachine",
			valueobject.UserAgentWindows,
		)
	}

	// Initialize repositories
	client.appRepo = appstore.NewApplicationClient(client.httpClient, client.store)
	client.chartRepo = appstore.NewChartClient(client.httpClient, client.store, client.appRepo)

	if client.credentials != nil {
		client.authRepo = appstore.NewAuthClient(client.httpClient, client.store, client.device)
		client.purchaseRepo = appstore.NewPurchaseClient(
			client.httpClient,
			client.store,
			client.credentials,
			client.device,
		)
	}

	// Initialize services
	client.chartService = &ChartService{
		useCase: usecase.NewGetTopCharts(client.chartRepo),
	}
	client.applicationService = &ApplicationService{
		getInfoUseCase:   usecase.NewGetApplicationInfo(client.appRepo),
		getRatingUseCase: usecase.NewGetRating(client.appRepo),
	}

	if client.authRepo != nil {
		client.authService = &AuthService{
			useCase: usecase.NewAuthenticate(client.authRepo),
			client:  client,
		}
	}

	if client.purchaseRepo != nil {
		client.purchaseService = &PurchaseService{
			useCase: usecase.NewPurchaseApplication(client.purchaseRepo),
		}
	}

	return client, nil
}

// Charts returns the chart service
func (c *Client) Charts() *ChartService {
	return c.chartService
}

// Applications returns the application service
func (c *Client) Applications() *ApplicationService {
	return c.applicationService
}

// Auth returns the authentication service
func (c *Client) Auth() *AuthService {
	if c.authService == nil {
		panic("authentication not available: credentials not provided")
	}
	return c.authService
}

// Purchase returns the purchase service
func (c *Client) Purchase() *PurchaseService {
	if c.purchaseService == nil {
		panic("purchase not available: credentials not provided")
	}
	return c.purchaseService
}

// Region returns the current region
func (c *Client) Region() string {
	return c.store.Region()
}

// SupportedRegions returns all supported regions
func (c *Client) SupportedRegions() []string {
	return c.storeRegistry.GetAllRegions()
}

// IsAuthenticated returns true if the client has valid credentials
func (c *Client) IsAuthenticated() bool {
	return c.credentials != nil && c.credentials.IsAuthenticated()
}

// CanPurchase returns true if the client can make purchases
func (c *Client) CanPurchase() bool {
	return c.credentials != nil && c.credentials.CanPurchase()
}
