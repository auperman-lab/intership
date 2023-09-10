package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"pkg/db/go/internal/controller/amq_transport"
	"pkg/db/go/internal/controller/chat"
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
	grpcServer()
}

func wsServer() {
	hub := chat.NewHub()
	go chat.Run(hub)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.HandleWebSocket(hub, w, r)
	})

	log.Println("WebSocket server is running on :8080")
	http.ListenAndServe(":8080", nil)

}

func rbmqServer() {

	userRepo := repository.NewUserRepository(db.DB)
	tokenRepo := repository.NewTokenRepository(db.RDB)

	userService := service.NewUserService(userRepo)
	tokenService := service.NewTokenService(tokenRepo)

	userController := amq_transport.NewUserControlleramqp(userService, tokenService)

	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	queueName := "registration_queue"
	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	message1 := models.UserModel{
		Email:    "danone",
		Password: "damedame",
	}
	message1.ID = 10

	PublishMessage(ch, queueName, message1)
	//userController.Register(ch, queueName)
	userController.Delete(ch, queueName)
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

func grpcServer() {
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
	pb.RegisterUserControllerServer(s, userController)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func init() {
	db.DB = db.InitPostgres()
	err := db.DB.AutoMigrate(models.UserModel{})
	if err != nil {
		log.Fatalf("failed to migrate user model\n")
	}

	db.RDB, _ = db.InitRedis()
	if err != nil {
		log.Fatalf("failed to initialize Redis: %v\n", err)
	}

	db.InitMongo()

}

func PublishMessage(ch *amqp.Channel, queueName string, messageBody models.UserModel) error {
	body, _ := json.Marshal(messageBody)
	err := ch.Publish(
		"",        // Exchange
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	fmt.Printf("Message sent: %s\n", body)
	return nil
}
