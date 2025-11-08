package valueobject

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
		return nil, ErrEmptyAppleID
	}

	return &Credentials{
		appleID: appleID,
	}, nil
}

// NewCredentialsWithTokens creates credentials with authentication tokens.
func NewCredentialsWithTokens(appleID, passwordToken, dsid string) (*Credentials, error) {
	if appleID == "" {
		return nil, ErrEmptyAppleID
	}

	if passwordToken == "" {
		return nil, ErrEmptyPasswordToken
	}

	if dsid == "" {
		return nil, ErrEmptyDSID
	}

	return &Credentials{
		appleID:       appleID,
		passwordToken: passwordToken,
		dsid:          dsid,
	}, nil
}

// AppleID returns the Apple ID.
func (c *Credentials) AppleID() string { return c.appleID }

// PasswordToken returns the password token.
func (c *Credentials) PasswordToken() string { return c.passwordToken }

// DSID returns the DSID.
func (c *Credentials) DSID() string { return c.dsid }

// Kbsync returns the kbsync certificate.
func (c *Credentials) Kbsync() string { return c.kbsync }

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
