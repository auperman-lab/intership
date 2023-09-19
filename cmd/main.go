package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"intership/internal/controller"
	"intership/internal/pb"
	"intership/internal/repository"
	"intership/internal/services"
	"intership/pkg/db"
	"log"
	"net"
)

func main() {

	clientOptions := options.Client().ApplyURI("mongodb://danu:macrii@localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	userRepo := repository.NewUserRepository(client, "mongo_db", "users")
	tokenRepo := repository.NewTokenRepository(db.RDB)

	userService := services.NewUserService(userRepo)
	tokenService := services.NewTokenService(tokenRepo)

	userController := controller.NewUserController(userService, tokenService)

	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, userController)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func init() {
	db.DB = db.InitPostgres()
	//err := db.DB.AutoMigrate(models.UserModel{})
	//if err != nil {
	//	log.Fatalf("failed to migrate user model\n")
	//}

	db.RDB, _ = db.InitRedis()

	db.InitMongo()

}
