package unit

import (
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/repository"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAppendEmailToFile_Success(t *testing.T) {
	// Arrange
	tempDir, err := ioutil.TempDir("", "test_files")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	const testFilePath = "test_file.txt"
	filePath := filepath.Join(tempDir, testFilePath)
	os.Setenv(constants.FILE_PATH, filePath)

	existingEmail := "existing_email@example.com"
	err = ioutil.WriteFile(filePath, []byte(existingEmail+"\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	newEmail := "new_email@example.com"

	// Act
	err = repository.AppendEmailToFile(newEmail)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	expectedContent := existingEmail + "\n" + newEmail + "\n"
	if string(fileContent) != expectedContent {
		t.Errorf("Expected file content:\n%s\nGot:\n%s", expectedContent, string(fileContent))
	}
}

func TestAppendEmailToFile_Error(t *testing.T) {
	// Arrange
	tempDir, err := ioutil.TempDir("", "test_files")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	const testFilePath = "test_file.txt"
	filePath := filepath.Join(tempDir, testFilePath)
	os.Setenv(constants.FILE_PATH, filePath)

	existingEmail := "existing_email@example.com"
	err = ioutil.WriteFile(filePath, []byte(existingEmail+"\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	err = repository.AppendEmailToFile(existingEmail)

	// Assert
	if err == nil {
		t.Error("Expected error, got nil")
	} else {
		expectedError := "email already exists: existing_email@example.com"
		if err.Error() != expectedError {
			t.Errorf("Expected error: %s, got: %s", expectedError, err.Error())
		}
	}
}

func TestGetAllEmails_FileNotExist(t *testing.T) {
	// Arrange
	tempDir, err := ioutil.TempDir("", "test_files")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	const testFilePath = "nonexistent_file.txt"
	filePath := filepath.Join(tempDir, testFilePath)
	os.Setenv(constants.FILE_PATH, filePath)

	// Act
	emails, err := repository.GetAllEmails()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(emails) != 0 {
		t.Errorf("Expected empty email list, got: %v", emails)
	}
}
