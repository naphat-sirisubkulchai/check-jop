package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("🌐 %s | %s | %d | %v | %s | %s %s | %s\n",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.ClientIP,
			param.StatusCode,
			param.Latency,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.ErrorMessage,
		)
	})
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("RequestID", requestID)
		c.Next()
	}
}

func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
