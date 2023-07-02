package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

type mockEmailService struct {
	SubscribeEmailErr error
}

func (m *mockEmailService) SubscribeEmail(email string) error {
	return m.SubscribeEmailErr
}

func (m *mockEmailService) SendRateForSubscribeEmails(ctx context.Context, coin string, currency string) error {
	return nil
}

func (m *mockEmailService) CreateLetters(coin string, currency string, currencyRate string, emails []string) []*mail.SGMailV3 {
	return nil
}

func TestSubscribeEmail(t *testing.T) {
	t.Run("Successful subscription", func(t *testing.T) {
		// Arrange
		router := gin.Default()

		emailService := &mockEmailService{}
		emailController := NewEmailController(emailService)

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
	})

	t.Run("Invalid email address", func(t *testing.T) {
		// Arrange
		router := gin.Default()

		emailService := &mockEmailService{}
		emailController := NewEmailController(emailService)

		router.POST("/email", emailController.SubscribeEmail)

		formData := url.Values{}
		formData.Set("email", "invalid_email")
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
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedResponseBody := `{"error":"Invalid email address."}`
		assert.Equal(t, expectedResponseBody, rec.Body.String())
	})

	t.Run("Email service error", func(t *testing.T) {
		// Arrange
		router := gin.Default()

		emailService := &mockEmailService{
			SubscribeEmailErr: errors.New("failed to subscribe email"),
		}
		emailController := NewEmailController(emailService)

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
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectedResponseBody := `{"error":"failed to subscribe email"}`
		assert.Equal(t, expectedResponseBody, rec.Body.String())
	})
}
