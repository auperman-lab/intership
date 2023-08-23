package service

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"pkg/db/go/internal/models"
)

type IUserRepository interface {
	Insert(user *models.UserModel) error
	Delete(user *models.UserModel) error
	CheckUserExists(userName string) (string, error)
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
		log.Print("Error hashing password")
		return err
	}
	user.Password = hashedPassword

	return svc.userRepo.Insert(user)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (svc *UserService) Login(guest *models.UserModel) error {

	pass, err := svc.userRepo.CheckUserExists(guest.Name)

	if err != nil {
		log.Print("Invalit credentials")
		return err
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(pass), []byte(guest.Password))

	if passErr != nil {
		log.Print("Invalit credentials")
		return err
	}

	return nil
}

func (svc *UserService) Delete(user *models.UserModel) error {

	return svc.Delete(user)
}
