package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"pkg/db/go/internal/models"
	"time"
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

	if c.Bind(&user) != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		return
	}

	err := ctrl.userService.Register(&user)
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (ctrl *UserController) Login(c *gin.Context) {

	var guest models.UserModel

	if c.Bind(&guest) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		return
	}
	err := ctrl.userService.Login(&guest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalit Credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  guest.ID,
		"exp": time.Now().Add(time.Hour * 8).Unix(),
	})

	secretKey := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*8, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})

}

func (ctrl *UserController) Delete(c *gin.Context) {
	var user models.UserModel

	if c.Bind(&user) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		return
	}

	err := ctrl.userService.Delete(&user)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid id"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

}

func (ctrl *UserController) Validate(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "Logged in"})

}
