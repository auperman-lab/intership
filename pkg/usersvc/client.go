package usersvc

import (
	"gateway/pkg/config"
	"gateway/pkg/usersvc/pb"
	"google.golang.org/grpc"
	"log"
)

type ServiceClient struct {
	Client pb.UserServiceClient
}

func InitUserServiceClient(c *config.Config) pb.UserServiceClient {
	cc, err := grpc.Dial(c.UserSvcUrl, grpc.WithInsecure())

	if err != nil {
		log.Fatal("Could not connect:", err)
	}

	return pb.NewUserServiceClient(cc)
}
