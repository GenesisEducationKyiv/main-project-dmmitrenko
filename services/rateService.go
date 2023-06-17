package services

import (
	constants "CurrencyRateApp/domain"
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

type ExchangeRateResponse struct {
	Rates map[string]map[string]float64 `json:"rates"`
}

func FetchExchangeRate(coins []string, currencies []string, precision uint) (ExchangeRateResponse, error) {
	url := constants.ApiBaseUrl + constants.SimplePriceEndpoint

	queryParams := map[string]string{
		coinParameters:     strings.Join(coins, ","),
		currencyParameters: strings.Join(currencies, ","),
		currencyPrecision:  strconv.Itoa(int(precision)),
	}

	resp, err := makeAPIRequest(url, queryParams)
	if err != nil {
		return ExchangeRateResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ExchangeRateResponse{}, err
	}

	var exchangeRates ExchangeRateResponse
	err = json.Unmarshal(body, &exchangeRates.Rates)
	if err != nil {
		return ExchangeRateResponse{}, err
	}

	return exchangeRates, nil
}
