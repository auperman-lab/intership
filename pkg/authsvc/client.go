package authsvc

import (
	"gateway/pkg/authsvc/pb"
	"gateway/pkg/config"
	"google.golang.org/grpc"
	"log"
)

type ServiceClient struct {
	Client pb.UserControllerClient
}

func InitUserControllerClient(c *config.Config) pb.UserControllerClient {
	cc, err := grpc.Dial(c.AuthSvcUrl, grpc.WithInsecure())

	if err != nil {
		log.Fatal("Could not connect:", err)
	}

	return pb.NewUserControllerClient(cc)
}
