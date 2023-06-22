package controller

import (
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetExchangeRate godoc
// @Summary Get BTC to UAH exchange rate
// @Description Returns the current BTC to UAH exchange rate
// @Tags rate
// @Accept json
// @Produce json
// @Success 200 {number} decimal
// @Router /exchange-rate [get]
func GetExchangeRate(c *gin.Context) {
	coins := []string{constants.Bitcoin}
	currencies := []string{constants.UAH}

	rates, err := service.FetchExchangeRate(coins, currencies, 2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchangeRate := rates.Rates[constants.Bitcoin]
	c.JSON(http.StatusOK, exchangeRate[constants.UAH])
}
