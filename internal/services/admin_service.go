package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAdminRepository interface {
	GetUsersID() ([]string, error)
	AssignAdmin(id primitive.ObjectID) error
}

type AdminService struct {
	userRepo IAdminRepository
}

func NewAdminService(userRepo IAdminRepository) *AdminService {
	return &AdminService{
		userRepo: userRepo,
	}
}

func (svc *AdminService) GetUserID() ([]string, error) {
	return svc.userRepo.GetUsersID()
}

func (svc *AdminService) AssignAdmin(ID primitive.ObjectID) error {
	return svc.userRepo.AssignAdmin(ID)
}
