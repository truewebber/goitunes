package service

import "testing"

func TestCurrencyService_ExtractCurrency(t *testing.T) {
	service := NewCurrencyService()

	tests := []struct {
		name           string
		price          float64
		formattedPrice string
		expected       string
	}{
		{"Free app", 0, "Free", ""},
		{"USD with dollar sign", 9.99, "$9.99", "$"},
		{"EUR with euro sign", 8.99, "€8,99", "€"},
		{"GBP with pound sign", 7.99, "£7.99", "£"},
		{"RUB with ruble sign", 599.00, "599,00 ₽", "₽"},
		{"JPY with yen sign", 1200, "¥1,200", "¥"},
		{"Get button", 0, "Get", ""},
		{"Download", 0, "Download", ""},
		{"Complex format", 12.99, "US$ 12.99", "US$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ExtractCurrency(tt.price, tt.formattedPrice)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCurrencyService_FormatPrice(t *testing.T) {
	service := NewCurrencyService()

	tests := []struct {
		name     string
		price    float64
		currency string
		expected string
	}{
		{"Free", 0, "", "Free"},
		{"Free with currency", 0, "$", "Free"},
		{"Paid with currency", 9.99, "$", "$"},
		{"No currency", 9.99, "", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.FormatPrice(tt.price, tt.currency)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
