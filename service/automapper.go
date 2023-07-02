package service

import (
	"CurrencyRateApp/domain/model"
	"errors"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type Mapper interface {
	MapToRate(input interface{}) (model.Rate, error)
}

type ExchangeRateMapper struct{}
type CoinMarkerExchangeRateResponseMapper struct{}

func (m *CoinMarkerExchangeRateResponseMapper) MapToRate(input interface{}) (model.Rate, error) {
	response, ok := input.(CoinMarkerExchangeRateResponse)
	if !ok {
		return model.Rate{}, errors.New("invalid input type")
	}

	rates := m.flattenCoinMarkerExchangeRateResponse(response.Data)

	return model.Rate{
		Rates: rates,
	}, nil
}

func (m *CoinMarkerExchangeRateResponseMapper) flattenCoinMarkerExchangeRateResponse(data map[string]struct {
	Quote map[string]struct {
		Price float64 `json:"price"`
	} `json:"quote"`
}) map[string]float64 {
	flattenedRates := make(map[string]float64)
	for currency, quote := range data {
		for targetCurrency, rate := range quote.Quote {
			key := currency + "/" + targetCurrency
			flattenedRates[key] = rate.Price
		}
	}
	return flattenedRates
}

type CoingeckoExchangeRateResponseMapper struct{}

func (m *CoingeckoExchangeRateResponseMapper) MapToRate(input interface{}) (model.Rate, error) {
	response, ok := input.(CoingeckoExchangeRateResponse)
	if !ok {
		return model.Rate{}, errors.New("invalid input type")
	}

	rates := m.flattenCoingeckoExchangeRateResponse(response.Rates)

	return model.Rate{
		Rates: rates,
	}, nil
}

func (m *CoingeckoExchangeRateResponseMapper) flattenCoingeckoExchangeRateResponse(rates map[string]map[string]float64) map[string]float64 {
	flattenedRates := make(map[string]float64)
	for currency, rateMap := range rates {
		for targetCurrency, rate := range rateMap {
			normalizedCurrency := model.NormalizeCurrency(strings.ToLower(currency))
			key := normalizedCurrency + "/" + strings.ToUpper(targetCurrency)
			flattenedRates[key] = rate
		}
	}
	return flattenedRates
}

func (m *ExchangeRateMapper) MapToRate(input interface{}) (model.Rate, error) {
	var rate model.Rate
	err := mapstructure.Decode(input, &rate)
	return rate, err
}
