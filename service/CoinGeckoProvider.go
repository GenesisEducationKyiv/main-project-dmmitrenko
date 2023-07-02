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
	Automapper Mapper
	ApiClient  ApiClientBase
}

func (r *CoingeckoProvider) FetchExchangeRate(ctx context.Context, coins []string, currencies []string, precision uint) (model.Rate, error) {
	url := constants.API_BASE_URL + constants.SIMPLE_PRICE_ENDPOINT

	queryParams := map[string]string{
		coinParameters:     strings.Join(coins, ","),
		currencyParameters: strings.Join(currencies, ","),
		currencyPrecision:  strconv.Itoa(int(precision)),
	}

	resp, err := r.ApiClient.MakeAPIRequest(ctx, url, queryParams, nil)
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
	rate, err = r.Automapper.MapToRate(exchangeRates)
	if err != nil {
		return model.Rate{}, err
	}

	return rate, nil
}
