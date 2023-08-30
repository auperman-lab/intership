package main

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
	"pkg/db/go/internal/controller/grpc_transport"
	"pkg/db/go/internal/controller/http_transport"
	"pkg/db/go/internal/middleware"
	"pkg/db/go/internal/models"
	"pkg/db/go/internal/pb"
	"pkg/db/go/internal/repository"
	service "pkg/db/go/internal/services"
	"pkg/db/go/pkg/db"
)

func main() {

	userRepo := repository.NewUserRepository(db.DB)
	tokenRepo := repository.NewTokenRepository(db.RDB)

	userService := service.NewUserService(userRepo)
	tokenService := service.NewTokenService(tokenRepo)

	userController := grpc_transport.NewUserControllergrpc(userService, tokenService)

	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//myServer := &grpc_transport.GrpcStruct{}
	pb.RegisterUserControllerServer(s, userController)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func httpServer() {
	r := gin.Default()
	r.Use(gin.Recovery())

	userRepo := repository.NewUserRepository(db.DB)
	tokenRepo := repository.NewTokenRepository(db.RDB)

	userService := service.NewUserService(userRepo)
	tokenService := service.NewTokenService(tokenRepo)

	userController := http_transport.NewUserControllerhttp(userService, tokenService)

	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)
	r.DELETE("/delete", userController.Delete)
	r.GET("/refresh", middleware.ValidateToken, userController.Refresh)

	_ = r.Run(":8080")

}

func init() {
	db.DB = db.Init()
	err := db.DB.AutoMigrate(models.UserModel{})
	if err != nil {
		log.Fatalf("failed to migrate user model\n")
	}

	db.RDB, _ = db.InitRedis()
	if err != nil {
		log.Fatalf("failed to initialize Redis: %v\n", err)
	}

}
