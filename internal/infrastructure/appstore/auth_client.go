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

	authResp, err := c.performAuthRequestWithRetry(ctx, appleID, password, 1)
	if err != nil {
		return nil, fmt.Errorf("perform auth request with retry: %w", err)
	}

	if validateErr := c.validateAuthResponse(authResp); validateErr != nil {
		return nil, fmt.Errorf("validate auth response: %w", validateErr)
	}

	credentials, err := valueobject.NewCredentialsWithTokens(appleID, authResp.PasswordToken, authResp.DSID)
	if err != nil {
		return nil, fmt.Errorf("credentials with token: %w", err)
	}

	return credentials, nil
}

const maxRetryAttempts = 4

// performAuthRequestWithRetry performs authentication with retry logic.
func (c *AuthClient) performAuthRequestWithRetry(
	ctx context.Context,
	appleID, password string,
	attempt int,
) (*model.AuthResponse, error) {
	if attempt > maxRetryAttempts {
		return nil, fmt.Errorf("%w: maximum retry attempts (%d) exceeded", ErrAuthenticationFailed, maxRetryAttempts)
	}

	authResp, err := c.performAuthRequest(ctx, appleID, password, attempt)
	if err != nil {
		return nil, fmt.Errorf("perform auth request: %w", err)
	}

	// Check for errors
	if authResp.FailureType != "" {
		// First attempt with invalid credentials - retry once
		if authResp.FailureType == "-5000" && attempt == 1 {
			return c.performAuthRequestWithRetry(ctx, appleID, password, attempt+1)
		}

		// Other failures - return error
		errorMsg := authResp.CustomerMessage
		if errorMsg == "" {
			errorMsg = fmt.Sprintf("authentication failed (failureType: %s)", authResp.FailureType)
		}

		return nil, fmt.Errorf("%w: %s", ErrAuthenticationFailed, errorMsg)
	}

	return authResp, nil
}

// performAuthRequest performs the authentication HTTP request.
func (c *AuthClient) performAuthRequest(
	ctx context.Context,
	appleID, password string,
	attempt int,
) (*model.AuthResponse, error) {
	return c.performAuthRequestWithPod(ctx, appleID, password, attempt, c.store.HostPrefix())
}

// performAuthRequestWithPod performs the authentication HTTP request with specific pod number.
func (c *AuthClient) performAuthRequestWithPod(
	ctx context.Context,
	appleID, password string,
	attempt int,
	pod int,
) (*model.AuthResponse, error) {
	loginURL, err := c.buildLoginURL(pod)
	if err != nil {
		return nil, err
	}

	req, err := c.createAuthRequest(ctx, loginURL, appleID, password, attempt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Handle redirect (302) - Apple may redirect to a different pod based on geolocation/IP
	if resp.StatusCode == http.StatusFound {
		return c.handleRedirectResponse(ctx, resp, appleID, password, attempt)
	}

	if resp.StatusCode != http.StatusOK {
		//nolint:errcheck // Error from Close is not critical here
		_ = resp.Body.Close()

		return nil, fmt.Errorf("%w: unexpected status code %d",
			ErrUnexpectedStatusCode, resp.StatusCode)
	}

	return c.parseAuthResponse(resp)
}

// buildLoginURL builds the login URL with pod number and query parameters.
func (c *AuthClient) buildLoginURL(pod int) (string, error) {
	loginURL := fmt.Sprintf(config.LoginURLTemplate, pod)

	loginURLParsed, err := url.Parse(loginURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse login URL: %w", err)
	}

	c.addPodQueryParams(loginURLParsed, pod)

	return loginURLParsed.String(), nil
}

// handleRedirectResponse handles HTTP 302 redirect response.
func (c *AuthClient) handleRedirectResponse(
	ctx context.Context,
	resp *http.Response,
	appleID, password string,
	attempt int,
) (*model.AuthResponse, error) {
	redirectURL := resp.Header.Get("Location")
	if redirectURL == "" {
		//nolint:errcheck // Error from Close is not critical here
		_ = resp.Body.Close()

		return nil, fmt.Errorf("%w: redirect location not found", ErrUnexpectedStatusCode)
	}

	// Read and close response body
	//nolint:errcheck // Error from Copy and Close are not critical here
	_, _ = io.Copy(io.Discard, resp.Body)
	//nolint:errcheck // Error from Close is not critical here
	_ = resp.Body.Close()

	// Parse redirect URL to extract pod
	redirectURLParsed, parseErr := url.Parse(strings.TrimSpace(redirectURL))
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse redirect URL: %w", parseErr)
	}

	// Extract pod from redirect URL (query params or hostname)
	redirectPod := c.extractPodFromRedirectURL(redirectURLParsed)

	// Retry request with redirect pod and incremented attempt
	return c.performAuthRequestWithPod(ctx, appleID, password, attempt+1, redirectPod)
}

// createAuthRequest creates an authentication HTTP request.
func (c *AuthClient) createAuthRequest(
	ctx context.Context,
	requestURL string,
	appleID, password string,
	attempt int,
) (*http.Request, error) {
	body := c.buildLoginBody(appleID, password, attempt)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	c.setAuthHeaders(req, "")

	return req, nil
}

// setAuthHeaders sets required headers for authentication requests.
func (c *AuthClient) setAuthHeaders(req *http.Request, referer string) {
	req.Header.Set(config.HeaderContentType, config.ContentTypeFormEncoded)
	req.Header.Set(config.HeaderUserAgent, "Configurator/2.17 (Macintosh; OS X 15.2; 24C5089c) AppleWebKit/0620.1.16.11.6")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	if referer != "" {
		req.Header.Set("Referer", referer)
	}
}

// buildLoginBody creates the request body for login.
func (c *AuthClient) buildLoginBody(appleID, password string, attempt int) *strings.Reader {
	params := url.Values{
		"appleId":  {appleID},
		"attempt":  {fmt.Sprintf("%d", attempt)},
		"guid":     {c.device.GUID()},
		"password": {password},
		"rmp":      {"0"},
		"why":      {"signIn"},
	}

	return strings.NewReader(params.Encode())
}

// extractPodFromRedirectURL extracts pod number from redirect URL.
func (c *AuthClient) extractPodFromRedirectURL(u *url.URL) int {
	// Try to extract from query params first (Apple always provides it)
	query := u.Query()

	podStr := query.Get("Pod")
	if podStr != "" {
		var pod int
		if _, err := fmt.Sscanf(podStr, "%d", &pod); err == nil {
			return pod
		}
	}

	// Fallback to extracting from hostname
	return c.extractPodFromHost(u.Host)
}

// extractPodFromHost extracts pod number from hostname (e.g., p71-buy.itunes.apple.com -> 71).
func (c *AuthClient) extractPodFromHost(host string) int {
	if !strings.HasPrefix(host, "p") {
		return c.store.HostPrefix()
	}

	parts := strings.Split(host, "-")
	if len(parts) == 0 || !strings.HasPrefix(parts[0], "p") {
		return c.store.HostPrefix()
	}

	var pod int
	if _, err := fmt.Sscanf(parts[0], "p%d", &pod); err != nil {
		return c.store.HostPrefix()
	}

	return pod
}

// addPodQueryParams adds Pod and PRH query parameters to URL if not present.
func (c *AuthClient) addPodQueryParams(u *url.URL, pod int) {
	query := u.Query()
	if query.Get("Pod") == "" {
		query.Set("Pod", fmt.Sprintf("%d", pod))
	}

	if query.Get("PRH") == "" {
		query.Set("PRH", fmt.Sprintf("%d", pod))
	}

	u.RawQuery = query.Encode()
}

// parseAuthResponse parses the authentication response from HTTP response.
func (c *AuthClient) parseAuthResponse(resp *http.Response) (*model.AuthResponse, error) {
	defer func() {
		//nolint:errcheck // Error from Close in defer is not critical
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d, URL: %s",
			ErrUnexpectedStatusCode, resp.StatusCode, resp.Request.URL.String())
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
