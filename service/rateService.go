package service

import (
	"CurrencyRateApp/domain/model"
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type RateService struct {
	Providers []RateProvider
	Logger    *logrus.Logger
}

type RateProvider interface {
	FetchExchangeRate(ctx context.Context, coins []string, currencies []string, precision uint) (model.Rate, error)
}

func NewRateService(logger *logrus.Logger, providers ...RateProvider) *RateService {
	return &RateService{
		Providers: providers,
		Logger:    logger,
	}
}

func (s *RateService) FetchExchangeRate(ctx context.Context, coins []string, currencies []string, precision uint) (model.Rate, error) {
	var rate model.Rate
	var err error

	for _, provider := range s.Providers {
		rate, err = provider.FetchExchangeRate(ctx, coins, currencies, precision)
		if err == nil {
			return rate, nil
		}

		s.Logger.WithError(err).Warn("Error in getting the rate:")
	}

	return model.Rate{}, fmt.Errorf("невозможно получить курс обмена от доступных провайдеров")
}
