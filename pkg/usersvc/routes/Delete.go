package routes

import (
	"context"
	"gateway/pkg/usersvc/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Delete(ctx *gin.Context, c pb.UserServiceClient) {
	id := ctx.Param("id")

	res, err := c.DeleteUser(context.Background(), &pb.DeleteUserRequest{
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
