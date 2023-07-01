package e2e

// import (
// 	"CurrencyRateApp/api/controller"
// 	"CurrencyRateApp/service"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGetExchangeRate(t *testing.T) {
// 	// Arrange
// 	rateService := &service.RateService{}
// 	rateController := controller.NewRateController(rateService)

// 	gin.SetMode(gin.TestMode)
// 	router := gin.Default()
// 	router.GET("/exchange-rate", rateController.GetExchangeRate)

// 	req, err := http.NewRequest("GET", "/exchange-rate", nil)
// 	assert.NoError(t, err)

// 	recorder := httptest.NewRecorder()

// 	// Act
// 	router.ServeHTTP(recorder, req)

// 	// Assert
// 	assert.Equal(t, http.StatusOK, recorder.Code)

// 	var response float64
// 	err = json.Unmarshal(recorder.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// }
