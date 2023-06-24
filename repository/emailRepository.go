package repository

import (
	constants "CurrencyRateApp/domain"
	"fmt"
	"os"
	"strings"
)

type EmailRepository struct{}

func NewEmailRepository() *EmailRepository {
	return &EmailRepository{}
}

func (r *EmailRepository) AppendEmailToFile(email string) error {
	err := createFileIfNotExists()
	if err != nil {
		return err
	}

	emails, err := r.GetAllEmails()
	if err != nil {
		return err
	}

	for _, e := range emails {
		if e == email {
			return fmt.Errorf("email already exists: %s", email)
		}
	}

	file, err := os.OpenFile(os.Getenv(constants.FILE_PATH), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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

func (r *EmailRepository) GetAllEmails() ([]string, error) {
	var err error
	_, err = os.Stat(os.Getenv(constants.FILE_PATH))
	if os.IsNotExist(err) {
		return []string{}, nil
	}

	data, err := os.ReadFile(os.Getenv(constants.FILE_PATH))
	if err != nil {
		return nil, err
	}

	emails := strings.Split(strings.TrimSpace(string(data)), "\n")

	return emails, nil
}

func createFileIfNotExists() error {
	if _, err := os.Stat(os.Getenv(constants.FILE_PATH)); os.IsNotExist(err) {
		_, err := os.Create(os.Getenv(constants.FILE_PATH))
		if err != nil {
			return err
		}
	}
	return nil
}
