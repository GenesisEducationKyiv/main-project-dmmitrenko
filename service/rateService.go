package service

import (
	"CurrencyRateApp/domain/model"
	"context"
	"fmt"
)

type RateService struct {
	providers []RateProvider
}

type RateProvider interface {
	FetchExchangeRate(ctx context.Context, coins []string, currencies []string, precision uint) (model.Rate, error)
}

func NewRateService(providers ...RateProvider) *RateService {
	return &RateService{
		providers: providers,
	}
}

func (s *RateService) FetchExchangeRate(ctx context.Context, coins []string, currencies []string, precision uint) (model.Rate, error) {
	var rate model.Rate
	var err error

	for _, provider := range s.providers {
		rate, err = provider.FetchExchangeRate(ctx, coins, currencies, precision)
		if err == nil {
			return rate, nil
		}

		fmt.Printf("Ошибка при получении курса: %v\n", err)
	}

	return model.Rate{}, fmt.Errorf("невозможно получить курс обмена от доступных провайдеров")
}
