package service

import (
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/repository"
	"context"
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

type EmailService struct {
	EmailRepository repository.EmailRepository
	RateService     *RateService
	APIClient       ApiClientBase
	Logger          *logrus.Logger
}

func NewEmailService(emailRepository repository.EmailRepository, rateService *RateService, apiClient ApiClientBase, logger *logrus.Logger) *EmailService {
	return &EmailService{
		EmailRepository: emailRepository,
		RateService:     rateService,
		APIClient:       apiClient,
		Logger:          logger,
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

	var options = ExchangeRateOptions{
		Coins:      []string{coin},
		Currencies: []string{currency},
		Precision:  2,
	}

	rates, err := r.RateService.FetchExchangeRate(ctx, options)
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
		r.Logger.Log(logrus.InfoLevel, response.StatusCode)
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
