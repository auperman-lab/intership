package usersvc

import (
	"gateway/pkg/config"
	"gateway/pkg/usersvc/routes"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config) *ServiceClient {
	svc := &ServiceClient{
		Client: InitUserServiceClient(c),
	}

	routes := r.Group("/user")
	routes.GET("/:id", svc.Get)
	routes.POST("/:id", svc.Delete)

	return svc
}

func (svc *ServiceClient) Get(ctx *gin.Context) {
	routes.Get(ctx, svc.Client)
}

func (svc *ServiceClient) Delete(ctx *gin.Context) {
	routes.Delete(ctx, svc.Client)
}
