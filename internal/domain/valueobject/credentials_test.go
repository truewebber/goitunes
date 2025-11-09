package valueobject_test

import (
	"testing"

	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
)

const (
	testAppleID = "test@example.com"
)

func TestNewCredentials(t *testing.T) {
	t.Parallel()

	appleID := testAppleID

	creds, err := valueobject.NewCredentials(appleID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if creds.AppleID() != appleID {
		t.Errorf("Expected appleID %s, got %s", appleID, creds.AppleID())
	}

	if creds.IsAuthenticated() {
		t.Error("New credentials should not be authenticated")
	}
}

func TestNewCredentials_EmptyAppleID(t *testing.T) {
	t.Parallel()

	_, err := valueobject.NewCredentials("")
	if err == nil {
		t.Error("Expected error for empty appleID")
	}
}

func TestNewCredentialsWithTokens(t *testing.T) {
	t.Parallel()

	appleID := testAppleID
	token := "test_token"
	dsid := "123456"

	creds, err := valueobject.NewCredentialsWithTokens(appleID, token, dsid)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if creds.AppleID() != appleID {
		t.Errorf("Expected appleID %s, got %s", appleID, creds.AppleID())
	}

	if creds.PasswordToken() != token {
		t.Errorf("Expected token %s, got %s", token, creds.PasswordToken())
	}

	if creds.DSID() != dsid {
		t.Errorf("Expected DSID %s, got %s", dsid, creds.DSID())
	}

	if !creds.IsAuthenticated() {
		t.Error("Credentials with tokens should be authenticated")
	}
}

func TestNewCredentialsWithTokens_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		appleID     string
		token       string
		dsid        string
		expectError bool
	}{
		{"Valid", testAppleID, "token", "123", false},
		{"Empty appleID", "", "token", "123", true},
		{"Empty token", testAppleID, "", "123", true},
		{"Empty DSID", testAppleID, "token", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := valueobject.NewCredentialsWithTokens(tt.appleID, tt.token, tt.dsid)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCredentials_CanPurchase(t *testing.T) {
	t.Parallel()

	// Without kbsync
	creds, err := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
	if err != nil {
		t.Fatalf("Failed to create credentials: %v", err)
	}

	if creds.CanPurchase() {
		t.Error("Credentials without kbsync should not be able to purchase")
	}

	// With kbsync
	creds.SetKbsync("test_kbsync")

	if !creds.CanPurchase() {
		t.Error("Credentials with kbsync should be able to purchase")
	}

	// Not authenticated
	creds2, err := valueobject.NewCredentials(testAppleID)
	if err != nil {
		t.Fatalf("Failed to create credentials: %v", err)
	}

	creds2.SetKbsync("test_kbsync")

	if creds2.CanPurchase() {
		t.Error("Unauthenticated credentials should not be able to purchase")
	}
}

func TestCredentials_SetMethods(t *testing.T) {
	t.Parallel()

	t.Run("SetPasswordToken", func(t *testing.T) {
		t.Parallel()

		creds, err := valueobject.NewCredentials(testAppleID)
		if err != nil {
			t.Fatalf("Failed to create credentials: %v", err)
		}

		newToken := "new_token"
		result := creds.SetPasswordToken(newToken)

		if result.PasswordToken() != newToken {
			t.Errorf("Expected token %s, got %s", newToken, result.PasswordToken())
		}

		// Check method chaining
		if result != creds {
			t.Error("SetPasswordToken should return the same instance")
		}
	})

	t.Run("SetDSID", func(t *testing.T) {
		t.Parallel()

		creds, err := valueobject.NewCredentials(testAppleID)
		if err != nil {
			t.Fatalf("Failed to create credentials: %v", err)
		}

		newDSID := "999888"
		result := creds.SetDSID(newDSID)

		if result.DSID() != newDSID {
			t.Errorf("Expected DSID %s, got %s", newDSID, result.DSID())
		}
	})

	t.Run("SetKbsync", func(t *testing.T) {
		t.Parallel()

		creds, err := valueobject.NewCredentials(testAppleID)
		if err != nil {
			t.Fatalf("Failed to create credentials: %v", err)
		}

		newKbsync := "kbsync_data"
		result := creds.SetKbsync(newKbsync)

		if result.Kbsync() != newKbsync {
			t.Errorf("Expected kbsync %s, got %s", newKbsync, result.Kbsync())
		}
	})

	t.Run("Set empty values", func(t *testing.T) {
		t.Parallel()

		creds, err := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
		if err != nil {
			t.Fatalf("Failed to create credentials: %v", err)
		}

		// Should accept empty values (validation happens at creation)
		creds.SetPasswordToken("")
		creds.SetDSID("")
		creds.SetKbsync("")

		if creds.PasswordToken() != "" || creds.DSID() != "" || creds.Kbsync() != "" {
			t.Error("Empty values should be accepted")
		}
	})
}

func TestCredentials_StateTransitions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		setup             func() *valueobject.Credentials
		isAuthenticated   bool
		canPurchase       bool
		description       string
	}{
		{
			name: "unauthenticated",
			setup: func() *valueobject.Credentials {
				creds, _ := valueobject.NewCredentials(testAppleID)
				return creds
			},
			isAuthenticated: false,
			canPurchase:     false,
			description:     "New credentials should be unauthenticated",
		},
		{
			name: "authenticated without kbsync",
			setup: func() *valueobject.Credentials {
				creds, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				return creds
			},
			isAuthenticated: true,
			canPurchase:     false,
			description:     "Authenticated credentials without kbsync cannot purchase",
		},
		{
			name: "authenticated with kbsync",
			setup: func() *valueobject.Credentials {
				creds, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				creds.SetKbsync("kbsync")
				return creds
			},
			isAuthenticated: true,
			canPurchase:     true,
			description:     "Authenticated credentials with kbsync can purchase",
		},
		{
			name: "unauthenticated with kbsync",
			setup: func() *valueobject.Credentials {
				creds, _ := valueobject.NewCredentials(testAppleID)
				creds.SetKbsync("kbsync")
				return creds
			},
			isAuthenticated: false,
			canPurchase:     false,
			description:     "Kbsync alone is not enough to purchase",
		},
		{
			name: "authenticated then reset token",
			setup: func() *valueobject.Credentials {
				creds, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				creds.SetPasswordToken("")
				return creds
			},
			isAuthenticated: false,
			canPurchase:     false,
			description:     "Clearing token should make credentials unauthenticated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			creds := tt.setup()

			if creds.IsAuthenticated() != tt.isAuthenticated {
				t.Errorf("%s: Expected IsAuthenticated=%v, got %v", 
					tt.description, tt.isAuthenticated, creds.IsAuthenticated())
			}

			if creds.CanPurchase() != tt.canPurchase {
				t.Errorf("%s: Expected CanPurchase=%v, got %v", 
					tt.description, tt.canPurchase, creds.CanPurchase())
			}
		})
	}
}

func TestCredentials_Equals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		creds1   func() *valueobject.Credentials
		creds2   func() *valueobject.Credentials
		expected bool
	}{
		{
			name: "identical credentials",
			creds1: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				return c
			},
			creds2: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				return c
			},
			expected: true,
		},
		{
			name: "different appleID",
			creds1: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				return c
			},
			creds2: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens("other@example.com", "token", "123")
				return c
			},
			expected: false,
		},
		{
			name: "different passwordToken",
			creds1: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token1", "123")
				return c
			},
			creds2: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token2", "123")
				return c
			},
			expected: false,
		},
		{
			name: "different dsid",
			creds1: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				return c
			},
			creds2: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "456")
				return c
			},
			expected: false,
		},
		{
			name: "different kbsync",
			creds1: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				c.SetKbsync("kbsync1")
				return c
			},
			creds2: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				c.SetKbsync("kbsync2")
				return c
			},
			expected: false,
		},
		{
			name: "nil comparison",
			creds1: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				return c
			},
			creds2: func() *valueobject.Credentials {
				return nil
			},
			expected: false,
		},
		{
			name: "same instance",
			creds1: func() *valueobject.Credentials {
				c, _ := valueobject.NewCredentialsWithTokens(testAppleID, "token", "123")
				return c
			},
			creds2: func() *valueobject.Credentials {
				return nil // Will be set to creds1 in test
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			creds1 := tt.creds1()
			creds2 := tt.creds2()
			
			// Handle same instance case
			if tt.name == "same instance" {
				creds2 = creds1
			}

			result := creds1.Equals(creds2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
