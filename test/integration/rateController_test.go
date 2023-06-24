package integration

import (
	"CurrencyRateApp/api/controller"
	"CurrencyRateApp/service"
	"net/http"
	"net/http/httptest"
	"testing"

	constants "CurrencyRateApp/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockRateService struct {
	mockRates service.ExchangeRateResponse
	mockError error
}

func (m *MockRateService) FetchExchangeRate(coins []string, currencies []string, precision uint) (service.ExchangeRateResponse, error) {
	return m.mockRates, m.mockError
}

func TestIntegrationGetExchangeRate(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()

	// Create a mock service
	mockRates := service.ExchangeRateResponse{
		Rates: map[string]map[string]float64{
			constants.BITCOIN: {
				constants.UAH: 40000,
			},
		},
	}
	mockService := &MockRateService{
		mockRates: mockRates,
		mockError: nil,
	}

	// Create a new rate controller with the mock service
	controller := controller.NewRateController(mockService)

	// Register the route with the controller's handler function
	router.GET("/exchange-rate", controller.GetExchangeRate)

	// Create a test request
	req, err := http.NewRequest("GET", "/exchange-rate", nil)
	assert.NoError(t, err)

	// Create a test recorder to capture the response
	recorder := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(recorder, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `40000`, recorder.Body.String())
}
