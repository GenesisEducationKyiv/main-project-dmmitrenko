package service

import (
	constants "CurrencyRateApp/domain"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type RateService struct {
	APIClient APIClient
}

type ApiClient interface {
	MakeAPIRequest(ctx context.Context, url string, queryParams map[string]string) (*http.Response, error)
}

func NewRateService(apiClient APIClient) *RateService {
	return &RateService{
		APIClient: apiClient,
	}
}

const (
	coinParameters     = "ids"
	currencyParameters = "vs_currencies"
	currencyPrecision  = "precision"
)

type ExchangeRateResponse struct {
	Rates map[string]map[string]float64 `json:"rates"`
}

func (r *RateService) FetchExchangeRate(ctx context.Context, coins []string, currencies []string, precision uint) (ExchangeRateResponse, error) {
	url := constants.API_BASE_URL + constants.SIMPLE_PRICE_ENDPOINT

	queryParams := map[string]string{
		coinParameters:     strings.Join(coins, ","),
		currencyParameters: strings.Join(currencies, ","),
		currencyPrecision:  strconv.Itoa(int(precision)),
	}

	resp, err := r.APIClient.MakeAPIRequest(ctx, url, queryParams)
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
