package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get("X-Request-ID") // Сначала проверяем заголовок
		if reqID == "" {
			reqID = c.Query("request_id")
		}
		if reqID == "" {
			reqID = uuid.New().String() // Генерируем UUID только если нет заголовка или query
		}
		ctx := context.WithValue(c.Request.Context(), RequestIDKey, reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Set("request_id", reqID) // Сохраняем в контексте Gin для Logger
		c.Writer.Header().Set("X-Request-ID", reqID)
		c.Next()
	}
}
