package service

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/repository"
)

func AddEmail(email string) error {
	return repository.AppendEmailToFile(email)
}

func SendRateForSubscribedEmails(coin string, currency string) error {
	emails, err := repository.GetAllEmails()
	if err != nil {
		return err
	}

	rates, err := FetchExchangeRate([]string{coin}, []string{currency}, 2)
	if err != nil {
		return err
	}

	from := mail.NewEmail(os.Getenv(constants.Nickname), os.Getenv(constants.EmailSender))
	htmlContentTemplate := "<p>RATE: %.2f</p>"

	client := sendgrid.NewSendClient(os.Getenv(constants.ApiKey))

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
