package http_transport

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"pkg/db/go/internal/models"
	"pkg/db/go/pkg/db"
)

type IUserService interface {
	Register(user *models.UserModel) error
	Login(user *models.UserModel) (uint, error)
	Delete(user *models.UserModel) error
}
type ITokenService interface {
	GenerateTokenPair(ID uint) (string, string, error)
}

type UserController struct {
	userService  IUserService
	tokenService ITokenService
}

func NewUserController(userService IUserService, tokenService ITokenService) *UserController {
	return &UserController{
		userService:  userService,
		tokenService: tokenService,
	}
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

	id, err := ctrl.userService.Login(&guest)
	guest.ID = id

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalit Credentials"})
	}

	accesTokenString, refreshTokenString, err := ctrl.tokenService.GenerateTokenPair(guest.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot create tokens"})
		return
	}
	ctx := context.Background()

	db.RDB.Set(ctx, refreshTokenString, guest.ID, 60*15)

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("acces_token", accesTokenString, 60*5, "", "", false, true)
	c.SetCookie("refresh_token", refreshTokenString, 60*15, "", "", false, true)
	c.SetCookie("logged", "true", 0, "", "", false, true)

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

func (ctrl *UserController) Refresh(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "good"})

}
