package routes

import (
	"context"
	"gateway/pkg/authsvc/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Create(ctx *gin.Context, c pb.UserControllerClient) {
	b := RegisterRequestBody{}

	if err := ctx.BindJSON(&b); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	res, err := c.CreateUser(context.Background(), &pb.CreateUserRequest{
		Email:    b.Email,
		Password: b.Password,
	})

	if err != nil || res.Status != http.StatusOK {
		ctx.AbortWithError(http.StatusBadGateway, err)
		return
	}

	ctx.JSON(http.StatusCreated, &res)
}
