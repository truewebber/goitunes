package entity_test

import (
	"testing"

	"github.com/truewebber/goitunes/v2/internal/domain/entity"
)

func TestRating(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value float64
		count int
	}{
		{"valid rating", 4.5, 1000},
		{"perfect rating", 5.0, 500},
		{"zero rating", 0.0, 0},
		{"negative rating value", -1.0, 100},
		{"rating above 5", 6.0, 100},
		{"decimal rating", 4.567, 12345},
		{"negative count", 3.5, -1},
		{"zero count", 4.0, 0},
		{"large count", 4.5, 999999999},
		{"minimum values", 0.0, 0},
		{"very precise rating", 4.123456789, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rating := &entity.Rating{
				Value: tt.value,
				Count: tt.count,
			}

			if rating.Value != tt.value {
				t.Errorf("Expected value %f, got %f", tt.value, rating.Value)
			}

			if rating.Count != tt.count {
				t.Errorf("Expected count %d, got %d", tt.count, rating.Count)
			}
		})
	}
}

func TestRating_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("zero value zero count", func(t *testing.T) {
		t.Parallel()

		rating := &entity.Rating{Value: 0, Count: 0}
		if rating.Value != 0 || rating.Count != 0 {
			t.Error("Zero values should be accepted")
		}
	})

	t.Run("high precision value", func(t *testing.T) {
		t.Parallel()

		rating := &entity.Rating{Value: 4.999999999, Count: 1}
		if rating.Value != 4.999999999 {
			t.Error("High precision values should be preserved")
		}
	})

	t.Run("very large count", func(t *testing.T) {
		t.Parallel()

		rating := &entity.Rating{Value: 5.0, Count: 2147483647}
		if rating.Count != 2147483647 {
			t.Error("Large count should be accepted")
		}
	})

	t.Run("negative value", func(t *testing.T) {
		t.Parallel()

		// Rating doesn't validate, so negative values are accepted
		rating := &entity.Rating{Value: -5.0, Count: 100}
		if rating.Value != -5.0 {
			t.Error("Negative value should be stored as-is")
		}
	})

	t.Run("value above maximum", func(t *testing.T) {
		t.Parallel()

		// Rating doesn't validate, so values > 5 are accepted
		rating := &entity.Rating{Value: 10.0, Count: 100}
		if rating.Value != 10.0 {
			t.Error("Value above 5 should be stored as-is")
		}
	})
}

func TestRating_PublicFields(t *testing.T) {
	t.Parallel()

	const (
		testRatingValue = 4.8
		testRatingCount = 250
	)

	// Test that Rating fields are public and can be accessed directly
	rating := &entity.Rating{}

	rating.Value = testRatingValue
	rating.Count = testRatingCount

	if rating.Value != testRatingValue {
		t.Error("Value field should be directly accessible")
	}

	if rating.Count != testRatingCount {
		t.Error("Count field should be directly accessible")
	}
}

func TestRating_Modification(t *testing.T) {
	t.Parallel()

	const (
		newRatingValue = 4.5
		newRatingCount = 200
	)

	rating := &entity.Rating{Value: 4.0, Count: 100}

	// Modify values
	rating.Value = newRatingValue
	rating.Count = newRatingCount

	if rating.Value != newRatingValue {
		t.Error("Value should be modifiable")
	}

	if rating.Count != newRatingCount {
		t.Error("Count should be modifiable")
	}
}

func TestRating_DefaultValues(t *testing.T) {
	t.Parallel()

	// Test zero values
	rating := &entity.Rating{}

	if rating.Value != 0 {
		t.Errorf("Default Value should be 0, got %f", rating.Value)
	}

	if rating.Count != 0 {
		t.Errorf("Default Count should be 0, got %d", rating.Count)
	}
}

