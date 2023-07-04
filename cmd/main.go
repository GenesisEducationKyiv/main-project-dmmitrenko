package main

import (
	_ "CurrencyRateApp/api/docs"
	"context"
)

func main() {
	app := App()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}
