package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Censor Logger middleware started") // Отладка
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = c.Query("request_id")
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		log.Printf("[%s] %s %s", requestID, c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("[%s] %s %s %d", requestID, c.Request.Method, c.Request.URL.Path, c.Writer.Status())
		log.Println("Censor Logger middleware finished") // Отладка
	}
}
