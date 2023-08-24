package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"pkg/db/go/internal/controller"
	"pkg/db/go/internal/models"
	"pkg/db/go/internal/repository"
	service "pkg/db/go/internal/services"
	"pkg/db/go/pkg/db"
)

var DB *gorm.DB

func main() {
	r := gin.Default()
	//r.Use(gin.Recovery())

	userRepo := repository.NewUserRepository(DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	r.DELETE("/delete", userController.Delete)
	r.GET("/validate", userController.Validate)

	_ = r.Run(":8080")

}

func init() {
	DB = db.Init()
	err := DB.AutoMigrate(models.UserModel{})
	if err != nil {
		log.Fatalf("failed to migrate user model\n")
	}
}
