package service

import (
	"CurrencyRateApp/repository"
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

type SenderOptions struct {
	Nickname    string `json:"Nickname"`
	EmailSender string `json:"EmailSender"`
	ApiKey      string `json:"ApiKey"`
}

type EmailService struct {
	emailRepository repository.EmailRepository
	rateService     *RateService
	apiClient       ApiClientBase
	logger          *logrus.Logger
	senderSettings  SenderOptions
}

func NewEmailService(
	emailRepository repository.EmailRepository,
	rateService *RateService,
	apiClient ApiClientBase,
	logger *logrus.Logger,
	senderSettings SenderOptions) *EmailService {
	return &EmailService{
		emailRepository: emailRepository,
		rateService:     rateService,
		apiClient:       apiClient,
		logger:          logger,
		senderSettings:  senderSettings,
	}
}

func (r *EmailService) SubscribeEmail(email string) error {
	return r.emailRepository.AppendEmailToFile(email)
}

func (r *EmailService) SendRateForSubscribeEmails(ctx context.Context, coin string, currency string) error {
	emails, err := r.emailRepository.GetAllEmails()
	if err != nil {
		return err
	}

	var options = ExchangeRateOptions{
		Coins:      []string{coin},
		Currencies: []string{currency},
		Precision:  2,
	}

	rates, err := r.rateService.FetchExchangeRate(ctx, options)
	if err != nil {
		return err
	}

	emailsToSend := r.CreateLetters(coin, currency, fmt.Sprintf("%b", rates.Rates[""]), emails)
	client := sendgrid.NewSendClient(r.senderSettings.ApiKey)

	for _, emailToSend := range emailsToSend {
		response, err := client.Send(emailToSend)
		if err != nil {
			return err
		}
		r.logger.Log(logrus.InfoLevel, response.StatusCode)
	}

	return nil
}

func (r *EmailService) CreateLetters(coin string, currency string, currencyRate string, emails []string) []*mail.SGMailV3 {
	from := mail.NewEmail(r.senderSettings.Nickname, r.senderSettings.EmailSender)

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
