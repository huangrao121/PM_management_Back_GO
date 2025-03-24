package controller

import (
	"pm_go_version/app/service"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	LoginUser(*gin.Context)
	SignupUser(*gin.Context)
	GetUserInfo(*gin.Context)
}

type UserControllerImpl struct {
	Svc service.UserService
}

func (u *UserControllerImpl) LoginUser(c *gin.Context) {
	u.Svc.LoginUser(c)
}

func (u *UserControllerImpl) SignupUser(c *gin.Context) {
	u.Svc.SignupUser(c)
}

func (u *UserControllerImpl) GetUserInfo(c *gin.Context) {
	u.Svc.GetUserInfo(c)
}

func UserControllerInit(svc service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		Svc: svc,
	}
}
