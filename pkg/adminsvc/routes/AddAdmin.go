package routes

import (
	"context"
	"fmt"
	"gateway/pkg/adminsvc/pb"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddAdmin(ctx *gin.Context, c pb.AdminServiceClient) {
	id := ctx.Param("id")

	fmt.Println("hello addadmin")

	res, err := c.AssignAdmin(context.Background(), &pb.AssignAdminRequest{
		Id: id,
	})

	fmt.Println("this is res", res)

	if err != nil || res.Status != http.StatusOK {

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
