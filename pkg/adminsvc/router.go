package adminsvc

import (
	"fmt"
	"gateway/pkg/adminsvc/routes"
	"gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config) *ServiceClient {
	svc := &ServiceClient{
		Client: InitAdminServiceClient(c),
	}
	fmt.Println("oi bl router")

	route := r.Group("/admin")
	route.POST("/add/:id", svc.AddAdmin)
	route.GET("/get/", svc.GetUsers)

	return svc
}

func (svc *ServiceClient) AddAdmin(ctx *gin.Context) {

	routes.AddAdmin(ctx, svc.Client)
}

func (svc *ServiceClient) GetUsers(ctx *gin.Context) {

	routes.GetUser(ctx, svc.Client)
}
