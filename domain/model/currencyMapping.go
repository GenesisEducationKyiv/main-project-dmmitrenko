package model

var currencyMapping = map[string]string{
	"bitcoin":  "BTC",
	"litecoin": "LTC",
	"ethereum": "ETH",
}

func NormalizeCurrency(currency string) string {
	if normalized, ok := currencyMapping[currency]; ok {
		return normalized
	}
	return currency
}
