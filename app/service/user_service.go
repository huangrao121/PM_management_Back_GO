package service

import (
	"context"
	"net/http"
	"pm_go_version/app/constant"
	"pm_go_version/app/domain"
	"pm_go_version/app/domain/entity"
	"pm_go_version/app/pkg"
	"pm_go_version/app/pkg/cache"
	"pm_go_version/app/repository"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	GetMe(c *gin.Context)
	LoginUser(c *gin.Context)
	SignupUser(c *gin.Context)
	GetUserInfo(c *gin.Context)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
}

type UserServiceImpl struct {
	Ur    repository.UserRepository
	cache cache.Cache
}

func (usv *UserServiceImpl) GetMe(c *gin.Context) {
	defer pkg.PanicHandler(c)
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
		"message": "Verified",
	}))
}

func (usv *UserServiceImpl) GetUserInfo(c *gin.Context) {
	defer pkg.PanicHandler(c)

	log.Info("Start to execute get all user in process")

	// data, err := usv.Ur.GetUserInfo()
	// if err != nil {
	// 	log.Error("Happened error getting list of users. Error", err)
	// 	pkg.PanicException(constant.DataNotFound)
	// }
	userId, _ := c.Get("parse_id")
	userName, _ := c.Get("parse_username")
	email, _ := c.Get("parse_email")
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
		"id":       userId,
		"userName": userName,
		"email":    email,
	}))
}

func (usv *UserServiceImpl) LoginUser(c *gin.Context) {
	defer pkg.PanicHandler(c)

	log.Info("Start to executing user log in process")
	var request entity.User
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Error when try to logged in and convert the json format")
		pkg.PanicException(constant.InvalidRequest)
	}
	log.Info(request)
	data, err := usv.Ur.Login(&request)

	if err != nil {
		log.Error("Happened error when logged in to the database. Error", err)
		pkg.PanicException(constant.DataNotFound)
	}
	token, err := pkg.GenerateJWT(&data)
	if err != nil {
		log.Error("Happened error when generate jwt token from user info", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.SetCookie("token", token, (24 * 60 * 60), "/", "localhost", false, true)
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
		"userid":    data.ID,
		"user_name": data.UserName,
		"email":     data.Email,
		"token":     token,
	}))
}

func (usv *UserServiceImpl) SignupUser(c *gin.Context) {
	defer pkg.PanicHandler(c)

	log.Info("Start to executing sign up process")
	var request entity.User
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Error when try to sign up in the database and convert to json format, error is: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	data, err := usv.Ur.Save(&request)
	if err != nil {
		log.Error("Happened error when adding user to database. Error", err)
		pkg.PanicException(constant.UnknownError)
	}
	token, err := pkg.GenerateJWT(&data)
	if err != nil {
		log.Error("Happened error when generate jwt token from user info", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.SetCookie("token", token, (24 * 60 * 60), "/", "localhost", false, true)
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
		"userid":    data.ID,
		"user_name": data.UserName,
		"email":     data.Email,
		"token":     token,
	}))

}

func UserServiceInit(userRepository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		Ur:    userRepository,
		cache: cache.NewRedisCache(),
	}
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	// 尝试从缓存获取
	var user domain.User
	cacheKey := "user:" + id
	err := s.cache.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// 缓存未命中，从数据库获取
	user, err = s.Ur.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 将结果存入缓存，设置过期时间为1小时
	err = s.cache.Set(ctx, cacheKey, user, time.Hour)
	if err != nil {
		// 缓存错误不影响主流程
	}

	return &user, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	// 更新数据库
	err := s.Ur.Update(ctx, user)
	if err != nil {
		return err
	}

	// 更新缓存
	cacheKey := "user:" + user.ID
	err = s.cache.Delete(ctx, cacheKey)
	if err != nil {
		// 缓存错误不影响主流程
	}

	return nil
}
