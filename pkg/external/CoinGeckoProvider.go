package external

import (
	"CurrencyRateApp/internal/model"

	"context"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

const (
	coinParameters     = "ids"
	currencyParameters = "vs_currencies"
	currencyPrecision  = "precision"
)

type CoingeckoExchangeRateResponse struct {
	Rates map[string]map[string]float64 `json:"rates"`
}

type CoingeckoOptions struct {
	Host                string `json:"Host"`
	GetExchangeEndpoint string `json:"GetExchangeEndpoint"`
}

type CoingeckoProvider struct {
	apiClient         APIClient
	coingeckoSettings CoingeckoOptions
}

func NewCoingeckoProvider(apiClient APIClient, coingeckoSettings CoingeckoOptions) *CoingeckoProvider {
	return &CoingeckoProvider{
		apiClient:         apiClient,
		coingeckoSettings: coingeckoSettings,
	}
}

func (r *CoingeckoProvider) FetchExchangeRate(ctx context.Context, options ExchangeRateOptions) (model.Rate, error) {
	url := r.coingeckoSettings.Host + r.coingeckoSettings.GetExchangeEndpoint

	queryParams := map[string]string{
		coinParameters:     strings.Join(options.Coins, ","),
		currencyParameters: strings.Join(options.Currencies, ","),
		currencyPrecision:  strconv.Itoa(int(options.Precision)),
	}

	resp, err := r.apiClient.MakeAPIRequest(ctx, url, queryParams, nil)
	if err != nil {
		return model.Rate{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Rate{}, err
	}

	var exchangeRates CoingeckoExchangeRateResponse
	err = json.Unmarshal(body, &exchangeRates.Rates)
	if err != nil {
		return model.Rate{}, err
	}

	var rate model.Rate
	rate.Rates = exchangeRates.Rates

	return rate, nil
}
