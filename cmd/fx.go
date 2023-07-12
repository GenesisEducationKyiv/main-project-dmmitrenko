package main

import (
	"CurrencyRateApp/api/controller"
	_ "CurrencyRateApp/api/docs"
	"CurrencyRateApp/api/route"
	constants "CurrencyRateApp/domain"
	"CurrencyRateApp/repository"
	"CurrencyRateApp/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func App() *fx.App {
	return fx.New(
		fx.Provide(
			NewLogger,
			NewEmailRepository,
			NewAPIClient,
			NewCoinMarketProvider,
			NewCoingeckoProvider,
			NewRateService,
			NewEmailService,
			NewEmailController,
			NewRateController,
			RateProviderSlice,
			route.SetupRouter,
		),
		fx.Invoke(startServer),
	)
}

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
	return logger
}

func NewEmailRepository() *repository.EmailRepository {
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

	return repository.NewEmailRepository(writeFile, file)
}

func NewAPIClient(logger *logrus.Logger) *service.ApiClientBase {
	return &service.ApiClientBase{
		Logger: logger,
	}
}

func NewCoinMarketProvider(apiClient *service.ApiClientBase, automapper service.Mapper) *service.CoinMarketProvider {
	return service.NewCoinMarketProvider(automapper, apiClient)
}

func NewCoingeckoProvider(apiClient *service.ApiClientBase, automapper service.Mapper) *service.CoingeckoProvider {
	return service.NewCoingeckoProvider(automapper, apiClient)
}

func RateProviderSlice() []service.RateProvider {
	return []service.RateProvider{
		&service.CoinMarketProvider{},
		&service.CoingeckoProvider{},
	}
}

func NewRateService(logger *logrus.Logger, providers []service.RateProvider) *service.RateService {
	return &service.RateService{
		Providers: providers,
		Logger:    logger,
	}
}

func NewEmailService(emailRepository *repository.EmailRepository, rateService *service.RateService, apiClient *service.ApiClientBase, logger *logrus.Logger) *service.EmailService {
	return &service.EmailService{
		EmailRepository: *emailRepository,
		RateService:     rateService,
		APIClient:       *apiClient,
		Logger:          logger,
	}
}

func NewEmailController(emailService *service.EmailService) *controller.EmailController {
	return controller.NewEmailController(emailService)
}

func NewRateController(rateService *service.RateService) *controller.RateController {
	return controller.NewRateController(rateService)
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

func startServer(router *gin.Engine) {
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
