package main

import (
	"CurrencyRateApp/pkg/external"
	"CurrencyRateApp/pkg/repository"
	"CurrencyRateApp/pkg/service"
	"context"
)

type Configuration struct {
	CoinMarketSettings external.CoinMarketOptions `json:"CoinMarketOptions"`
	CoingeckoSettings  external.CoingeckoOptions  `json:"CoingeckoOptions"`
	SenderSettings     service.SenderOptions      `json:"SenderOptions"`
	FileSettings       repository.FileOptions     `json:"FileOptions"`
}

func main() {
	app := App()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}
