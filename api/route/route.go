package route

import (
	"CurrencyRateApp/api/controller"
	"CurrencyRateApp/api/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.ExceptionMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/exchange-rate", controller.GetExchangeRate)
	router.POST("/email", controller.SubscribeEmail)
	router.POST("/subscribe", controller.SendEmails)

	return router
}
