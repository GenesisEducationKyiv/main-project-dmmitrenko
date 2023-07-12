package main

// import (
// 	"CurrencyRateApp/api/controller"
// 	_ "CurrencyRateApp/api/docs"
// 	"CurrencyRateApp/api/route"
// 	"CurrencyRateApp/repository"
// 	"CurrencyRateApp/service"
// 	"encoding/json"
// 	"io"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"github.com/sirupsen/logrus"
// 	"go.uber.org/fx"
// )

// func App() *fx.App {
// 	return fx.New(
// 		fx.Provide(
// 			NewAppSettings,
// 			NewCoingeckoSettings,
// 			NewCoinMarketSettings,
// 			NewFileSettings,
// 			NewSenderSettings,
// 			NewLogger,
// 			NewEmailRepository,
// 			NewAPIClient,
// 			NewCoinMarketProvider,
// 			NewCoingeckoProvider,
// 			NewRateService,
// 			NewEmailService,
// 			NewEmailController,
// 			NewRateController,
// 			RateProviderSlice,
// 			route.SetupRouter,
// 		),
// 		fx.Invoke(startServer),
// 	)
// }

// func NewLogger() *logrus.Logger {
// 	logger := logrus.New()
// 	logger.SetFormatter(&logrus.TextFormatter{
// 		ForceColors: true,
// 	})
// 	logger.SetLevel(logrus.DebugLevel)
// 	logger.SetOutput(os.Stdout)
// 	return logger
// }

// func NewEmailRepository(fileSettings repository.FileOptions) *repository.EmailRepository {
// 	filePath := fileSettings.Path
// 	createFileIfNotExists(filePath)

// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	writeFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return repository.NewEmailRepository(writeFile, file, fileSettings)
// }

// func NewAPIClient(logger *logrus.Logger) *service.ApiClientBase {
// 	return service.NewAPIClient(logger)
// }

// func NewCoinMarketProvider(apiClient *service.ApiClientBase, automapper service.Mapper, coinMarketSettings service.CoinMarketOptions) *service.CoinMarketProvider {
// 	return service.NewCoinMarketProvider(automapper, apiClient, coinMarketSettings)
// }

// func NewCoingeckoProvider(apiClient *service.ApiClientBase, automapper service.Mapper, coingeckoSettings service.CoingeckoOptions) *service.CoingeckoProvider {
// 	return service.NewCoingeckoProvider(automapper, apiClient, coingeckoSettings)
// }

// func RateProviderSlice() []service.RateProvider {
// 	return []service.RateProvider{
// 		&service.CoinMarketProvider{},
// 		&service.CoingeckoProvider{},
// 	}
// }

// func NewRateService(logger *logrus.Logger, providers []service.RateProvider) *service.RateService {
// 	return service.NewRateService(logger, providers...)
// }

// func NewEmailService(emailRepository *repository.EmailRepository, rateService *service.RateService, apiClient *service.ApiClientBase, logger *logrus.Logger, senderSettings service.SenderOptions) *service.EmailService {
// 	return service.NewEmailService(*emailRepository, rateService, *apiClient, logger, senderSettings)
// }

// func NewEmailController(emailService *service.EmailService) *controller.EmailController {
// 	return controller.NewEmailController(emailService)
// }

// func NewRateController(rateService *service.RateService) *controller.RateController {
// 	return controller.NewRateController(rateService)
// }

// func NewCoinMarketSettings(configuration Configuration) service.CoinMarketOptions {
// 	return configuration.CoinMarketSettings
// }

// func NewCoingeckoSettings(configuration Configuration) service.CoingeckoOptions {
// 	return configuration.CoingeckoSettings
// }

// func NewSenderSettings(configuration Configuration) service.SenderOptions {
// 	return configuration.SenderSettings
// }

// func NewFileSettings(Configuration Configuration) repository.FileOptions {
// 	return Configuration.FileSettings
// }

// func NewAppSettings() Configuration {
// 	configFile, err := os.Open("appsettings.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer configFile.Close()

// 	config, err := io.ReadAll(configFile)
// 	if err != nil {
// 		panic(err)
// 	}

// 	appConfig := Configuration{}
// 	err = json.Unmarshal(config, &appConfig)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return appConfig
// }

// func createFileIfNotExists(filePath string) {
// 	_, err := os.Stat(filePath)
// 	if os.IsNotExist(err) {
// 		_, err := os.Create(filePath)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }

// func startServer(router *gin.Engine) {
// 	err := router.Run(":8080")
// 	if err != nil {
// 		panic(err)
// 	}
// }
