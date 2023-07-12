package external

import (
	"CurrencyRateApp/internal/helper"
	"CurrencyRateApp/internal/model"
	"context"
	"encoding/json"
	"strings"
)

type ExchangeRateOptions struct {
	Coins      []string
	Currencies []string
	Precision  uint
}

type CoinMarketOptions struct {
	ApiKey              string `json:"ApiKey"`
	Host                string `json:"Host"`
	GetExchangeEndpoint string `json:"GetExchangeEndpoint"`
}

type CoinMarkerExchangeRateResponse struct {
	Data map[string]struct {
		Quote map[string]struct {
			Price float64 `json:"price"`
		} `json:"quote"`
	} `json:"data"`
}

type CoinMarketProvider struct {
	apiClient          APIClient
	coinMarketSettings CoinMarketOptions
}

func NewCoinMarketProvider(apiClient APIClient, coinMarketSettings CoinMarketOptions) *CoinMarketProvider {
	return &CoinMarketProvider{
		apiClient:          apiClient,
		coinMarketSettings: coinMarketSettings,
	}
}

func (r *CoinMarketProvider) FetchExchangeRate(ctx context.Context, options ExchangeRateOptions) (model.Rate, error) {
	apiKey := r.coinMarketSettings.ApiKey

	url := r.coinMarketSettings.Host + r.coinMarketSettings.GetExchangeEndpoint

	normalizedCoins := make([]string, len(options.Coins))
	for i, coin := range options.Coins {
		normalizedCoins[i] = helper.NormalizeCurrency(strings.ToLower(coin))
	}

	headers := map[string]string{
		"X-CMC_PRO_API_KEY": apiKey,
	}

	queryParams := map[string]string{
		"symbol":  strings.Join(normalizedCoins, ","),
		"convert": strings.Join(options.Currencies, ","),
	}

	resp, err := r.apiClient.MakeAPIRequest(ctx, url, queryParams, headers)
	if err != nil {
		return model.Rate{}, err
	}

	var exchangeRateResponse CoinMarkerExchangeRateResponse
	err = json.NewDecoder(resp.Body).Decode(&exchangeRateResponse)
	if err != nil {
		return model.Rate{}, err
	}

	var rate model.Rate

	for currency, quote := range exchangeRateResponse.Data {
		for targetCurrency, price := range quote.Quote {
			rate.Rates[currency][targetCurrency] = price.Price
		}
	}

	return rate, err
}
