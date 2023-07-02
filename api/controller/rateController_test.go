package controller

import (
	"CurrencyRateApp/service"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockRateService struct{}

func (m *mockRateService) FetchExchangeRate(ctx context.Context, coins []string, currencies []string, precision uint) (service.ExchangeRateResponse, error) {
	return service.ExchangeRateResponse{
		Rates: map[string]map[string]float64{
			"bitcoin": {
				"uah": 450000.0,
			},
		},
	}, nil
}

func TestGetBitcoinToUahExchangeRate(t *testing.T) {
	// Arrange
	router := gin.Default()

	rateService := &mockRateService{}

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
