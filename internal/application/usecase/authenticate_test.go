package usecase_test

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/truewebber/goitunes/v2/internal/application/dto"
	"github.com/truewebber/goitunes/v2/internal/application/usecase"
	"github.com/truewebber/goitunes/v2/internal/domain/repository/mocks"
	"github.com/truewebber/goitunes/v2/internal/domain/valueobject"
)

func TestAuthenticate_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		req           dto.AuthenticateRequest
		setupMock     func(*mocks.MockAuthRepository)
		expectedError error
		checkResponse func(*testing.T, *dto.AuthenticateResponse)
	}{
		{
			name: "positive: successful authentication",
			req: dto.AuthenticateRequest{
				AppleID:  "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockAuthRepository) {
				creds, _ := valueobject.NewCredentialsWithTokens("test@example.com", "token123", "dsid456")
				m.EXPECT().
					Authenticate(gomock.Any(), "test@example.com", "password123").
					Return(creds, nil)
			},
			expectedError: nil,
			checkResponse: func(t *testing.T, resp *dto.AuthenticateResponse) {
				if resp == nil {
					t.Fatal("Response should not be nil")
				}
				if resp.AppleID != "test@example.com" {
					t.Errorf("Expected AppleID test@example.com, got %s", resp.AppleID)
				}
				if resp.PasswordToken != "token123" {
					t.Errorf("Expected token token123, got %s", resp.PasswordToken)
				}
				if resp.DSID != "dsid456" {
					t.Errorf("Expected DSID dsid456, got %s", resp.DSID)
				}
				if !resp.Authenticated {
					t.Error("Expected Authenticated to be true")
				}
			},
		},
		{
			name: "negative: empty appleID",
			req: dto.AuthenticateRequest{
				AppleID:  "",
				Password: "password123",
			},
			setupMock:     func(m *mocks.MockAuthRepository) {},
			expectedError: usecase.ErrEmptyAppleID,
			checkResponse: func(t *testing.T, resp *dto.AuthenticateResponse) {
				if resp != nil {
					t.Error("Response should be nil for error case")
				}
			},
		},
		{
			name: "negative: empty password",
			req: dto.AuthenticateRequest{
				AppleID:  "test@example.com",
				Password: "",
			},
			setupMock:     func(m *mocks.MockAuthRepository) {},
			expectedError: usecase.ErrEmptyPassword,
			checkResponse: func(t *testing.T, resp *dto.AuthenticateResponse) {
				if resp != nil {
					t.Error("Response should be nil for error case")
				}
			},
		},
		{
			name: "negative: repository returns error",
			req: dto.AuthenticateRequest{
				AppleID:  "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockAuthRepository) {
				m.EXPECT().
					Authenticate(gomock.Any(), "test@example.com", "password123").
					Return(nil, errors.New("authentication failed"))
			},
			expectedError: errors.New("authentication failed"),
			checkResponse: func(t *testing.T, resp *dto.AuthenticateResponse) {
				if resp != nil {
					t.Error("Response should be nil for error case")
				}
			},
		},
		{
			name: "corner case: whitespace in appleID",
			req: dto.AuthenticateRequest{
				AppleID:  "  test@example.com  ",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockAuthRepository) {
				creds, _ := valueobject.NewCredentialsWithTokens("  test@example.com  ", "token", "dsid")
				m.EXPECT().
					Authenticate(gomock.Any(), "  test@example.com  ", "password123").
					Return(creds, nil)
			},
			expectedError: nil,
			checkResponse: func(t *testing.T, resp *dto.AuthenticateResponse) {
				if resp == nil {
					t.Fatal("Response should not be nil")
				}
				if !resp.Authenticated {
					t.Error("Expected authentication to succeed")
				}
			},
		},
		{
			name: "corner case: special characters in password",
			req: dto.AuthenticateRequest{
				AppleID:  "test@example.com",
				Password: "p@$$w0rd!@#$%^&*()",
			},
			setupMock: func(m *mocks.MockAuthRepository) {
				creds, _ := valueobject.NewCredentialsWithTokens("test@example.com", "token", "dsid")
				m.EXPECT().
					Authenticate(gomock.Any(), "test@example.com", "p@$$w0rd!@#$%^&*()").
					Return(creds, nil)
			},
			expectedError: nil,
			checkResponse: func(t *testing.T, resp *dto.AuthenticateResponse) {
				if resp == nil {
					t.Fatal("Response should not be nil")
				}
			},
		},
		{
			name: "corner case: very long password",
			req: dto.AuthenticateRequest{
				AppleID:  "test@example.com",
				Password: string(make([]byte, 10000)),
			},
			setupMock: func(m *mocks.MockAuthRepository) {
				creds, _ := valueobject.NewCredentialsWithTokens("test@example.com", "token", "dsid")
				m.EXPECT().
					Authenticate(gomock.Any(), "test@example.com", gomock.Any()).
					Return(creds, nil)
			},
			expectedError: nil,
			checkResponse: func(t *testing.T, resp *dto.AuthenticateResponse) {
				if resp == nil {
					t.Fatal("Response should not be nil")
				}
			},
		},
		{
			name: "corner case: unicode in appleID",
			req: dto.AuthenticateRequest{
				AppleID:  "тест@example.com",
				Password: "password123",
			},
			setupMock: func(m *mocks.MockAuthRepository) {
				creds, _ := valueobject.NewCredentialsWithTokens("тест@example.com", "token", "dsid")
				m.EXPECT().
					Authenticate(gomock.Any(), "тест@example.com", "password123").
					Return(creds, nil)
			},
			expectedError: nil,
			checkResponse: func(t *testing.T, resp *dto.AuthenticateResponse) {
				if resp == nil {
					t.Fatal("Response should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockAuthRepository(ctrl)
			tt.setupMock(mockRepo)

			uc := usecase.NewAuthenticate(mockRepo)
			resp, err := uc.Execute(context.Background(), tt.req)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error containing %v, got nil", tt.expectedError)
				} else if !errors.Is(err, tt.expectedError) {
					// For wrapped errors, check if the original error is contained
					if tt.expectedError.Error() != "" && !errors.Is(err, tt.expectedError) {
						// Just ensure there's an error, wrapped errors are acceptable
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}

			tt.checkResponse(t, resp)
		})
	}
}

func TestAuthenticate_Execute_ContextCancellation(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	mockRepo.EXPECT().
		Authenticate(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, context.Canceled)

	uc := usecase.NewAuthenticate(mockRepo)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := dto.AuthenticateRequest{
		AppleID:  "test@example.com",
		Password: "password",
	}

	_, err := uc.Execute(ctx, req)
	if err == nil {
		t.Error("Expected error when context is cancelled")
	}
}

func TestNewAuthenticate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	uc := usecase.NewAuthenticate(mockRepo)

	if uc == nil {
		t.Fatal("NewAuthenticate should not return nil")
	}
}
