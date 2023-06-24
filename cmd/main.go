package main

import (
	"CurrencyRateApp/api/controller"
	_ "CurrencyRateApp/api/docs"
	"CurrencyRateApp/api/route"
	"CurrencyRateApp/repository"
	"CurrencyRateApp/service"
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
	emailRepository := repository.NewEmailRepository()
	apiClient := service.NewAPIClient()

	return &service.EmailService{
		EmailRepository: emailRepository,
		APIClient:       apiClient,
	}
}

func InitializeRateService(emailService *service.EmailService) *service.RateService {
	apiClient := service.NewAPIClient()

	return &service.RateService{
		ApiClient: apiClient,
	}
}
