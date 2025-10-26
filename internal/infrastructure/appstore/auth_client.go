package appstore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/micromdm/plist"

	"github.com/truewebber/goitunes/internal/domain/valueobject"
	"github.com/truewebber/goitunes/internal/infrastructure/appstore/model"
	"github.com/truewebber/goitunes/internal/infrastructure/config"
	infrahttp "github.com/truewebber/goitunes/internal/infrastructure/http"
)

// AuthClient implements AuthRepository interface
type AuthClient struct {
	httpClient infrahttp.Client
	store      *valueobject.Store
	device     *valueobject.Device
}

// NewAuthClient creates a new authentication client
func NewAuthClient(
	httpClient infrahttp.Client,
	store *valueobject.Store,
	device *valueobject.Device,
) *AuthClient {
	return &AuthClient{
		httpClient: httpClient,
		store:      store,
		device:     device,
	}
}

// Authenticate performs authentication with Apple ID and password
func (c *AuthClient) Authenticate(
	ctx context.Context,
	appleID, password string,
) (*valueobject.Credentials, error) {
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	loginURL := fmt.Sprintf(config.LoginURLTemplate, c.store.HostPrefix())
	
	body := c.buildLoginBody(appleID, password)
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add(config.HeaderContentType, config.ContentTypeFormEncoded)
	req.Header.Add(config.HeaderUserAgent, c.device.UserAgent())
	req.Header.Add(config.HeaderXAppleStoreFront, c.store.XAppleStoreFront())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed with status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var authResp model.AuthResponse
	if err := plist.Unmarshal(data, &authResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if authResp.PasswordToken == "" {
		return nil, fmt.Errorf("password token not found in response")
	}
	if authResp.DSID == "" {
		return nil, fmt.Errorf("DSID not found in response")
	}

	credentials, err := valueobject.NewCredentialsWithTokens(appleID, authResp.PasswordToken, authResp.DSID)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %w", err)
	}

	return credentials, nil
}

// buildLoginBody creates the request body for login
func (c *AuthClient) buildLoginBody(appleID, password string) *strings.Reader {
	params := url.Values{
		"machineName":   {c.device.MachineName()},
		"why":           {"signin"},
		"attempt":       {"1"},
		"createSession": {"true"},
		"guid":          {c.device.GUID()},
		"appleId":       {appleID},
		"password":      {password},
	}

	return strings.NewReader(params.Encode())
}

