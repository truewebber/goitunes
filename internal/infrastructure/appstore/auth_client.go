package appstore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/micromdm/plist"

	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/appstore/model"
	"github.com/truewebber/goitunes/v2/internal/infrastructure/config"
	infrahttp "github.com/truewebber/goitunes/v2/internal/infrastructure/http"
)

// AuthClient implements AuthRepository interface.
type AuthClient struct {
	httpClient infrahttp.Client
	store      *valueobject.Store
	device     *valueobject.Device
}

// NewAuthClient creates a new authentication client.
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

// Authenticate performs authentication with Apple ID and password.
func (c *AuthClient) Authenticate(
	ctx context.Context,
	appleID, password string,
) (*valueobject.Credentials, error) {
	if password == "" {
		return nil, ErrEmptyPassword
	}

	authResp, err := c.performAuthRequest(ctx, appleID, password)
	if err != nil {
		return nil, err
	}

	//nolint:gocritic // err is already declared, using = to avoid shadow
	if err = c.validateAuthResponse(authResp); err != nil {
		return nil, err
	}

	credentials, err := valueobject.NewCredentialsWithTokens(appleID, authResp.PasswordToken, authResp.DSID)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %w", err)
	}

	return credentials, nil
}

// performAuthRequest performs the authentication HTTP request.
func (c *AuthClient) performAuthRequest(
	ctx context.Context,
	appleID, password string,
) (*model.AuthResponse, error) {
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

	var authResp model.AuthResponse

	if err = plist.Unmarshal(data, &authResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &authResp, nil
}

// validateAuthResponse validates the authentication response.
func (c *AuthClient) validateAuthResponse(authResp *model.AuthResponse) error {
	if authResp.PasswordToken == "" {
		return ErrPasswordTokenNotFound
	}

	if authResp.DSID == "" {
		return ErrDSIDNotFound
	}

	return nil
}

// buildLoginBody creates the request body for login.
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
