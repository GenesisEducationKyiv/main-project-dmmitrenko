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

type ExchangeRateOptions struct {
	Coins      []string
	Currencies []string
	Precision  uint
}

type RateProvider interface {
	FetchExchangeRate(ctx context.Context, options ExchangeRateOptions) (model.Rate, error)
}

func NewRateService(logger *logrus.Logger, providers ...RateProvider) *RateService {
	return &RateService{
		Providers: providers,
		Logger:    logger,
	}
}

func (s *RateService) FetchExchangeRate(ctx context.Context, options ExchangeRateOptions) (model.Rate, error) {
	var rate model.Rate
	var err error

	for _, provider := range s.Providers {
		rate, err = provider.FetchExchangeRate(ctx, options)
		if err == nil {
			return rate, nil
		}

		s.Logger.WithError(err).Warn("error in getting the rate:")
	}

	return model.Rate{}, fmt.Errorf("it is impossible to get exchange rates from available providers")
}
