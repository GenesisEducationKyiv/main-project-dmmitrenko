package e2e

import (
	"CurrencyRateApp/api/controller"
	"CurrencyRateApp/service"
	"bytes"
	"log"

	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeEmail(t *testing.T) {
	// Arrange
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	filePath := createTempFilePath(t, tempDir)
	setEnvFilePath(t, filePath)

	emailService := &service.EmailService{}
	emailController := controller.NewEmailController(emailService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/email", emailController.SubscribeEmail)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	err := writer.WriteField("email", "test1@example.com")
	if err != nil {
		log.Printf("Writing error.")
	}

	writer.Close()

	req, err := http.NewRequest("POST", "/email", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	expectedResponse := `{"message":"Email is successfully subscribed to the newsletter."}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}

func TestSendEmails(t *testing.T) {
	// Arrange
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	filePath := createTempFilePath(t, tempDir)
	setEnvFilePath(t, filePath)

	emailService := &service.EmailService{}
	emailController := controller.NewEmailController(emailService)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/subscribe", emailController.SendEmails)

	req, err := http.NewRequest("POST", "/subscribe", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	// Act
	router.ServeHTTP(recorder, req)

	// Assert
	assert.Equal(t, http.StatusOK, recorder.Code)
	expectedResponse := `{"message":"Letters sent successfully."}`
	assert.Equal(t, expectedResponse, recorder.Body.String())
}

func createTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "test_files")
	if err != nil {
		t.Fatal(err)
	}
	return tempDir
}

func createTempFilePath(t *testing.T, tempDir string) string {
	testFilePath := "test_file.txt"
	filePath := filepath.Join(tempDir, testFilePath)
	return filePath
}

func setEnvFilePath(t *testing.T, filePath string) {
	err := os.Setenv("FILE_PATH", filePath)
	if err != nil {
		t.Fatal(err)
	}
}
