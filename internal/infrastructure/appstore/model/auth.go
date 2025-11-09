package model

// AuthResponse represents the authentication response from Apple.
type AuthResponse struct {
	Action          map[string]interface{} `plist:"action"`
	PasswordToken   string                 `plist:"passwordToken"`
	DSID            string                 `plist:"dsPersonId"`
	CustomerMessage string                 `plist:"customerMessage,omitempty"`
	FailureType     string                 `plist:"failureType,omitempty"`
	CreditBalance   string                 `plist:"creditBalance"`
	FreeSongBalance string                 `plist:"freeSongBalance"`
	JingleDocType   string                 `plist:"jingleDocType"`
	JingleAction    string                 `plist:"jingleAction"`
}
