package entity_test

import (
	"testing"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
)

func TestNewDownloadInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		bundleID    string
		url         string
		downloadKey string
	}{
		{
			name:        "valid download info",
			bundleID:    "com.example.app",
			url:         "https://example.com/download",
			downloadKey: "key123",
		},
		{
			name:        "empty bundleID",
			bundleID:    "",
			url:         "https://example.com/download",
			downloadKey: "key123",
		},
		{
			name:        "empty URL",
			bundleID:    "com.example.app",
			url:         "",
			downloadKey: "key123",
		},
		{
			name:        "empty downloadKey",
			bundleID:    "com.example.app",
			url:         "https://example.com/download",
			downloadKey: "",
		},
		{
			name:        "all empty",
			bundleID:    "",
			url:         "",
			downloadKey: "",
		},
		{
			name:        "very long URL",
			bundleID:    "com.example.app",
			url:         "https://example.com/" + string(make([]byte, 10000)),
			downloadKey: "key123",
		},
		{
			name:        "special characters in bundleID",
			bundleID:    "com.special.app",
			url:         "https://example.com/download",
			downloadKey: "key123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := entity.NewDownloadInfo(tt.bundleID, tt.url, tt.downloadKey)

			if info == nil {
				t.Fatal("NewDownloadInfo should not return nil")
			}

			if info.BundleID() != tt.bundleID {
				t.Errorf("Expected bundleID %s, got %s", tt.bundleID, info.BundleID())
			}

			if info.URL() != tt.url {
				t.Errorf("Expected URL %s, got %s", tt.url, info.URL())
			}

			if info.DownloadKey() != tt.downloadKey {
				t.Errorf("Expected downloadKey %s, got %s", tt.downloadKey, info.DownloadKey())
			}

			// Check default values
			if info.Sinf() != "" {
				t.Error("Sinf should be empty by default")
			}

			if info.Metadata() != "" {
				t.Error("Metadata should be empty by default")
			}

			if info.Headers() == nil {
				t.Error("Headers should be initialized")
			}

			if len(info.Headers()) != 0 {
				t.Error("Headers should be empty by default")
			}

			if info.DownloadID() != "" {
				t.Error("DownloadID should be empty by default")
			}

			if info.VersionID() != 0 {
				t.Error("VersionID should be zero by default")
			}

			if info.FileSize() != 0 {
				t.Error("FileSize should be zero by default")
			}
		})
	}
}

func TestDownloadInfo_SetSinf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sinf string
	}{
		{"valid sinf", "base64encodedsinf"},
		{"empty sinf", ""},
		{"very long sinf", string(make([]byte, 50000))},
		{"unicode sinf", "test-unicode-data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := entity.NewDownloadInfo("com.test", "url", "key")
			result := info.SetSinf(tt.sinf)

			if result != info {
				t.Error("SetSinf should return the same instance for chaining")
			}

			if info.Sinf() != tt.sinf {
				t.Errorf("Expected sinf %s, got %s", tt.sinf, info.Sinf())
			}
		})
	}
}

func TestDownloadInfo_SetMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		metadata string
	}{
		{"valid metadata", "base64encodedmetadata"},
		{"empty metadata", ""},
		{"very long metadata", string(make([]byte, 100000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := entity.NewDownloadInfo("com.test", "url", "key")
			result := info.SetMetadata(tt.metadata)

			if result != info {
				t.Error("SetMetadata should return the same instance for chaining")
			}

			if info.Metadata() != tt.metadata {
				t.Errorf("Expected metadata %s, got %s", tt.metadata, info.Metadata())
			}
		})
	}
}

func TestDownloadInfo_SetHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		headers map[string]string
	}{
		{
			name:    "valid headers",
			headers: map[string]string{"X-Token": "token123", "X-DSID": "dsid456"},
		},
		{
			name:    "empty headers",
			headers: map[string]string{},
		},
		{
			name:    "nil headers",
			headers: nil,
		},
		{
			name: "many headers",
			headers: map[string]string{
				"Header1": "Value1",
				"Header2": "Value2",
				"Header3": "Value3",
				"Header4": "Value4",
				"Header5": "Value5",
			},
		},
		{
			name:    "headers with special characters",
			headers: map[string]string{"X-Special": "test-unicode-header"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := entity.NewDownloadInfo("com.test", "url", "key")
			result := info.SetHeaders(tt.headers)

			if result != info {
				t.Error("SetHeaders should return the same instance for chaining")
			}

			headers := info.Headers()
			if tt.headers == nil && headers == nil {
				return // Both nil is acceptable
			}

			if tt.headers != nil && headers == nil {
				t.Error("Headers should not be nil when set with non-nil map")

				return
			}

			if len(headers) != len(tt.headers) {
				t.Errorf("Expected %d headers, got %d", len(tt.headers), len(headers))
			}

			for key, value := range tt.headers {
				if headers[key] != value {
					t.Errorf("Expected header %s=%s, got %s", key, value, headers[key])
				}
			}
		})
	}
}

func TestDownloadInfo_AddHeader(t *testing.T) {
	t.Parallel()

	const testTokenValue = "token123"

	t.Run("add single header", func(t *testing.T) {
		t.Parallel()

		info := entity.NewDownloadInfo("com.test", "url", "key")
		result := info.AddHeader("X-Token", testTokenValue)

		if result != info {
			t.Error("AddHeader should return the same instance for chaining")
		}

		if info.Headers()["X-Token"] != testTokenValue {
			t.Error("Header should be added")
		}
	})

	t.Run("add multiple headers", func(t *testing.T) {
		t.Parallel()

		info := entity.NewDownloadInfo("com.test", "url", "key")
		info.AddHeader("X-Token", testTokenValue).
			AddHeader("X-DSID", "dsid456").
			AddHeader("X-Store", "143441")

		headers := info.Headers()

		if headers["X-Token"] != testTokenValue {
			t.Error("Previous header should be preserved")
		}

		if headers["X-DSID"] != "dsid456" {
			t.Error("New header should be added")
		}

		if headers["X-Store"] != "143441" {
			t.Error("Chained header should be added")
		}
	})

	t.Run("overwrite existing header", func(t *testing.T) {
		t.Parallel()

		info := entity.NewDownloadInfo("com.test", "url", "key")
		info.AddHeader("X-Token", "oldtoken")
		info.AddHeader("X-Token", "newtoken")

		if info.Headers()["X-Token"] != "newtoken" {
			t.Error("Header should be overwritten")
		}
	})

	t.Run("add empty key", func(t *testing.T) {
		t.Parallel()

		info := entity.NewDownloadInfo("com.test", "url", "key")
		info.AddHeader("", "value")

		if info.Headers()[""] != "value" {
			t.Error("Empty key should be accepted")
		}
	})

	t.Run("add empty value", func(t *testing.T) {
		t.Parallel()

		info := entity.NewDownloadInfo("com.test", "url", "key")
		info.AddHeader("X-Empty", "")

		if info.Headers()["X-Empty"] != "" {
			t.Error("Empty value should be accepted")
		}
	})
}

func TestDownloadInfo_SetDownloadID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		downloadID string
	}{
		{"valid ID", "download123"},
		{"empty ID", ""},
		{"numeric ID", "123456789"},
		{"UUID ID", "123e4567-e89b-12d3-a456-426614174000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := entity.NewDownloadInfo("com.test", "url", "key")
			result := info.SetDownloadID(tt.downloadID)

			if result != info {
				t.Error("SetDownloadID should return the same instance for chaining")
			}

			if info.DownloadID() != tt.downloadID {
				t.Errorf("Expected downloadID %s, got %s", tt.downloadID, info.DownloadID())
			}
		})
	}
}

func TestDownloadInfo_SetVersionID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		versionID int64
	}{
		{"positive versionID", 123456789},
		{"zero versionID", 0},
		{"negative versionID", -1},
		{"max int64", 9223372036854775807},
		{"min int64", -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := entity.NewDownloadInfo("com.test", "url", "key")
			result := info.SetVersionID(tt.versionID)

			if result != info {
				t.Error("SetVersionID should return the same instance for chaining")
			}

			if info.VersionID() != tt.versionID {
				t.Errorf("Expected versionID %d, got %d", tt.versionID, info.VersionID())
			}
		})
	}
}

func TestDownloadInfo_SetFileSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fileSize int64
	}{
		{"small file", 1024},
		{"medium file", 50 * 1024 * 1024},
		{"large file", 5 * 1024 * 1024 * 1024},
		{"zero size", 0},
		{"negative size", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			info := entity.NewDownloadInfo("com.test", "url", "key")
			result := info.SetFileSize(tt.fileSize)

			if result != info {
				t.Error("SetFileSize should return the same instance for chaining")
			}

			if info.FileSize() != tt.fileSize {
				t.Errorf("Expected fileSize %d, got %d", tt.fileSize, info.FileSize())
			}
		})
	}
}

func TestDownloadInfo_MethodChaining(t *testing.T) {
	t.Parallel()

	info := entity.NewDownloadInfo("com.test", "url", "key")

	// Test method chaining
	result := info.
		SetSinf("sinf123").
		SetMetadata("metadata456").
		SetDownloadID("dl789").
		SetVersionID(100).
		SetFileSize(50000).
		AddHeader("X-Token", "token").
		AddHeader("X-DSID", "dsid")

	if result != info {
		t.Error("All setter methods should return the same instance")
	}

	// Verify all values were set
	if info.Sinf() != "sinf123" {
		t.Error("Sinf not set correctly")
	}

	if info.Metadata() != "metadata456" {
		t.Error("Metadata not set correctly")
	}

	if info.DownloadID() != "dl789" {
		t.Error("DownloadID not set correctly")
	}

	if info.VersionID() != 100 {
		t.Error("VersionID not set correctly")
	}

	if info.FileSize() != 50000 {
		t.Error("FileSize not set correctly")
	}

	if info.Headers()["X-Token"] != "token" {
		t.Error("Header not set correctly")
	}

	if info.Headers()["X-DSID"] != "dsid" {
		t.Error("Header not set correctly")
	}
}

