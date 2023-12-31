package route

import (
	"CurrencyRateApp/api/controller"
	"CurrencyRateApp/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(emailController *controller.EmailController, rateController *controller.RateController, logger *logrus.Logger) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.ExceptionMiddleware(logger))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/exchange-rate", rateController.GetBitcoinToUahExchangeRate)
	router.POST("/exchange-rate", rateController.GetCoinExchangeRate)
	router.POST("/email", emailController.SubscribeEmail)
	router.POST("/subscribe", emailController.SendEmails)

	return router
}
