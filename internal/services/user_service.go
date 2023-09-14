package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"intership/internal/models"
)

type IUserRepository interface {
	Delete(userID primitive.ObjectID) error
	CheckUserId(userID primitive.ObjectID) (models.UserModel, error)
}

type UserService struct {
	userRepo IUserRepository
}

func NewUserService(userRepo IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) Delete(user models.UserModel) error {

	return svc.userRepo.Delete(user.ID)
}

func (svc *UserService) Get(ID primitive.ObjectID) (models.UserModel, error) {

	return svc.userRepo.CheckUserId(ID)
}
