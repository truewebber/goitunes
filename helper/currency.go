package helper

import (
	"regexp"
	"strings"
)

func GetCurrency(price float64, priceFormatted string) string {
	var currency string
	if price != 0 {
		rgx := regexp.MustCompile("\\d")
		currency = rgx.ReplaceAllString(priceFormatted, "")
		currency = strings.Replace(currency, ".", "", -1)
		currency = strings.Replace(currency, ",", "", -1)

		currency = strings.TrimSpace(currency)
	}

	return currency
}
