package repository_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"CurrencyRateApp/repository"

	"github.com/stretchr/testify/assert"
)

type mockReader struct {
	data []byte
}

func (r *mockReader) Read(p []byte) (n int, err error) {
	copy(p, r.data)
	return len(r.data), io.EOF
}

type mockWriter struct {
	buffer bytes.Buffer
}

func (w *mockWriter) Write(p []byte) (n int, err error) {
	return w.buffer.Write(p)
}

func TestEmailRepository_AppendEmailToFile(t *testing.T) {
	var buf bytes.Buffer
	emailRepo := repository.NewEmailRepository(&buf, &buf)

	existingEmail := "existing@example.com"
	newEmail := "new@example.com"

	_, err := buf.WriteString(existingEmail + "\n")
	assert.NoError(t, err)

	err = emailRepo.AppendEmailToFile(newEmail)
	assert.NoError(t, err)

	assert.Contains(t, buf.String(), newEmail)
}

func TestEmailRepository_GetAllEmails(t *testing.T) {
	existingEmails := []string{"test1@example.com", "test2@example.com"}

	reader := &mockReader{data: []byte(strings.Join(existingEmails, "\n"))}

	repo := repository.NewEmailRepository(nil, reader)

	emails, err := repo.GetAllEmails()

	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, existingEmails, emails, "incorrect emails returned")
}

func TestEmailRepository_AppendEmailToFile_EmailAlreadyExists(t *testing.T) {
	existingEmails := []string{"test1@example.com", "test2@example.com"}
	existingEmail := existingEmails[0]

	writer := &mockWriter{}

	reader := &mockReader{data: []byte(strings.Join(existingEmails, "\n"))}

	repo := repository.NewEmailRepository(writer, reader)

	err := repo.AppendEmailToFile(existingEmail)

	assert.Error(t, err, "expected error")
}
