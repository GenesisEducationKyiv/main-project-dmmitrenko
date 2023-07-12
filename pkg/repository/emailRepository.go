package repository

import (
	"fmt"
	"io"
	"strings"
)

type FileOptions struct {
	Path string `json:"Path"`
}

type EmailRepository struct {
	writer       io.Writer
	reader       io.Reader
	fileSettings FileOptions
}

func NewEmailRepository(writer io.Writer, reader io.Reader, fileSettings FileOptions) *EmailRepository {
	return &EmailRepository{
		writer:       writer,
		reader:       reader,
		fileSettings: fileSettings,
	}
}

func (r *EmailRepository) AppendEmailToFile(email string) error {
	emails, err := r.GetAllEmails()
	if err != nil {
		return err
	}

	for _, e := range emails {
		if e == email {
			return fmt.Errorf("email already exists: %s", email)
		}
	}

	_, err = fmt.Fprintln(r.writer, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *EmailRepository) GetAllEmails() ([]string, error) {
	var err error

	emailsData, err := io.ReadAll(r.reader)
	if err != nil {
		return nil, err
	}

	emails := strings.Split(strings.TrimSpace(string(emailsData)), "\n")

	return emails, nil
}
