package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"pkg/db/go/internal/models"
)

type IUserService interface {
	Register(user *models.UserModel) error
	Login(user *models.UserModel) error
	Delete(user *models.UserModel) error
}

type UserController struct {
	userService IUserService
}

func NewUserController(userService IUserService) *UserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) Register(c *gin.Context) {
	var user models.UserModel

	if c.Bind(user) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		return
	}

	err := ctrl.userService.Register(&user)
	if err.Error != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (ctrl *UserController) Login(c *gin.Context) {

	var guest models.UserModel

	if c.Bind(guest) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		return
	}
	err := ctrl.userService.Login

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalit Credentials"})
	}

	token := jwt.New(jwt.SigningMethodEdDSA)

	c.JSON(http.StatusCreated, gin.H{"token": token})

}

func (ctrl *UserController) Delete(c *gin.Context) {
	var user models.UserModel

	if c.Bind(user) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		return
	}

	err := ctrl.userService.Delete

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid id"})
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "User deleted successfully"})

}
