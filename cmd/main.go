package main

import (
	"CurrencyRateApp/api/controller"
	_ "CurrencyRateApp/api/docs"
	"CurrencyRateApp/api/route"
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/repository"
	"CurrencyRateApp/service"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	logger.SetLevel((logrus.DebugLevel))
	logger.SetOutput(os.Stdout)

	emailService := InitializeEmailService(logger)
	rateService := InitializeRateService(emailService, logger)

	emailController := controller.NewEmailController(emailService)
	rateController := controller.NewRateController(rateService)

	router := route.SetupRouter(emailController, rateController, logger)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func InitializeEmailService(logger *logrus.Logger) *service.EmailService {
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
	apiClient := service.NewAPIClient(logger)

	coinMarketProvider := &service.CoinMarketProvider{
		Automapper: &service.CoinMarkerExchangeRateResponseMapper{},
		ApiClient:  service.NewAPIClient(logger),
	}
	coingeckoProvider := &service.CoingeckoProvider{
		Automapper: &service.CoingeckoExchangeRateResponseMapper{},
		ApiClient:  service.NewAPIClient(logger),
	}

	rateService := service.NewRateService(logger, coinMarketProvider, coingeckoProvider)

	return &service.EmailService{
		EmailRepository: *emailRepository,
		RateService:     *rateService,
		APIClient:       apiClient,
	}
}

func InitializeRateService(emailService *service.EmailService, logger *logrus.Logger) *service.RateService {
	coinMarketProvider := &service.CoinMarketProvider{
		Automapper: &service.CoinMarkerExchangeRateResponseMapper{},
		ApiClient:  service.NewAPIClient(logger),
	}
	coingeckoProvider := &service.CoingeckoProvider{
		Automapper: &service.CoingeckoExchangeRateResponseMapper{},
		ApiClient:  service.NewAPIClient(logger),
	}

	return service.NewRateService(logger, coingeckoProvider, coinMarketProvider)
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
