package controller

import (
	constants "CurrencyRateApp/domain"
	_ "CurrencyRateApp/domain/model"
	"CurrencyRateApp/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RateController struct {
	rateService service.RateProvider
}

func NewRateController(rateService service.RateProvider) *RateController {
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
func (r *RateController) GetBitcoinToUahExchangeRate(c *gin.Context) {
	coins := []string{constants.BITCOIN}
	currencies := []string{constants.UAH}
	precision := 2

	var options = service.ExchangeRateOptions{
		Coins:      coins,
		Currencies: currencies,
		Precision:  uint(precision),
	}

	rates, err := r.rateService.FetchExchangeRate(c, options)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exchangeRate := rates.Rates["BTC/UAH"]
	c.JSON(http.StatusOK, exchangeRate)
}

// GetCoinExchangeRate godoc
// @Summary Get the exchange rate for a crypto coin
// @Description Returns the current exchange rate for a crypto coin
// @Tags rate
// @Accept multipart/form-data
// @Produce json
// @Param coins formData string true "Comma-separated list of crypto coins"
// @Param currencies formData string true "Comma-separated list of currencies"
// @Param precision formData string true "Precision of the exchange rate"
// @Success 200
// @Failure 400
// @Router /exchange-rate [post]
func (r *RateController) GetCoinExchangeRate(c *gin.Context) {
	coins := c.PostForm("coins")
	currencies := c.PostForm("currencies")
	precisionStr := c.PostForm("precision")

	precision, err := strconv.ParseUint(precisionStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var options = service.ExchangeRateOptions{
		Coins:      strings.Split(coins, ","),
		Currencies: strings.Split(currencies, ","),
		Precision:  uint(precision),
	}

	rates, err := r.rateService.FetchExchangeRate(c, options)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rates)
}
