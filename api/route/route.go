package route

import (
	"CurrencyRateApp/api/controller"
	"CurrencyRateApp/api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(emailController *controller.EmailController, rateController *controller.RateController) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.ExceptionMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/exchange-rate", rateController.GetExchangeRate)
	router.POST("/email", emailController.SubscribeEmail)
	router.POST("/subscribe", emailController.SendEmails)

	return router
}
