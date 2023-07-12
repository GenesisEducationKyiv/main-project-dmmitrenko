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

	rates := m.flattenCoinMarkerExchangeRateResponse(response)

	return model.Rate{
		Rates: rates,
	}, nil
}

func (m *CoinMarkerExchangeRateResponseMapper) flattenCoinMarkerExchangeRateResponse(response CoinMarkerExchangeRateResponse) map[string]float64 {
	flattenedRates := make(map[string]float64)
	for currency, quote := range response.Data {
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

	rates := m.flattenCoingeckoExchangeRateResponse(response)

	return model.Rate{
		Rates: rates,
	}, nil
}

func (m *CoingeckoExchangeRateResponseMapper) flattenCoingeckoExchangeRateResponse(response CoingeckoExchangeRateResponse) map[string]float64 {
	flattenedRates := make(map[string]float64)
	for currency, rateMap := range response.Rates {
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
