package service_test

import (
	"testing"

	"github.com/truewebber/goitunes/v2/internal/domain/service"
)

func TestCurrencyService_ExtractCurrency(t *testing.T) {
	t.Parallel()

	currencyService := service.NewCurrencyService()

	tests := []struct {
		name           string
		formattedPrice string
		expected       string
		price          float64
	}{
		{name: "Free app", price: 0, formattedPrice: "Free", expected: ""},
		{name: "USD with dollar sign", price: 9.99, formattedPrice: "$9.99", expected: "$"},
		{name: "EUR with euro sign", price: 8.99, formattedPrice: "€8,99", expected: "€"},
		{name: "GBP with pound sign", price: 7.99, formattedPrice: "£7.99", expected: "£"},
		{name: "RUB with ruble sign", price: 599.00, formattedPrice: "599,00 ₽", expected: "₽"},
		{name: "JPY with yen sign", price: 1200, formattedPrice: "¥1,200", expected: "¥"},
		{name: "Get button", price: 0, formattedPrice: "Get", expected: ""},
		{name: "Download", price: 0, formattedPrice: "Download", expected: ""},
		{name: "Complex format", price: 12.99, formattedPrice: "US$ 12.99", expected: "US$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := currencyService.ExtractCurrency(tt.price, tt.formattedPrice)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCurrencyService_ExtractCurrency_EdgeCases(t *testing.T) {
	t.Parallel()

	currencyService := service.NewCurrencyService()

	tests := []struct {
		name           string
		formattedPrice string
		expected       string
		price          float64
	}{
		// More unicode currencies
		{name: "Indian Rupee", price: 99, formattedPrice: "₹99", expected: "₹"},
		{name: "Won", price: 1000, formattedPrice: "₩1,000", expected: "₩"},
		{name: "Bitcoin", price: 0.001, formattedPrice: "₿0.001", expected: "₿"},
		{name: "Chinese Yuan", price: 15, formattedPrice: "¥15", expected: "¥"},

		// Edge formatting
		{name: "Price with spaces", price: 9.99, formattedPrice: "$ 9.99", expected: "$"},
		{name: "Currency at end", price: 9.99, formattedPrice: "9.99$", expected: "$"},
		{name: "Multiple symbols", price: 9.99, formattedPrice: "$€¥9.99", expected: "$€¥"},
		{name: "Mixed unicode", price: 100, formattedPrice: "Цена: 100₽", expected: "Цена:₽"},

		// Special cases
		{name: "Empty formatted", price: 9.99, formattedPrice: "", expected: ""},
		{name: "Only numbers", price: 9.99, formattedPrice: "999", expected: ""},
		{name: "Only spaces", price: 9.99, formattedPrice: "   ", expected: ""},
		{name: "Negative price", price: -9.99, formattedPrice: "-$9.99", expected: "-$"},
		{name: "Very large price", price: 999999.99, formattedPrice: "$999,999.99", expected: "$"},

		// Different number formats
		{name: "Dot thousands", price: 1999.99, formattedPrice: "$1.999,99", expected: "$"},
		{name: "Space thousands", price: 1999.99, formattedPrice: "$1 999,99", expected: "$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := currencyService.ExtractCurrency(tt.price, tt.formattedPrice)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCurrencyService_FormatPrice(t *testing.T) {
	t.Parallel()

	currencyService := service.NewCurrencyService()

	tests := []struct {
		name     string
		currency string
		expected string
		price    float64
	}{
		{name: "Free", price: 0, currency: "", expected: "Free"},
		{name: "Free with currency", price: 0, currency: "$", expected: "Free"},
		{name: "Paid with currency", price: 9.99, currency: "$", expected: "$"},
		{name: "No currency", price: 9.99, currency: "", expected: "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := currencyService.FormatPrice(tt.price, tt.currency)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCurrencyService_FormatPrice_EdgeCases(t *testing.T) {
	t.Parallel()

	currencyService := service.NewCurrencyService()

	tests := []struct {
		name     string
		currency string
		expected string
		price    float64
	}{
		// Unicode currencies
		{name: "With euro", price: 9.99, currency: "€", expected: "€"},
		{name: "With ruble", price: 100, currency: "₽", expected: "₽"},
		{name: "With won", price: 1000, currency: "₩", expected: "₩"},
		{name: "With rupee", price: 99, currency: "₹", expected: "₹"},

		// Edge prices
		{name: "Negative price", price: -9.99, currency: "$", expected: "$"},
		{name: "Very small price", price: 0.01, currency: "$", expected: "$"},
		{name: "Very large price", price: 999999.99, currency: "$", expected: "$"},
		{name: "Exactly zero", price: 0.0, currency: "$", expected: "Free"},

		// Special currency strings
		{name: "Long currency code", price: 9.99, currency: "USD", expected: "USD"},
		{name: "Multiple symbols", price: 9.99, currency: "$€¥", expected: "$€¥"},
		{name: "Currency with spaces", price: 9.99, currency: "$ ", expected: "$ "},
		{name: "Unicode text currency", price: 9.99, currency: "долларов", expected: "долларов"},
		{name: "Empty currency with paid", price: 9.99, currency: "", expected: "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := currencyService.FormatPrice(tt.price, tt.currency)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestNewCurrencyService(t *testing.T) {
	t.Parallel()

	svc := service.NewCurrencyService()
	if svc == nil {
		t.Fatal("NewCurrencyService should not return nil")
	}

	// Test that multiple instances work independently
	svc1 := service.NewCurrencyService()
	svc2 := service.NewCurrencyService()

	result1 := svc1.ExtractCurrency(9.99, "$9.99")
	result2 := svc2.ExtractCurrency(9.99, "€9.99")

	if result1 != "$" {
		t.Errorf("svc1: Expected $, got %q", result1)
	}

	if result2 != "€" {
		t.Errorf("svc2: Expected €, got %q", result2)
	}
}
