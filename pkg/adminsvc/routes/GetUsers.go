package routes

import (
	"context"
	"fmt"
	"gateway/pkg/adminsvc/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(ctx *gin.Context, c pb.AdminServiceClient) {

	res, err := c.GetUsersID(context.Background(), &pb.GetUsersIDsRequest{})

	fmt.Println("this is res", res)
	fmt.Println("this is err", err)

	if err != nil || res.Err != "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ""})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
