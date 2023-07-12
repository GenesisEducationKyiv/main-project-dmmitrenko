package helper

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

func GetAvailableCoins() []string {
	availableCoins := make([]string, 0, len(currencyMapping))
	for coin := range currencyMapping {
		availableCoins = append(availableCoins, coin)
	}
	return availableCoins
}
