package service

import (
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/repository"
	"context"
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	EmailRepository repository.EmailRepository
	RateService     RateService
	APIClient       ApiClientBase
}

func NewEmailService(emailRepository repository.EmailRepository, rateService RateService, apiClient ApiClientBase) *EmailService {
	return &EmailService{
		EmailRepository: emailRepository,
		RateService:     rateService,
		APIClient:       apiClient,
	}
}

func (r *EmailService) SubscribeEmail(email string) error {
	return r.EmailRepository.AppendEmailToFile(email)
}

func (r *EmailService) SendRateForSubscribeEmails(ctx context.Context, coin string, currency string) error {
	emails, err := r.EmailRepository.GetAllEmails()
	if err != nil {
		return err
	}

	rates, err := r.RateService.FetchExchangeRate(ctx, []string{coin}, []string{currency}, 2)
	if err != nil {
		return err
	}

	emailsToSend := r.CreateLetters(coin, currency, fmt.Sprintf("%b", rates.Rates[""]), emails)
	client := sendgrid.NewSendClient(os.Getenv(constants.APIKEY))

	for _, emailToSend := range emailsToSend {
		response, err := client.Send(emailToSend)
		if err != nil {
			return err
		}
		fmt.Println(response.StatusCode)
	}

	return nil
}

func (r *EmailService) CreateLetters(coin string, currency string, currencyRate string, emails []string) []*mail.SGMailV3 {
	from := mail.NewEmail(os.Getenv(constants.NICKNAME), os.Getenv(constants.EMAIL_SENDER))

	htmlContentTemplate := "<p>RATE: %s</p>"
	plainTextContent := fmt.Sprintf("RATE: %s", currencyRate)
	htmlContent := fmt.Sprintf(htmlContentTemplate, currencyRate)
	subjectTemplate := fmt.Sprintf("%s to %s rate", coin, currency)

	var letters []*mail.SGMailV3
	for _, email := range emails {
		to := mail.NewEmail("", email)
		letter := mail.NewSingleEmail(from, subjectTemplate, to, plainTextContent, htmlContent)
		letters = append(letters, letter)
	}

	return letters
}
