package services

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"intership/internal/models"
)

type IUserRepository interface {
	Insert(user *models.UserModel) error
	CheckUserEmail(userEmail string) (*models.UserModel, error)
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) Register(user *models.UserModel) error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		fmt.Println("Error hashing password")
		return err
	}
	user.Password = hashedPassword

	return svc.userRepo.Insert(user)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (svc *UserService) Login(guest *models.UserModel) (primitive.ObjectID, error) {
	user, err := svc.userRepo.CheckUserEmail(guest.Email)
	if err != nil {
		fmt.Println("Invalid credentials")
		return primitive.NilObjectID, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(guest.Password))
	if err != nil {
		fmt.Println("Invalid credentials")
		return primitive.NilObjectID, err
	}

	return user.ID, nil
}
