package grpc_transport

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"pkg/db/go/internal/models"
	"pkg/db/go/internal/pb"
	"pkg/db/go/pkg/db"
)

type IUserService interface {
	Register(user *models.UserModel) error
	Login(user *models.UserModel) (uint, error)
	Delete(user *models.UserModel) error
	Get(ID uint) (models.UserModel, error)
}
type ITokenService interface {
	GenerateTokenPair(ID uint) (string, string, error)
	Validate(accesToken string, refreshToken string) (string, string, uint)
}

type UserControllergrpc struct {
	pb.UnimplementedUserControllerServer
	userService  IUserService
	tokenService ITokenService
}

func NewUserControllergrpc(userService IUserService, tokenService ITokenService) *UserControllergrpc {
	return &UserControllergrpc{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (ctrl *UserControllergrpc) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &models.UserModel{
		Email:    request.Email,
		Password: request.Password,
	}

	err := ctrl.userService.Register(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating user: %v", err)
	}

	response := &pb.CreateUserResponse{
		Status: 200,
		Error:  "",
	}

	return response, nil
}

func (ctrl *UserControllergrpc) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	var guest models.UserModel

	guest.Email = request.Email
	guest.Password = request.Password

	id, err := ctrl.userService.Login(&guest)
	guest.ID = id

	if err != nil {

		response := &pb.LoginUserResponse{
			Status:       http.StatusUnauthorized,
			Error:        "invalit credentials",
			AccesToken:   "",
			RefreshToken: "",
		}

		return response, status.Errorf(codes.InvalidArgument, "Invalid Credentials")
	}

	accesTokenString, refreshTokenString, err := ctrl.tokenService.GenerateTokenPair(guest.ID)
	if err != nil {
		response := &pb.LoginUserResponse{
			Status:       http.StatusUnauthorized,
			Error:        "cannot create tokens",
			AccesToken:   accesTokenString,
			RefreshToken: refreshTokenString,
		}

		return response, status.Errorf(codes.Internal, "Cannot create tokens")
	}

	ctb := context.Background()
	db.RDB.Set(ctb, refreshTokenString, guest.ID, 60*15)

	response := &pb.LoginUserResponse{
		Status:       http.StatusOK,
		Error:        "",
		AccesToken:   accesTokenString,
		RefreshToken: refreshTokenString,
	}

	return response, nil
}

func (ctrl *UserControllergrpc) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	var user models.UserModel
	user.ID = uint(request.Id)

	err := ctrl.userService.Delete(&user)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Invalid id")
	}

	return &pb.DeleteUserResponse{}, nil
}

func (ctrl *UserControllergrpc) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := uint(request.Id)

	user, err := ctrl.userService.Get(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	response := &pb.GetUserResponse{
		User: &pb.User{
			Id:        int32(user.ID),
			Email:     user.Email,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
			DeletedAt: nil,
		},
	}

	return response, nil
}

func (ctrl *UserControllergrpc) Validate(ctx context.Context, request *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	accestoken := request.AccesToken
	refreshtoken := request.RefreshToken

	newAccesToken, newRefreshToken, id := ctrl.tokenService.Validate(accestoken, refreshtoken)

	if newAccesToken == "" || newRefreshToken == "" || id == 0 {
		response := &pb.ValidateUserResponse{
			Status:       http.StatusUnauthorized,
			Error:        "cannot generate tokens",
			UserId:       int64(id),
			AccesToken:   newAccesToken,
			RefreshToken: newRefreshToken,
		}
		return response, fmt.Errorf("errror creating tokens or invalid id")
	}
	response := &pb.ValidateUserResponse{
		Status:       http.StatusOK,
		Error:        "",
		UserId:       int64(id),
		AccesToken:   newAccesToken,
		RefreshToken: newRefreshToken,
	}
	return response, nil

}
