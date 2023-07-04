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
	Automapper Mapper
	ApiClient  *ApiClientBase
}

func (r *CoinMarketProvider) FetchExchangeRate(ctx context.Context, coins []string, currencies []string, precision uint) (model.Rate, error) {
	apiKey := "8f5685ff-4148-40ad-8d88-21d3e5b8d068"

	url := coinMarketCapAPIURL

	normalizedCoins := make([]string, len(coins))
	for i, coin := range coins {
		normalizedCoins[i] = model.NormalizeCurrency(strings.ToLower(coin))
	}

	headers := map[string]string{
		"X-CMC_PRO_API_KEY": apiKey,
	}

	queryParams := map[string]string{
		"symbol":  strings.Join(normalizedCoins, ","),
		"convert": strings.Join(currencies, ","),
	}

	resp, err := r.ApiClient.MakeAPIRequest(ctx, url, queryParams, headers)
	if err != nil {
		return model.Rate{}, err
	}

	var exchangeRateResponse CoinMarkerExchangeRateResponse
	err = json.NewDecoder(resp.Body).Decode(&exchangeRateResponse)
	if err != nil {
		return model.Rate{}, err
	}

	var rate model.Rate
	r.Automapper = &CoinMarkerExchangeRateResponseMapper{}
	rate, err = r.Automapper.MapToRate(exchangeRateResponse)

	return rate, err
}
