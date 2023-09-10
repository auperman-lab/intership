package authsvc

import (
	"gateway/pkg/authsvc/routes"
	"gateway/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *config.Config) *ServiceClient {
	svc := &ServiceClient{
		Client: InitUserControllerClient(c),
	}

	routes := r.Group("/auth")
	routes.POST("/register", svc.Create)
	routes.POST("/login", svc.Login)

	return svc
}

func (svc *ServiceClient) Create(ctx *gin.Context) {
	routes.Create(ctx, svc.Client)
}

func (svc *ServiceClient) Login(ctx *gin.Context) {
	routes.Login(ctx, svc.Client)
}
