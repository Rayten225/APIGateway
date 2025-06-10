package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем request_id из заголовка
		requestID := c.Request.Header.Get("X-Request-ID")
		if requestID == "" {
			// Если заголовок пуст, генерируем новый UUID
			requestID = uuid.New().String()
		}
		// Сохраняем request_id в контексте
		c.Set("request_id", requestID)
		// Устанавливаем заголовок в ответе
		c.Header("X-Request-ID", requestID)
		// Логируем начало запроса
		log.Printf("[%s] %s %s", requestID, c.Request.Method, c.Request.URL.Path)
		c.Next()
		// Логируем завершение запроса
		log.Printf("[%s] %s %s %d %s", requestID, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), c.Writer.Header().Get("X-Request-ID"))
	}
}
