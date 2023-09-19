package adminsvc

import (
	"gateway/pkg/adminsvc/pb"
	"gateway/pkg/config"

	"google.golang.org/grpc"
	"log"
)

type ServiceClient struct {
	Client pb.AdminServiceClient
}

func InitAdminServiceClient(c *config.Config) pb.AdminServiceClient {
	cc, err := grpc.Dial(c.AdminSvcUrl, grpc.WithInsecure())

	if err != nil {
		log.Fatal("Could not connect:", err)
	}

	return pb.NewAdminServiceClient(cc)
}
