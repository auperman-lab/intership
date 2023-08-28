package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"pkg/db/go/internal/models"
)

type IUserRepository interface {
	Insert(user *models.UserModel) error
	Delete(user *models.UserModel) error
	CheckUserExists(userEmail string) (string, uint, error)
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

func (svc *UserService) Login(guest *models.UserModel) (uint, error) {

	passFromDB, idFromDB, err := svc.userRepo.CheckUserExists(guest.Email)

	if err != nil {
		fmt.Println("Invalit credentials")
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passFromDB), []byte(guest.Password))

	if err != nil {
		fmt.Println("Invalit credentials")
		return 0, err
	}

	return idFromDB, nil
}

func (svc *UserService) Delete(user *models.UserModel) error {

	return svc.userRepo.Delete(user)
}
