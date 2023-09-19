package routes

import (
	"context"
	"gateway/pkg/authsvc/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Validate(ctx *gin.Context, c pb.AuthServiceClient) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	res, err := c.Validate(context.Background(), &pb.ValidateUserRequest{
		RefreshToken: refreshToken,
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, res)
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.SetCookie("refresh_token", res.RefreshToken, 60*15, "", "", false, true)

}
