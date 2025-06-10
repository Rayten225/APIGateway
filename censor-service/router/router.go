package router

import (
	"github.com/gin-gonic/gin"
	"micronews/censor-service/handlers"
	"micronews/censor-service/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.POST("/censor", handlers.Censor)
	return r
}
