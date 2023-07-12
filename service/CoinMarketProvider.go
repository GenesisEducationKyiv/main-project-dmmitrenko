package service

import (
	"CurrencyRateApp/domain/model"
	"context"
	"encoding/json"
	"strings"
)

const (
	coinMarketCapAPIURL = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest"
)

type CoinMarkerExchangeRateResponse struct {
	Data map[string]struct {
		Quote map[string]struct {
			Price float64 `json:"price"`
		} `json:"quote"`
	} `json:"data"`
}

type CoinMarketProvider struct {
	automapper Mapper
	apiClient  APIClient
}

func NewCoinMarketProvider(automapper Mapper, apiClient APIClient) *CoinMarketProvider {
	return &CoinMarketProvider{
		automapper: automapper,
		apiClient:  apiClient,
	}
}

func (r *CoinMarketProvider) FetchExchangeRate(ctx context.Context, options ExchangeRateOptions) (model.Rate, error) {
	apiKey := "8f5685ff-4148-40ad-8d88-21d3e5b8d068"

	url := coinMarketCapAPIURL

	normalizedCoins := make([]string, len(options.Coins))
	for i, coin := range options.Coins {
		normalizedCoins[i] = model.NormalizeCurrency(strings.ToLower(coin))
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
	r.automapper = &CoinMarkerExchangeRateResponseMapper{}
	rate, err = r.automapper.MapToRate(exchangeRateResponse)

	return rate, err
}
