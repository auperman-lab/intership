package main

import (
	"gateway/pkg/adminsvc"
	"gateway/pkg/authsvc"
	"gateway/pkg/config"
	"gateway/pkg/usersvc"
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
	usersvc.RegisterRoutes(r, &c)
	adminsvc.RegisterRoutes(r, &c)

	r.Run(":8080")

}
