package controller

import (
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RateServiceInterface interface {
	FetchExchangeRate(coins []string, currencies []string, precision uint) (service.ExchangeRateResponse, error)
}

type RateController struct {
	rateService RateServiceInterface
}

func NewRateController(rateService RateServiceInterface) *RateController {
	return &RateController{
		rateService: rateService,
	}
}

// GetExchangeRate godoc
// @Summary Get BTC to UAH exchange rate
// @Description Returns the current BTC to UAH exchange rate
// @Tags rate
// @Accept json
// @Produce json
// @Success 200 {number} decimal
// @Router /exchange-rate [get]
func (r *RateController) GetExchangeRate(c *gin.Context) {
	coins := []string{constants.BITCOIN}
	currencies := []string{constants.UAH}

	rates, err := r.rateService.FetchExchangeRate(coins, currencies, 2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchangeRate := rates.Rates[constants.BITCOIN]
	c.JSON(http.StatusOK, exchangeRate[constants.UAH])
}
