package controller

type IUserService interface {
	GetByName(name string) string
}

type UserController struct {
	userSvc IUserService
}

func (ctrl *UserController) GetByName() {

}
