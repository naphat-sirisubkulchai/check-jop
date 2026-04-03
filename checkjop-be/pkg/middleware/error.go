package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("💥 Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "Something went wrong",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func ValidationErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("⚠️ Request error: %v", err.Error())
		}
	}
}
