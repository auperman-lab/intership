package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pkg/db/go/internal/models"
	"pkg/db/go/pkg/db"
)

func main() {
	r := gin.Default()

	dbInstance := db.Init()
	r.GET("/", func(c *gin.Context) {

		var users []models.UserModel
		dbInstance.Find(&users)

		c.JSON(http.StatusOK, gin.H{"data": users})
	})

	r.Run()

}
