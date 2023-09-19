package controller

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"intership/internal/pb"
	"net/http"
)

type IAdminService interface {
	GetUserID() ([]string, error)
	AssignAdmin(ID primitive.ObjectID) error
}

type AdminController struct {
	pb.UnimplementedAdminServiceServer
	adminService IAdminService
}

func NewAdminController(adminService IAdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

func (ctrl *AdminController) GetUsersID(ctx context.Context, request *pb.GetUsersIDsRequest) (*pb.GetUsersIDsResponse, error) {

	fmt.Println("hello get users")
	ids, err := ctrl.adminService.GetUserID()
	if err != nil {
		fmt.Println("bl error in controller", err.Error())
		response := &pb.GetUsersIDsResponse{
			Ids:    ids,
			Status: http.StatusInternalServerError,
			Err:    "some error" + err.Error(),
		}
		return response, err
	}

	fmt.Println("kzdm")

	return &pb.GetUsersIDsResponse{
		Ids:    ids,
		Status: http.StatusOK,
		Err:    "",
	}, nil
}

func (ctrl *AdminController) AssignAdmin(ctx context.Context, request *pb.AssignAdminRequest) (*pb.AssignAdminResponse, error) {

	id, _ := primitive.ObjectIDFromHex(request.Id)
	err := ctrl.adminService.AssignAdmin(id)
	if err != nil {
		fmt.Println(err)
		return &pb.AssignAdminResponse{
			Status: http.StatusInternalServerError,
			Err:    "sorry cant create" + err.Error(),
		}, nil
	}

	return &pb.AssignAdminResponse{
		Status: http.StatusOK,
		Err:    "",
	}, nil
}
