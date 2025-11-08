package valueobject

import "fmt"

// Credentials represents authentication credentials for App Store.
type Credentials struct {
	appleID       string
	passwordToken string // X-Token
	dsid          string // Directory Services ID
	kbsync        string // Certificate for buying
}

// NewCredentials creates a new Credentials value object.
func NewCredentials(appleID string) (*Credentials, error) {
	if appleID == "" {
		return nil, fmt.Errorf("appleID cannot be empty")
	}
	return &Credentials{
		appleID: appleID,
	}, nil
}

// NewCredentialsWithTokens creates credentials with authentication tokens.
func NewCredentialsWithTokens(appleID, passwordToken, dsid string) (*Credentials, error) {
	if appleID == "" {
		return nil, fmt.Errorf("appleID cannot be empty")
	}
	if passwordToken == "" {
		return nil, fmt.Errorf("passwordToken cannot be empty")
	}
	if dsid == "" {
		return nil, fmt.Errorf("dsid cannot be empty")
	}

	return &Credentials{
		appleID:       appleID,
		passwordToken: passwordToken,
		dsid:          dsid,
	}, nil
}

// Getters.
func (c *Credentials) AppleID() string       { return c.appleID }
func (c *Credentials) PasswordToken() string { return c.passwordToken }
func (c *Credentials) DSID() string          { return c.dsid }
func (c *Credentials) Kbsync() string        { return c.kbsync }

// SetPasswordToken sets the password token (after authentication).
func (c *Credentials) SetPasswordToken(token string) *Credentials {
	c.passwordToken = token
	return c
}

// SetDSID sets the DSID (after authentication).
func (c *Credentials) SetDSID(dsid string) *Credentials {
	c.dsid = dsid
	return c
}

// SetKbsync sets the kbsync certificate (for purchases).
func (c *Credentials) SetKbsync(kbsync string) *Credentials {
	c.kbsync = kbsync
	return c
}

// IsAuthenticated returns true if credentials have authentication tokens.
func (c *Credentials) IsAuthenticated() bool {
	return c.passwordToken != "" && c.dsid != ""
}

// CanPurchase returns true if credentials can be used for purchases.
func (c *Credentials) CanPurchase() bool {
	return c.IsAuthenticated() && c.kbsync != ""
}

// Equals checks if two credentials are equal.
func (c *Credentials) Equals(other *Credentials) bool {
	if other == nil {
		return false
	}
	return c.appleID == other.appleID &&
		c.passwordToken == other.passwordToken &&
		c.dsid == other.dsid &&
		c.kbsync == other.kbsync
}
