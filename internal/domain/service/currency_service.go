package service

import (
	"regexp"
	"strings"
)

// CurrencyService provides currency-related operations.
type CurrencyService struct{}

// NewCurrencyService creates a new CurrencyService.
func NewCurrencyService() *CurrencyService {
	return &CurrencyService{}
}

// ExtractCurrency extracts the currency label from a formatted price string.
// Examples: "$1.99" -> "$", "Free" -> "", "Get" -> "".
func (s *CurrencyService) ExtractCurrency(price float64, formattedPrice string) string {
	if price == 0 {
		return ""
	}

	// Remove all digits, spaces, commas, and dots
	re := regexp.MustCompile(`[0-9\s,.]`)
	currency := re.ReplaceAllString(formattedPrice, "")
	currency = strings.TrimSpace(currency)

	return currency
}

// FormatPrice formats a price with currency.
func (s *CurrencyService) FormatPrice(price float64, currency string) string {
	if price == 0 {
		return "Free"
	}

	if currency == "" {
		return "0"
	}

	return currency
}
