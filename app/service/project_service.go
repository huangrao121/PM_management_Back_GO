package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pm_go_version/app/constant"
	"pm_go_version/app/domain/entity"
	"pm_go_version/app/pkg"
	"pm_go_version/app/pkg/redis_config"
	"pm_go_version/app/repository"
	"strconv"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

type ProjectService interface {
	GetListofProjects(c *gin.Context)
	CreateProject(c *gin.Context)
	GetProjectById(c *gin.Context)
	UpdateProjectById(c *gin.Context)
}

type ProjectServiceImpl struct {
	Pr  repository.ProjectRepository
	Rdb *redis_config.RedisCache
}

func (ps *ProjectServiceImpl) GetListofProjects(c *gin.Context) {
	//Check if the user is member of a workspace, then get list of projects of that workspace
	defer pkg.PanicHandler(c)
	user_id, workspace_id := GetUnWIds(c)

	redisKey := fmt.Sprintf("user:%v workspace:%v projects", user_id, workspace_id)
	data, err := ps.Rdb.GetStructValue(c, redisKey, func() (interface{}, error) {
		return ps.Pr.GetListofProjects(user_id, workspace_id)
	})
	//result, err := ps.Pr.GetListofProjects(user_id, workspace_id)
	if err != nil {
		log.Error("Failed to get all projects of a workspace, error is: ", err)
		pkg.PanicException(constant.DataNotFound)
	}

	var result interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Error("Failed to unmarshal workspace data: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, result))
}

func (ps *ProjectServiceImpl) CreateProject(c *gin.Context) {
	var project entity.Project
	if err := c.ShouldBind(&project); err != nil {
		log.Error("Failed to bind user request to project object, error is: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	value, isExist := c.Get("parse_id")
	if !isExist {
		log.Error("Error Try to fetch user id from middleware, error: ")
		pkg.PanicException(constant.UnknownError)
	}
	//log.Info("The data type is ", reflect.TypeOf(value))
	userId, _ := ConvertAnyToInt(value)

	newFileName, save_err := SavetoLocalWithNewName(c, "project_image", "project_image")
	if save_err != nil {
		log.Error("Error when try to bind the form format to struct, error is: ", save_err)
		pkg.PanicException(constant.UnknownError)
	}
	project.ImageUrl = newFileName
	//log.Info(err)

	isCreated, err := ps.Pr.CreateProject(userId, &project)
	if !isCreated {
		log.Error("Failed to create new project, error is: ", err)
		pkg.PanicException(constant.DataNotFound)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, project))
}

func (ps *ProjectServiceImpl) GetProjectById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	// 获取并验证参数
	projectId, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		log.Error("Invalid project ID format: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	workspaceId, err := strconv.ParseUint(c.Param("workspaceId"), 10, 32)
	if err != nil {
		log.Error("Invalid workspace ID format: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	// 获取并验证用户ID
	userIDValue, exists := c.Get("parse_id")
	if !exists {
		log.Error("User ID not found in context")
		pkg.PanicException(constant.UnknownError)
	}

	userId, isInt := ConvertAnyToInt(userIDValue)
	if !isInt {
		log.Error("Failed to convert user ID: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	result, err := ps.Pr.GetProjectById(userId, uint(workspaceId), uint(projectId))
	if err != nil {
		log.Error("Failed to get project by id: ", err)
		pkg.PanicException(constant.DataNotFound)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, result))
}

func (ps *ProjectServiceImpl) UpdateProjectById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	projectId, err := strconv.ParseUint(c.Param("projectId"), 10, 32)
	if err != nil {
		log.Error("Invalid project ID format: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	workspaceId, err := strconv.ParseUint(c.Param("workspaceId"), 10, 32)
	if err != nil {
		log.Error("Invalid workspace ID format: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	// 获取并验证用户ID
	userIDValue, exists := c.Get("parse_id")
	if !exists {
		log.Error("User ID not found in context")
		pkg.PanicException(constant.UnknownError)
	}

	userId, isInt := ConvertAnyToInt(userIDValue)
	if !isInt {
		log.Error("Failed to convert user ID: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	var project entity.Project
	if err := c.ShouldBind(&project); err != nil {
		log.Error("Failed to bind user request to project object, error is: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	isUpdated, err := ps.Pr.UpdateProjectById(userId, uint(workspaceId), uint(projectId), &project)
	if !isUpdated {
		log.Error("Failed to update project by id: ", err)
		pkg.PanicException(constant.DataNotFound)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, project))
}

func ProjectServiceInit(pr repository.ProjectRepository) *ProjectServiceImpl {
	return &ProjectServiceImpl{
		Pr:  pr,
		Rdb: redis_config.GetRedisCache(),
	}
}
