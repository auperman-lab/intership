package amq_transport

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"pkg/db/go/internal/models"
)

type IUserService interface {
	Register(user *models.UserModel) error
	Login(user *models.UserModel) (uint, error)
	Delete(user *models.UserModel) error
	Get(ID uint) (models.UserModel, error)
}
type ITokenService interface {
	GenerateTokenPair(ID uint) (string, string, error)
}

type UserControlleramqp struct {
	userService  IUserService
	tokenService ITokenService
}

func NewUserControlleramqp(userService IUserService, tokenService ITokenService) *UserControlleramqp {
	return &UserControlleramqp{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (ctrl *UserControlleramqp) Register(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for msg := range msgs {
		var user models.UserModel

		if err := json.Unmarshal(msg.Body, &user); err != nil {
			fmt.Printf("Error unmarshalling message: %v", err)
			continue
		}

		if err := ctrl.userService.Register(&user); err != nil {
			fmt.Printf("User registration failed: %v", err)
			continue
		}

		fmt.Printf("User registered successfully: %+v\n", &user)
	}
}

func (ctrl *UserControlleramqp) Delete(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	for msg := range msgs {
		var user models.UserModel

		if err := json.Unmarshal(msg.Body, &user); err != nil {
			fmt.Printf("Error unmarshalling message: %v", err)
			continue
		}

		if err := ctrl.userService.Delete(&user); err != nil {
			fmt.Printf("User registration failed: %v", err)
			continue
		}

		fmt.Printf("User deleted successfully: %+v\n", &user)
	}
}
