package router

import (
	"github.com/gin-gonic/gin"
	"micronews/comment-service/handlers"
	"micronews/comment-service/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.POST("/comments", handlers.CreateComment)
	r.GET("/comments", handlers.ListComments)
	return r
}
