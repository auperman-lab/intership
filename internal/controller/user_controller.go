package controller

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"intership/internal/models"
	"intership/internal/pb"
	"intership/pkg/db"
	"net/http"
)

type IUserService interface {
	Register(user *models.UserModel) error
	Login(user *models.UserModel) (primitive.ObjectID, error)
}
type ITokenService interface {
	GenerateTokenPair(ID primitive.ObjectID) (string, string, error)
	Validate(refreshToken string) (string, string, error)
}

type UserController struct {
	pb.UnimplementedAuthServiceServer
	userService  IUserService
	tokenService ITokenService
}

func NewUserController(userService IUserService, tokenService ITokenService) *UserController {
	return &UserController{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (ctrl *UserController) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &models.UserModel{
		Email:    request.Email,
		Password: request.Password,
	}

	err := ctrl.userService.Register(user)
	if err != nil {
		response := &pb.CreateUserResponse{
			Status: http.StatusInternalServerError,
			Error:  "Error creating user" + err.Error(),
		}

		return response, status.Errorf(codes.Internal, "Error creating user: %v", err)
	}

	response := &pb.CreateUserResponse{
		Status: http.StatusOK,
		Error:  "",
	}

	return response, nil
}

func (ctrl *UserController) LoginUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	var guest models.UserModel

	guest.Email = request.Email
	guest.Password = request.Password

	id, err := ctrl.userService.Login(&guest)
	guest.ID = id

	if err != nil {

		response := &pb.LoginUserResponse{
			Status:       http.StatusUnauthorized,
			Error:        "invalit credentials" + err.Error(),
			RefreshToken: "",
		}

		return response, status.Errorf(codes.InvalidArgument, "Invalid Credentials")
	}

	accesTokenString, refreshTokenString, err := ctrl.tokenService.GenerateTokenPair(guest.ID)
	if err != nil {
		response := &pb.LoginUserResponse{
			Status:       http.StatusUnauthorized,
			Error:        "cannot create tokens" + err.Error(),
			RefreshToken: refreshTokenString,
		}

		return response, status.Errorf(codes.Internal, "Cannot create tokens")
	}

	ctb := context.Background()
	db.RDB.Set(ctb, refreshTokenString, guest.ID, 60*15)

	response := &pb.LoginUserResponse{
		Status:       http.StatusOK,
		Error:        "",
		RefreshToken: refreshTokenString,
	}

	md := metadata.Pairs("authorization", "Bearer "+accesTokenString)
	grpc.SendHeader(ctx, md)

	return response, nil
}

func (ctrl *UserController) Validate(ctx context.Context, request *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	refreshtoken := request.RefreshToken

	newAccesToken, newRefreshToken, err := ctrl.tokenService.Validate(refreshtoken)

	if err != nil || newRefreshToken == "" {
		response := &pb.ValidateUserResponse{
			Status:       http.StatusUnauthorized,
			Error:        "cannot generate tokens",
			RefreshToken: newRefreshToken,
		}
		return response, err
	}
	response := &pb.ValidateUserResponse{
		Status:       http.StatusOK,
		Error:        "",
		RefreshToken: newRefreshToken,
	}

	md := metadata.Pairs("authorization", "Bearer "+newAccesToken)
	grpc.SendHeader(ctx, md)

	return response, nil

}
