package api

import (
	"CurrencyRateApp/internal/model"
	"CurrencyRateApp/pkg/external"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type stubRateService struct{}

func (m *stubRateService) FetchExchangeRate(ctx context.Context, options external.ExchangeRateOptions) (model.Rate, error) {
	return model.Rate{
		Rates: map[string]map[string]float64{
			"BTC": {
				"UAH": 450000.0,
			},
		},
	}, nil
}

func TestGetBitcoinToUahExchangeRate(t *testing.T) {
	// Arrange
	router := gin.Default()

	rateService := &stubRateService{}

	rateController := NewRateController(rateService)

	router.GET("/exchange-rate", rateController.GetBitcoinToUahExchangeRate)

	req, err := http.NewRequest("GET", "/exchange-rate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	expectedResponseBody := `450000`
	assert.Equal(t, expectedResponseBody, rec.Body.String())
}
