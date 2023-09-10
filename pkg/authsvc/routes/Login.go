package routes

import (
	"context"
	"gateway/pkg/authsvc/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context, c pb.UserControllerClient) {
	b := LoginRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.LoginUser(context.Background(), &pb.LoginUserRequest{
		Email:    b.Email,
		Password: b.Password,
	})

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Set cookies and return the response without sensitive tokens
	refreshToken := res.RefreshToken
	accessToken := res.AccesToken

	ctx.SetCookie("acces_token", accessToken, 60*5, "", "", false, true)
	ctx.SetCookie("refresh_token", refreshToken, 60*15, "", "", false, true)
	res.RefreshToken = ""
	res.AccesToken = ""

	ctx.JSON(http.StatusCreated, res)
}
