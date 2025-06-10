package router

import (
	"github.com/gin-gonic/gin"
	"micronews/api-gateway/handlers"
	"micronews/api-gateway/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Logger())

	api := r.Group("/")
	api.GET("news", handlers.ListNews)
	api.GET("news/:id", handlers.GetNewsDetail)
	api.POST("comments", handlers.CreateComment)
	api.POST("news", handlers.CreateNews)
	api.POST("/censor", handlers.CensorText)

	return r
}
