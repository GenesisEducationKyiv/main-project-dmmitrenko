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
