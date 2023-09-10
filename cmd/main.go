package main

import (
	"gateway/pkg/authsvc"
	"gateway/pkg/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	r := gin.Default()

	authsvc.RegisterRoutes(r, &c)

	r.Run(":8080")

}
