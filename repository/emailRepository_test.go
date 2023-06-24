package repository

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendEmailToFile_Success(t *testing.T) {
	// Arrange
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	filePath := createTempFilePath(t, tempDir)
	setEnvFilePath(t, filePath)

	existingEmail := "existing_email@example.com"
	writeToFile(t, filePath, existingEmail+"\n")

	newEmail := "new_email@example.com"

	repo := NewEmailRepository()

	// Act
	err := repo.AppendEmailToFile(newEmail)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	fileContent := readFile(t, filePath)

	expectedContent := existingEmail + "\n" + newEmail + "\n"
	if fileContent != expectedContent {
		t.Errorf("Expected file content:\n%s\nGot:\n%s", expectedContent, fileContent)
	}
}

func TestAppendEmailToFile_Error(t *testing.T) {
	// Arrange
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	filePath := createTempFilePath(t, tempDir)
	setEnvFilePath(t, filePath)

	existingEmail := "existing_email@example.com"
	writeToFile(t, filePath, existingEmail+"\n")

	repo := NewEmailRepository()

	// Act
	err := repo.AppendEmailToFile(existingEmail)

	// Assert
	expectedError := "email already exists: existing_email@example.com"
	if err == nil {
		t.Error("Expected error, got nil")
	} else if err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestGetAllEmails_FileNotExist(t *testing.T) {
	// Arrange
	tempDir := createTempDir(t)
	defer os.RemoveAll(tempDir)

	filePath := createTempFilePath(t, tempDir)
	setEnvFilePath(t, filePath)

	repo := NewEmailRepository()

	// Act
	emails, err := repo.GetAllEmails()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(emails) != 0 {
		t.Errorf("Expected empty email list, got: %v", emails)
	}
}

// Helper functions

func createTempDir(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "test_files")
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

func writeToFile(t *testing.T, filePath, content string) {
	err := ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, filePath string) string {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return string(fileContent)
}
