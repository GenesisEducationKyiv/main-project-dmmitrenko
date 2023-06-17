package repositories

import (
	constants "CurrencyRateApp/domain"
	"fmt"
	"os"
	"strings"
)

func AppendEmailToFile(email string) error {
	emails, err := GetAllEmails()
	if err != nil {
		return err
	}

	for _, e := range emails {
		if e == email {
			return fmt.Errorf("email already exists: %s", email)
		}
	}

	file, err := os.OpenFile(constants.FilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(email + "\n")
	if err != nil {
		return err
	}

	return nil
}

func GetAllEmails() ([]string, error) {
	if _, err := os.Stat(constants.FilePath); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		_, err := os.Create(constants.FilePath)
		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(constants.FilePath)
	if err != nil {
		return nil, err
	}

	emails := strings.Split(strings.TrimSpace(string(data)), "\n")

	return emails, nil
}
