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
