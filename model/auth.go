package model

type (
	AuthResponse struct {
		PasswordToken   string `plist:"passwordToken"`
		DSID            string `plist:"dsPersonId"`
		CreditBalance   string `plist:"creditBalance"`
		FreeSongBalance string `plist:"freeSongBalance"`
	}
)
