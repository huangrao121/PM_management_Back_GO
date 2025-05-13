package service

import (
	"encoding/json"
	"net/http"
	"os"
	"pm_go_version/app/constant"
	"pm_go_version/app/domain/entity"
	"pm_go_version/app/pkg"
	"pm_go_version/app/pkg/cache"
	"pm_go_version/app/repository"
	"pm_go_version/app/utils"

	"golang.org/x/oauth2"
	//"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	GetMe(c *gin.Context)
	LoginUser(c *gin.Context)
	SignupUser(c *gin.Context)
	GetUserInfo(c *gin.Context)
	GoogleLogin(c *gin.Context)
	GoogleCallback(c *gin.Context, code, state, cookieState string)
	// GithubLogin(c *gin.Context)
	// GithubCallback(c *gin.Context)
}

type UserServiceImpl struct {
	Ur           repository.UserRepository
	cache        cache.Cache
	GoogleConfig *oauth2.Config
	GithubConfig *oauth2.Config
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

func (usv *UserServiceImpl) GoogleLogin(c *gin.Context) {
	defer pkg.PanicHandler(c)
	log.Info("Start to executing google login process")
	state := utils.GenerateSecureState()
	c.SetCookie("oauth_state", state, 60*60, "/", "localhost", false, true)
	authUrl := usv.GoogleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, authUrl)
}

func (usv *UserServiceImpl) GoogleCallback(c *gin.Context, code, state, cookieState string) {
	log.Info("Start to executing google callback process")
	if state != cookieState {
		log.Error("Invalid state")
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=invalid_state")
		return
	}
	token, err := usv.GoogleConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		log.Error("Failed to exchange code for token", err)
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=failed_to_exchange_code_for_token")
		return
	}
	client := usv.GoogleConfig.Client(c.Request.Context(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		log.Error("Failed to get user info", err)
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=failed_to_get_user_info")
		return
	}
	defer response.Body.Close()

	var userInfo map[string]any
	json.NewDecoder(response.Body).Decode(&userInfo)
	user, getErr := usv.Ur.GetUserByEmail(userInfo["email"].(string))
	if getErr == nil {
		_, hasErr := usv.Ur.GetOAuthIdentity(userInfo["email"].(string), userInfo["sub"].(string), "google")
		if hasErr == nil {
			jwtToken, err := pkg.GenerateJWT(user)
			if err != nil {
				log.Error("Failed to generate JWT token", err)
				c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=failed_to_generate_JWT_token")
				return
			}
			c.SetCookie("token", jwtToken, (24 * 60 * 60), "/", "localhost", false, true)
			c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/")
			return
		}

		// Has user but not oauth identity
		tempToken, err := pkg.GenerateTempJWT(userInfo["email"].(string), userInfo["name"].(string), userInfo["sub"].(string), "google")
		if err != nil {
			log.Error("Failed to generate temp JWT token", err)
			c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=failed_to_generate_temp_JWT_token")
			return
		}
		c.SetCookie("temp_token", tempToken, (30 * 60), "/", "localhost", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/oauth/binding")
		return
	}

	user1, err := usv.Ur.CreateUserOAuthIdentity(userInfo["email"].(string), userInfo["name"].(string), userInfo["sub"].(string), "google")
	if err != nil {
		log.Error("Failed to create user oauth identity", err)
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=failed_to_create_user_and_oauth_identity")
		return
	}
	jwtToken, err := pkg.GenerateJWT(user1)
	if err != nil {
		log.Error("Failed to generate JWT token", err)
		c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/login?error=failed_to_generate_JWT_token")
		return
	}
	c.SetCookie("token", jwtToken, (24 * 60 * 60), "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/")
}

func UserServiceInit(userRepository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		Ur:    userRepository,
		cache: cache.NewRedisCache(),
		GoogleConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		GithubConfig: &oauth2.Config{
			// ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			// ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			// RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
			// Scopes:       []string{"user:email"},
			// Endpoint:     github.Endpoint,
		},
	}
}

// func (s *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
// 	// 尝试从缓存获取
// 	var user domain.User
// 	cacheKey := "user:" + id
// 	err := s.cache.Get(ctx, cacheKey, &user)
// 	if err == nil {
// 		return &user, nil
// 	}

// 	// 缓存未命中，从数据库获取
// 	user, err = s.Ur.FindByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 将结果存入缓存，设置过期时间为1小时
// 	err = s.cache.Set(ctx, cacheKey, user, time.Hour)
// 	if err != nil {
// 		// 缓存错误不影响主流程
// 	}

// 	return &user, nil
// }

// func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
// 	// 更新数据库
// 	err := s.Ur.Update(ctx, user)
// 	if err != nil {
// 		return err
// 	}

// 	// 更新缓存
// 	cacheKey := "user:" + user.ID
// 	err = s.cache.Delete(ctx, cacheKey)
// 	if err != nil {
// 		// 缓存错误不影响主流程
// 	}

// 	return nil
// }
