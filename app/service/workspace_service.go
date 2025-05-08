package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"pm_go_version/app/constant"
	"pm_go_version/app/domain/entity"
	"pm_go_version/app/pkg"
	"pm_go_version/app/pkg/redis_config"
	"pm_go_version/app/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Code struct {
	InviteCode string `json:"invite_code"`
}

type WorkspaceService interface {
	GetListofWorkspaces(c *gin.Context)
	CreateWorkspace(c *gin.Context)
	GetWorkspacesById(c *gin.Context)
	DeleteWorkspaceById(c *gin.Context)
	UpdateWorkspaceById(c *gin.Context)
	GetSingleWorkspaceById(c *gin.Context)
	ResetInvite(c *gin.Context)
	JoinWorkspace(c *gin.Context)
	GetWorkspaceInfo(c *gin.Context)
}

type WorkspaceServiceImpl struct {
	Wr  repository.WorkspaceRepository
	Rdb *redis_config.RedisCache
}

func (ws *WorkspaceServiceImpl) GetListofWorkspaces(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var data, err = ws.Wr.GetWorkspaces()
	if err != nil {
		log.Error("Happened error getting list of users. Error", err)
		pkg.PanicException(constant.DataNotFound)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, data))
}

/*
service create workspace
*/
func (ws *WorkspaceServiceImpl) CreateWorkspace(c *gin.Context) {
	defer pkg.PanicHandler(c)

	value, isExist := c.Get("parse_id")
	createrName, isExist2 := c.Get("parse_username")
	if !isExist || !isExist2 {
		log.Error("Error Try to fetch user id from middleware, error: ")
		pkg.PanicException(constant.UnknownError)
	}
	userId, _ := ConvertAnyToInt(value)

	newFileName, save_err := SavetoLocalWithNewName(c, "workspace_image", "workspace_image")
	if save_err != nil {
		log.Error("Error when try to bind the form format to struct, error is: ", save_err)
		pkg.PanicException(constant.UnknownError)
	}

	var request entity.Workspace
	if err := c.ShouldBind(&request); err != nil {
		log.Error("Error when try to bind the form format to struct, error is: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	name, _ := ConvertAnyToString(createrName)
	request.CreaterID = userId
	request.CreaterName = name
	request.ImageUrl = newFileName
	request.InviteCode = pkg.GenerateInviteCode(10)
	id, _ := ws.Wr.CreateWorkspace(&request)

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
		"workspaceID": id,
	}))
}

func (ws *WorkspaceServiceImpl) GetWorkspacesById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	value, isExist := c.Get("parse_id")
	if !isExist {
		log.Error("Error Try to fetch user id from middleware, error: ")
		pkg.PanicException(constant.UnknownError)
	}

	userId, _ := ConvertAnyToInt(value)
	redisKey := "user:" + fmt.Sprintf("%v", value) + "workspace"

	var result interface{}
	data, err := ws.Rdb.GetStructValue(c, redisKey, func() (interface{}, error) {
		return ws.Wr.GetWorkspacesById(userId)
	})
	if err != nil {
		log.Error("Failed to get workspace data: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Error("Failed to unmarshal workspace data: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, result))
}

func (ws *WorkspaceServiceImpl) DeleteWorkspaceById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	workspace_id := c.Param("workspaceId")
	value, isExist := c.Get("parse_id")
	if !isExist {
		log.Error("Error Try to fetch user id from middleware, error: ")
		pkg.PanicException(constant.UnknownError)
	}
	//convert user id from middleware to int 64
	userId, _ := ConvertAnyToInt(value)

	//convert workspace id from string to uint
	num, err := strconv.ParseInt(workspace_id, 10, 64)
	if err != nil {
		log.Error("Error Try to parse the workspace id to integer, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	isSuccess, _ := ws.Wr.DeleteWorkspaceById(userId, uint(num))
	if !isSuccess {
		log.Error("Error try to delete workspace by id")
		pkg.PanicException(constant.InvalidRequest)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
		"workspaceID": num,
	}))
}

func (ws *WorkspaceServiceImpl) UpdateWorkspaceById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	workspace_id := c.Param("workspaceId")
	value, isExist := c.Get("parse_id")
	if !isExist {
		log.Error("Error Try to fetch user id from middleware, error: ")
		pkg.PanicException(constant.UnknownError)
	}
	//convert user id from middleware to uint
	userId, _ := ConvertAnyToInt(value)

	//Convert workspace id from string to uint
	workspaceId, err := strconv.Atoi(workspace_id)
	if err != nil {
		log.Error("Error Try to parse the workspace id to integer, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	var request entity.UpdateWorkspace
	if err := c.ShouldBind(&request); err != nil {
		log.Error("Error when try to bind the form format to struct, error is: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	//request.ID = uint(num)
	oldFileName, err2 := ws.Wr.GetImageName(userId, uint(workspaceId))
	if err2 != nil {
		log.Error("Failed to fetch the file name, error is: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	newFileName, err3 := SavetoLocalWithNewName(c, "workspace_image", "workspace_image")
	if err3 != nil {
		log.Error("Failed to save the updated image file to local disk, error is: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	request.ImageUrl = &newFileName
	log.Debug("request: ", request)
	//Query the update content to the database
	isSuccess, _ := ws.Wr.UpdateWorkspaceById(userId, workspaceId, &request)
	if !isSuccess {
		log.Error("Failed to save the updated image file to local disk, error is: ", err)
		pkg.PanicException(constant.UnknownError)
	} else {
		os.Remove("public/" + oldFileName)
		c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
			"workspaceID": workspaceId,
			"name":        request.Name,
			"image_url":   request.ImageUrl,
		}))
	}
}

func (ws *WorkspaceServiceImpl) GetSingleWorkspaceById(c *gin.Context) {
	defer pkg.PanicHandler(c)

	workspace_id := c.Param("workspaceId")
	value, isExist := c.Get("parse_id")

	if !isExist {
		log.Error("Error Try to fetch user id from middleware, error: ")
		pkg.PanicException(constant.UnknownError)
	}
	//convert user id from middleware to int 64
	userId, _ := ConvertAnyToInt(value)

	//convert workspace id from string to uint
	num, err := strconv.ParseInt(workspace_id, 10, 64)
	if err != nil {
		log.Error("Error Try to parse the workspace id to integer, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	redisKey := "user:" + fmt.Sprintf("%v", value) + "workspace:" + fmt.Sprintf("%v", workspace_id)
	data, err := ws.Rdb.GetStructValue(c, redisKey, func() (interface{}, error) {
		return ws.Wr.GetSingleWorkspaceById(userId, uint(num))
	})
	//data, err := ws.Wr.GetSingleWorkspaceById(userId, uint(num))
	if err != nil {
		log.Error("Failed to find the workspace, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	var result interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Error("Failed to unmarshal workspace data: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, result))
}

func (ws *WorkspaceServiceImpl) ResetInvite(c *gin.Context) {
	defer pkg.PanicHandler(c)

	workspace_id := c.Param("workspaceId")
	value, isExist := c.Get("parse_id")
	if !isExist {
		log.Error("Error Try to fetch user id from middleware, error: ")
		pkg.PanicException(constant.UnknownError)
	}
	//convert user id from middleware to int 64
	userId, _ := ConvertAnyToInt(value)

	//convert workspace id from string to uint
	num, err := strconv.ParseInt(workspace_id, 10, 64)
	if err != nil {
		log.Error("Error Try to parse the workspace id to integer, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	inviteCode := pkg.GenerateInviteCode(10)
	isUpdated, err2 := ws.Wr.ResetInvite(userId, uint(num), inviteCode)
	if !isUpdated {
		log.Error("Failed to find the workspace, error: ", err2)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
		"invite_code": inviteCode,
	}))
}

func (ws *WorkspaceServiceImpl) JoinWorkspace(c *gin.Context) {
	defer pkg.PanicHandler(c)

	// var code Code
	// if err := c.ShouldBind(&code); err != nil {
	// 	log.Error("Can't fetch the invite code, error: ", err)
	// 	pkg.PanicException(constant.UnknownError)
	// }
	inviteCode, err := c.GetRawData()
	if err != nil {
		log.Error("Can't fetch the invite code, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	inviteString := string(inviteCode)
	//log.Info("From workspace service joinworkspace, :", inviteString)
	workspace_id := c.Param("workspaceId")
	value, isExist := c.Get("parse_id")
	if !isExist {
		log.Error("Error Try to fetch user id from middleware, error: ")
		pkg.PanicException(constant.UnknownError)
	}
	//convert user id from middleware to int 64
	userId, _ := ConvertAnyToInt(value)

	//convert workspace id from string to uint
	num, err := strconv.ParseInt(workspace_id, 10, 64)
	if err != nil {
		log.Error("Error Try to parse the workspace id to integer, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	result, err := ws.Wr.JoinWorkspace(userId, uint(num), inviteString)
	if err != nil {
		log.Error("Failed to join the workspace, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, result))
}

func (ws *WorkspaceServiceImpl) GetWorkspaceInfo(c *gin.Context) {
	defer pkg.PanicHandler(c)

	workspace_value := c.Param("workspaceId")
	workspace_id, err := strconv.ParseInt(workspace_value, 10, 64)
	if err != nil {
		log.Error("Error Try to parse the workspace id to integer, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	result, err := ws.Wr.GetWorkspaceInfo(uint(workspace_id))
	if err != nil {
		log.Error("Failed to get the workspace info, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	// c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
	// 	"name": result.Name,
	// }))
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, pkg.BuildResponse(constant.Success,
		gin.H{"name": result.Name},
	)))
}

func ConvertAnyToInt(value any) (uint, bool) {
	userId, ok := value.(uint)
	if !ok {
		log.Error("Type conversion assertion error, parse id: ", userId)
		return 0, false
	}
	return uint(userId), true
}

func ConvertAnyToString(value any) (string, bool) {
	convertedValue, ok := value.(string)
	if !ok {
		log.Error("Type conversion assertion error, parse id: ", convertedValue)
		return "", false
	}
	return convertedValue, true
}

func WorkspaceServiceInit(wr repository.WorkspaceRepository) *WorkspaceServiceImpl {
	return &WorkspaceServiceImpl{
		Wr:  wr,
		Rdb: redis_config.GetRedisCache(),
	}
}
