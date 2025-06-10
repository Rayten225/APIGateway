package router

import (
	"github.com/gin-gonic/gin"
	"micronews/news-service/handlers"
	"micronews/news-service/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.GET("/news", handlers.ListNews)
	r.GET("/news/:id", handlers.GetNews)
	r.POST("/news", handlers.CreateNews)
	return r
}
