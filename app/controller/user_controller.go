package controller

import (
	"net/http"
	"pm_go_version/app/service"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetMe(*gin.Context)
	LoginUser(*gin.Context)
	SignupUser(*gin.Context)
	GetUserInfo(*gin.Context)
	GoogleLogin(*gin.Context)
	GoogleCallback(*gin.Context)
	GithubLogin(*gin.Context)
	GithubCallback(*gin.Context)
}

type UserControllerImpl struct {
	Svc service.UserService
}

func (u *UserControllerImpl) GetMe(c *gin.Context) {
	u.Svc.GetMe(c)
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

func (u *UserControllerImpl) GoogleLogin(c *gin.Context) {
	u.Svc.GoogleLogin(c)
}

func (u *UserControllerImpl) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	cookieState, err := c.Cookie("oauth_state")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid state"})
		return
	}

	u.Svc.GoogleCallback(c, code, state, cookieState)
}

func (u *UserControllerImpl) GithubLogin(c *gin.Context) {
	// u.Svc.GithubLogin(c)
}

func (u *UserControllerImpl) GithubCallback(c *gin.Context) {
	// u.Svc.GithubCallback(c)
}

func UserControllerInit(svc service.UserService) *UserControllerImpl {
	return &UserControllerImpl{
		Svc: svc,
	}
}
