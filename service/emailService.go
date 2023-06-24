package service

import (
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/repository"
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type IEmailService interface {
	AddEmail(email string) error
	SendRateForSubscribedEmails(coin string, currency string) error
}

type EmailService struct {
	EmailRepository *repository.EmailRepository
	RateService     *RateService
	APIClient       *APIClient
}

func NewEmailService(emailRepository *repository.EmailRepository, rateService *RateService, apiClient *APIClient) *EmailService {
	return &EmailService{
		EmailRepository: emailRepository,
		RateService:     rateService,
		APIClient:       apiClient,
	}
}

func (r *EmailService) AddEmail(email string) error {
	return r.EmailRepository.AppendEmailToFile(email)
}

func (r *EmailService) SendRateForSubscribedEmails(coin string, currency string) error {
	emails, err := r.EmailRepository.GetAllEmails()
	if err != nil {
		return err
	}

	rates, err := r.RateService.FetchExchangeRate([]string{coin}, []string{currency}, 2)
	if err != nil {
		return err
	}

	from := mail.NewEmail(os.Getenv(constants.NICKNAME), os.Getenv(constants.EMAIL_SENDER))
	htmlContentTemplate := "<p>RATE: %.2f</p>"

	client := sendgrid.NewSendClient(os.Getenv(constants.APIKEY))

	for _, email := range emails {
		to := mail.NewEmail("", email)

		plainTextContent := fmt.Sprintf("RATE: %.2f", rates.Rates[coin][currency])
		htmlContent := fmt.Sprintf(htmlContentTemplate, rates.Rates[coin][currency])
		subjectTemplate := fmt.Sprintf("%s to %s rate", coin, currency)

		message := mail.NewSingleEmail(from, subjectTemplate, to, plainTextContent, htmlContent)

		response, err := client.Send(message)
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return nil
}
