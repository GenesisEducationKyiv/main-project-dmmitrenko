package main

import (
	"CurrencyRateApp/pkg/api"
	"CurrencyRateApp/pkg/external"
	"CurrencyRateApp/pkg/repository"
	"CurrencyRateApp/pkg/service"
	"encoding/json"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

func App() *fx.App {
	return fx.New(
		fx.Provide(
			NewAppSettings,
			NewCoingeckoSettings,
			NewCoinMarketSettings,
			NewFileSettings,
			NewSenderSettings,
			NewLogger,
			NewEmailRepository,
			NewAPIClient,
			NewCoinMarketProvider,
			NewCoingeckoProvider,
			RateProviderSlice,
			NewRateService,
			NewEmailService,
			NewEmailController,
			NewRateController,
			api.SetupRouter,
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

func NewEmailRepository(configuration Configuration) *repository.EmailRepository {
	filePath := configuration.FileSettings.Path
	createFileIfNotExists(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	writeFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	return repository.NewEmailRepository(writeFile, file, configuration.FileSettings)
}

func NewAPIClient(logger *logrus.Logger) *external.ApiClientBase {
	return external.NewAPIClient(logger)
}

func NewCoinMarketProvider(apiClient *external.ApiClientBase, coinMarketSettings external.CoinMarketOptions) *external.CoinMarketProvider {
	return external.NewCoinMarketProvider(apiClient, coinMarketSettings)
}

func NewCoingeckoProvider(apiClient *external.ApiClientBase, coingeckoSettings external.CoingeckoOptions) *external.CoingeckoProvider {
	return external.NewCoingeckoProvider(apiClient, coingeckoSettings)
}

func RateProviderSlice(coingeckoProvider *external.CoingeckoProvider, coinMarketProvider *external.CoinMarketProvider) []service.RateProvider {
	return []service.RateProvider{
		coingeckoProvider,
		coinMarketProvider,
	}
}

func NewRateService(logger *logrus.Logger, providers []service.RateProvider) *service.RateService {
	return service.NewRateService(logger, providers...)
}

func NewEmailService(emailRepository *repository.EmailRepository, rateService *service.RateService, apiClient *external.ApiClientBase, logger *logrus.Logger, senderSettings service.SenderOptions) *service.EmailService {
	return service.NewEmailService(*emailRepository, rateService, *apiClient, logger, senderSettings)
}

func NewEmailController(emailService *service.EmailService) *api.EmailController {
	return api.NewEmailController(emailService)
}

func NewRateController(rateService *service.RateService) *api.RateController {
	return api.NewRateController(rateService)
}

func NewCoinMarketSettings(configuration Configuration) external.CoinMarketOptions {
	return configuration.CoinMarketSettings
}

func NewCoingeckoSettings(configuration Configuration) external.CoingeckoOptions {
	return configuration.CoingeckoSettings
}

func NewSenderSettings(configuration Configuration) service.SenderOptions {
	return configuration.SenderSettings
}

func NewFileSettings(configuration Configuration) repository.FileOptions {
	return configuration.FileSettings
}

func NewAppSettings() Configuration {
	configFile, err := os.Open("../config/appsettings.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	config, err := io.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	appConfig := Configuration{}
	err = json.Unmarshal(config, &appConfig)
	if err != nil {
		panic(err)
	}

	return appConfig
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
