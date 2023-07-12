package service

import (
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/domain/model"
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

type CoingeckoProvider struct {
	automapper Mapper
	apiClient  APIClient
}

func NewCoingeckoProvider(automapper Mapper, apiClient APIClient) *CoingeckoProvider {
	return &CoingeckoProvider{
		automapper: automapper,
		apiClient:  apiClient,
	}
}

func (r *CoingeckoProvider) FetchExchangeRate(ctx context.Context, options ExchangeRateOptions) (model.Rate, error) {
	url := constants.API_BASE_URL + constants.SIMPLE_PRICE_ENDPOINT

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
	rate, err = r.automapper.MapToRate(exchangeRates)
	if err != nil {
		return model.Rate{}, err
	}

	return rate, nil
}
