package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"micronews/censor-service/router"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter()
	if err := r.SetTrustedProxies([]string{"172.22.0.0/16"}); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}
	log.Println("[*] Censor Service listening on :8003")
	if err := r.Run(":8003"); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
