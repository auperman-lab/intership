package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"pkg/db/go/internal/controller"
	"pkg/db/go/internal/models"
	"pkg/db/go/internal/repository"
	service "pkg/db/go/internal/services"
	"pkg/db/go/pkg/db"
)

var DB *gorm.DB

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	fmt.Printf("POSTGRES_HOST: %s\n", os.Getenv("POSTGRES_HOST"))

	userRepo := repository.NewUserRepository(DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	userRouter := r.Group("/user")
	{
		userRouter.POST("/register", userController.Register)
		userRouter.POST("/login", userController.Login)
		userRouter.DELETE("/delete", userController.Delete)
		r.POST("/get", func(c *gin.Context) {

			var users []models.UserModel
			DB.Find(&users)

			c.JSON(http.StatusOK, gin.H{"data": users})
		})

	}

	_ = r.Run(":8888")

}

func init() {
	DB = db.Init()
	err := DB.AutoMigrate(models.UserModel{})
	if err != nil {
		log.Fatalf("failed to migrate user model\n")
	}
}
