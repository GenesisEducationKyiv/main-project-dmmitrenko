package integration

import (
	"CurrencyRateApp/api/controller"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockEmailService struct {
	mockError error
}

func (m *MockEmailService) AddEmail(email string) error {
	return m.mockError
}

func (m *MockEmailService) SendRateForSubscribedEmails(coin string, currency string) error {
	return m.mockError
}

func TestIntegrationSubscribeEmail_Success(t *testing.T) {
	// Arrange
	mockService := &MockEmailService{
		mockError: nil,
	}
	controller := controller.NewEmailController(mockService)
	router := gin.Default()
	router.POST("/email", controller.SubscribeEmail)
	recorder := httptest.NewRecorder()
	formData := url.Values{
		"email": []string{"test@example.com"},
	}

	// Act
	request, _ := http.NewRequest("POST", "/email", strings.NewReader(formData.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"message":"Email is successfully subscribed to the newsletter."}`, recorder.Body.String())
}

func TestIntegrationSubscribeEmail_InvalidEmail(t *testing.T) {
	// Arrange
	mockService := &MockEmailService{}
	controller := controller.NewEmailController(mockService)
	router := gin.Default()
	router.POST("/email", controller.SubscribeEmail)
	recorder := httptest.NewRecorder()
	formData := map[string][]string{
		"email": {"invalid_email"},
	}

	// Act
	request, _ := http.NewRequest("POST", "/email", nil)
	request.PostForm = formData
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.JSONEq(t, `{"error":"Invalid email address."}`, recorder.Body.String())
}

func TestIntegrationSendEmails_Success(t *testing.T) {
	// Arrange
	mockService := &MockEmailService{
		mockError: nil,
	}
	controller := controller.NewEmailController(mockService)
	router := gin.Default()
	router.POST("/subscribe", controller.SendEmails)
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/subscribe", nil)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"message":"Letters sent successfully."}`, recorder.Body.String())
}
