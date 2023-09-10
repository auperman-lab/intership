package authsvc

import (
	"context"
	"gateway/pkg/authsvc/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthMiddlewareConfig struct {
	svc *ServiceClient
}

func (c *AuthMiddlewareConfig) AuthRequired(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	accesToken, err := ctx.Cookie("acces_token")

	res, err := c.svc.Client.Validate(context.Background(), &pb.ValidateUserRequest{
		AccesToken:   accesToken,
		RefreshToken: refreshToken,
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.SetCookie("acces_token", res.AccesToken, 60*5, "", "", false, true)
	ctx.SetCookie("refresh_token", res.RefreshToken, 60*15, "", "", false, true)
	ctx.SetCookie("id", string(res.UserId), 5*60, "", "", false, true)

}
