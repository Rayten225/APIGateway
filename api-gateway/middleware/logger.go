package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем request_id из query или заголовка
		requestID := c.Query("request_id")
		if requestID == "" {
			requestID = c.Request.Header.Get("X-Request-ID")
		}
		if requestID == "" {
			requestID = uuid.New().String()
		}
		// Сохраняем в контексте
		c.Set("request_id", requestID)
		// Устанавливаем заголовок
		c.Header("X-Request-ID", requestID)
		// Логируем
		log.Printf("[%s] %s %s", requestID, c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf("[%s] %s %s %d %s", requestID, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), c.Writer.Header().Get("X-Request-ID"))
	}
}
