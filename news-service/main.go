package main

import (
	"log"
	"micronews/news-service/router"
)

func main() {
	r := router.SetupRouter()
	log.Println("[*] News Service listening on :8005")
	r.Run(":8005")
}
