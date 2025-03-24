package service

import (
	"net/http"
	"pm_go_version/app/constant"
	"pm_go_version/app/pkg"
	"pm_go_version/app/repository"

	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserWorkspaceService interface {
	GetListofMembersByWorkspaceId(c *gin.Context)
	DeleteMemberByWorkspaceId(c *gin.Context)
}

type UserWorkspaceServiceImpl struct {
	Uwr repository.UserWorkspaceRepository
}

func (uws *UserWorkspaceServiceImpl) GetListofMembersByWorkspaceId(c *gin.Context) {
	defer pkg.PanicHandler(c)
	user_id, workspace_id := GetUnWIds(c)

	data, err := uws.Uwr.GetListofMembersByWorkspaceId(user_id, workspace_id)
	if err != nil {
		log.Error("Failed to fetch the member list from database error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, data))
}

func (uws *UserWorkspaceServiceImpl) DeleteMemberByWorkspaceId(c *gin.Context) {
	defer pkg.PanicHandler(c)
	user_id, workspace_id := GetUnWIds(c)

	memberValue, err := c.GetRawData()
	if err != nil {
		log.Error("Can't fetch the member id, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	memberId := string(memberValue)
	member_id, err := strconv.ParseUint(memberId, 10, 64)
	if err != nil {
		log.Error("member ID type is incorrect, must be number, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	member_id_uint := uint(member_id)
	if user_id == member_id_uint {
		log.Error("You can't delete yourself")
		pkg.PanicException(constant.UnknownError)
	}
	isDeleted, _ := uws.Uwr.DeleteMemberByWorkspaceId(user_id, workspace_id, member_id_uint)
	if !isDeleted {
		log.Error("Failed to delete member from workspace, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, gin.H{
		"workspace_id": workspace_id,
		"member_id":    memberId,
	}))
}

func UserWorkspaceServiceInit(uwr repository.UserWorkspaceRepository) *UserWorkspaceServiceImpl {
	return &UserWorkspaceServiceImpl{Uwr: uwr}
}

func GetUnWIds(c *gin.Context) (uint, uint) {
	workspaceId := c.Param("workspaceId")
	value, isExist := c.Get("parse_id")
	if !isExist {
		log.Error("Error Try to fetch user id from middleware, error: ")
		pkg.PanicException(constant.UnknownError)
	}
	workspace_id, err := strconv.ParseUint(workspaceId, 10, 64)
	if err != nil {
		log.Error("Error Try to parse the workspace id to integer, error: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	userId, _ := ConvertAnyToInt(value)
	return userId, uint(workspace_id)
}
