package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExceptionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Recovered from panic:", err)

				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				c.Abort()
			}
		}()

		c.Next()
	}
}
