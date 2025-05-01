package router

import (
	"pm_go_version/app/constant"
	"pm_go_version/app/domain/entity"
	"pm_go_version/app/pkg"
	"pm_go_version/config"

	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CORSMiddleware 处理跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的域名列表
		allowedOrigins := []string{
			"http://localhost:3000", // React开发环境
			//"http://localhost:5173",   // Vite开发环境
			//"https://your-domain.com", // 生产环境域名
		}

		origin := c.Request.Header.Get("Origin")
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, PATCH, GET, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24小时
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer pkg.PanicHandler(c)

		// 1. 首先检查Authorization header
		token := c.GetHeader("Authorization")
		if len(token) > 7 {
			claims, err := pkg.Verfiy(token[7:])
			if err == nil {
				log.Info("Token verified from Authorization header")
				c.Set("parse_id", claims.ID)
				c.Set("parse_username", claims.UserName)
				c.Set("parse_email", claims.Email)
				c.Next()
				return
			}
		}

		// 2. 如果没有Authorization header或验证失败，检查cookie
		cookieToken, err := c.Cookie("token")
		//log.Info("cookieToken is: ", cookieToken)
		if err == nil && cookieToken != "" {
			claims, err := pkg.Verfiy(cookieToken)
			if err == nil {
				log.Info("Token verified from cookie")
				c.Set("parse_id", claims.ID)
				c.Set("parse_username", claims.UserName)
				c.Set("parse_email", claims.Email)
				c.Next()
				return
			}
		}

		// 3. 如果两种方式都失败，返回未授权错误
		log.Error("Authentication failed: no valid token found in header or cookie")
		pkg.PanicException(constant.InvalidRequest)
		c.Abort()
	}
}

func CheckOwnership(paramsKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer pkg.PanicHandler(c)
		log.Info("CheckOwnership middleware called")
		userIdValue, _ := c.Get("parse_id")
		userId := int(userIdValue.(uint))
		if paramsKey == "workspace" {
			workspaceId, err := strconv.Atoi(c.Param("workspaceId"))
			if err != nil {
				log.Error("Invalid parameter value: ", err)
				pkg.PanicException(constant.InvalidRequest)
			}
			if !CheckWorkspaceMembership(userId, workspaceId) {
				log.Error("User is not member of the workspace")
				pkg.PanicException(constant.Unauthorized)
				c.Abort()
			}
		} else if paramsKey == "project" {
			workspaceId, err := strconv.Atoi(c.Param("workspaceId"))
			if err != nil {
				log.Error("Invalid parameter value: ", err)
				pkg.PanicException(constant.InvalidRequest)
			}
			projectId, err := strconv.Atoi(c.Param("projectId"))
			if err != nil {
				log.Error("Invalid parameter value: ", err)
				pkg.PanicException(constant.InvalidRequest)
			}
			if !CheckProjectMembership(userId, workspaceId, projectId) {
				log.Error("User is not member of the workspace")
				pkg.PanicException(constant.Unauthorized)
			}
		} else if paramsKey == "task" {
			taskId, err := strconv.Atoi(c.Param("taskId"))
			if err != nil {
				log.Error("Invalid parameter value: ", err)
				pkg.PanicException(constant.InvalidRequest)
			}
			if !CheckTaskMembership(userId, taskId) {
				log.Error("User is not member of the workspace")
				pkg.PanicException(constant.Unauthorized)
				c.Abort()
			}
			log.Info("User is member of the workspace")
		} else if paramsKey == "batchTask" {
			log.Info("User is member of the workspace")
		}
		c.Next()
	}
}

func CheckWorkspaceOwnership(userId int, workspaceId int) bool {
	db := config.GetGdb()
	var uw entity.UserWorkspace
	r1 := db.Where("user_id=? AND workspace_id=?", userId, workspaceId).First(&uw)
	if r1.Error != nil {
		log.Error("Got an error when check if user have workspaces. Error: ", r1.Error)
		return false
	}
	if uw.UserMember != "Owner" {
		log.Error("User is not owner to delete workspace. Error: ")
		return false
	} else {
		return true
	}
}

func CheckWorkspaceMembership(userId int, workspaceId int) bool {
	db := config.GetGdb()
	var userWorkspace entity.UserWorkspace
	r1 := db.Where("user_id=? AND workspace_id=?", userId, workspaceId).First(&userWorkspace)
	if r1.Error != nil {
		log.Error("Got an error when check if user belongs to workspace. Error: ", r1.Error)
		return false
	}
	return true
}

func CheckProjectMembership(userId int, workspaceId int, projectId int) bool {
	if !CheckWorkspaceMembership(userId, workspaceId) {
		log.Error("User is not member of the workspace")
		return false
	}
	db := config.GetGdb()
	var project entity.Project
	r1 := db.Where("id=? AND workspace_id=?", projectId, workspaceId).First(&project)
	if r1.Error != nil {
		log.Error("Got an error when check if workspace has this project. Error: ", r1.Error)
		return false
	}
	return true
}

func CheckTaskMembership(userId int, taskId int) bool {
	db := config.GetGdb()
	var result map[string]any
	subQuery := db.Model(&entity.Task{}).Where("id=?", taskId)
	r1 := db.Debug().Table(`(?) u`, subQuery).
		Joins("JOIN workspaces w On u.workspace_id = w.id").
		Joins("JOIN user_workspaces us on u.workspace_id=us.workspace_id").
		Where("us.user_id=?", userId).Find(&result)

	if r1.Error != nil {
		log.Error("Got an error when check if workspace has this task. Error: ", r1.Error)
		return false
	}

	if len(result) == 0 {
		log.Error("User is not member of the workspace for the specific task")
		return false
	}
	log.Info("User is member of the workspace for the specific task")
	return true
}

// var project constant.Project
