package routes

import (
	"context"
	"gateway/pkg/usersvc/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Get(ctx *gin.Context, c pb.UserServiceClient) {
	id := ctx.Param("id")

	res, err := c.GetUser(context.Background(), &pb.GetUserRequest{
		Id: id,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "sorry?"})
		return
	}

	if res.Status != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "sorry?"})
		return
	}

	ctx.JSON(http.StatusOK, &res)

}
