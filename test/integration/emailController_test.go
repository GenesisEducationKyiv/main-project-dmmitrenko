package integration

import (
	"CurrencyRateApp/api/controller"
	"CurrencyRateApp/repository"
	"CurrencyRateApp/service"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) CreateLetters(coin string, currency string, currencyRate string, emails []string) []*mail.SGMailV3 {
	args := m.Called(emails)
	return args.Get(0).([]*mail.SGMailV3)
}

func (m *MockEmailService) SendRateForSubscribeEmails(ctx context.Context, coin string, currency string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockEmailService) SubscribeEmail(email string) error {
	args := m.Called()
	return args.Error(0)

}

func TestSubscribeEmailIntegration(t *testing.T) {
	// Arrange
	router := gin.Default()

	tempFile, err := os.CreateTemp("", "test_emails_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	filePath := tempFile.Name()
	defer os.Remove(filePath)

	writer, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	reader, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	logger := logrus.New()
	emailRepository := repository.NewEmailRepository(writer, reader)
	apiClient := service.NewAPIClient(logger)
	coinMarketProvider := service.NewCoinMarketProvider(&service.CoinMarkerExchangeRateResponseMapper{}, apiClient)
	coingeckoProvider := service.NewCoingeckoProvider(&service.CoingeckoExchangeRateResponseMapper{}, apiClient)

	rateService := service.NewRateService(logger, coingeckoProvider, coinMarketProvider)
	emailService := service.NewEmailService(*emailRepository, rateService, *apiClient, logger)

	emailController := controller.NewEmailController(emailService)

	router.POST("/email", emailController.SubscribeEmail)

	formData := url.Values{}
	formData.Set("email", "test@example.com")
	formDataReader := strings.NewReader(formData.Encode())

	req, err := http.NewRequest("POST", "/email", formDataReader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formData.Encode())))

	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	expectedResponseBody := `{"message":"Email is successfully subscribed to the newsletter."}`
	assert.Equal(t, expectedResponseBody, rec.Body.String())

	emails, err := emailRepository.GetAllEmails()
	assert.NoError(t, err)
	assert.Contains(t, emails, "test@example.com")
}

func TestSendEmailsIntegration(t *testing.T) {
	// Arrange
	router := gin.Default()

	tempFile, err := os.CreateTemp("", "test_emails_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	filePath := tempFile.Name()
	defer os.Remove(filePath)

	mockEmailService := &MockEmailService{}
	mockEmailService.On("SendRateForSubscribeEmails").Return(nil)

	emailController := controller.NewEmailController(mockEmailService)

	router.POST("/subscribe", emailController.SendEmails)

	req, err := http.NewRequest("POST", "/subscribe", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	expectedResponseBody := `{"message":"Letters sent successfully."}`
	assert.Equal(t, expectedResponseBody, rec.Body.String())

	mockEmailService.AssertCalled(t, "SendRateForSubscribeEmails")
}
