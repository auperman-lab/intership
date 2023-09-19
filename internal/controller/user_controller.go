package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"intership/internal/models"
	"intership/internal/pb"
	"net/http"
)

type IUserService interface {
	Delete(user models.UserModel) error
	Get(userID primitive.ObjectID) (models.UserModel, error)
}

type UserController struct {
	pb.UnimplementedUserServiceServer
	userService IUserService
}

func NewUserController(userService IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (ctrl *UserController) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	var user models.UserModel
	user.ID, _ = primitive.ObjectIDFromHex(request.Id)

	err := ctrl.userService.Delete(user)
	if err != nil {
		return &pb.DeleteUserResponse{
			Status: http.StatusInternalServerError,
			Error:  "sorry cant delete",
		}, nil
	}

	return &pb.DeleteUserResponse{
		Status: http.StatusOK,
		Error:  "",
	}, nil
}

func (ctrl *UserController) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id, _ := primitive.ObjectIDFromHex(request.Id)

	user, err := ctrl.userService.Get(id)
	if err != nil {
		return &pb.GetUserResponse{
			Status: http.StatusInternalServerError,
			Error:  "cannot find id",
		}, nil
	}

	response := &pb.GetUserResponse{
		Id:       (user.ID).Hex(),
		Email:    user.Email,
		Password: user.Password,
		Status:   http.StatusOK,
		Error:    "",
	}

	return response, nil
}
