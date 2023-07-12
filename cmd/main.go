package main

import (
	"CurrencyRateApp/api/controller"
	_ "CurrencyRateApp/api/docs"
	"CurrencyRateApp/api/route"
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/repository"
	"CurrencyRateApp/service"
	"os"
)

func main() {
	emailService := InitializeEmailService()
	rateService := InitializeRateService(emailService)

	emailController := controller.NewEmailController(emailService)
	rateController := controller.NewRateController(rateService)

	router := route.SetupRouter(emailController, rateController)
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func InitializeEmailService() *service.EmailService {
	filePath := os.Getenv(constants.FILE_PATH)
	createFileIfNotExists(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	writeFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	emailRepository := repository.NewEmailRepository(writeFile, file)
	apiClient := service.NewAPIClient()
	rateService := service.NewRateService(apiClient)

	return &service.EmailService{
		EmailRepository: *emailRepository,
		RateService:     *rateService,
		APIClient:       apiClient,
	}
}

func InitializeRateService(emailService *service.EmailService) *service.RateService {
	apiclient := service.NewAPIClient()

	return &service.RateService{
		APIClient: apiclient,
	}
}

func createFileIfNotExists(filePath string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
	}
}
