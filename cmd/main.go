package main

import (
	_ "CurrencyRateApp/api/docs"
	"CurrencyRateApp/repository"
	"CurrencyRateApp/service"
	"context"
)

type Configuration struct {
	CoinMarketSettings service.CoinMarketOptions `json:"CoinMarketOptions"`
	CoingeckoSettings  service.CoingeckoOptions  `json:"CoingeckoOptions"`
	SenderSettings     service.SenderOptions     `json:"SenderOptions"`
	FileSettings       repository.FileOptions    `json:"FileOptions"`
}

func main() {
	app := App()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}
