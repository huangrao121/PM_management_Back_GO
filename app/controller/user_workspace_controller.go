package controller

import (
	"pm_go_version/app/service"

	"github.com/gin-gonic/gin"
)

type UserWorkspaceController interface {
	GetListofMembersByWorkspaceId(c *gin.Context)
	DeleteMemberByWorkspaceId(c *gin.Context)
}

type UserWorkspaceControllerImpl struct {
	Uws service.UserWorkspaceService
}

func (uwc *UserWorkspaceControllerImpl) GetListofMembersByWorkspaceId(c *gin.Context) {
	uwc.Uws.GetListofMembersByWorkspaceId(c)
}

func (uwc *UserWorkspaceControllerImpl) DeleteMemberByWorkspaceId(c *gin.Context) {
	uwc.Uws.DeleteMemberByWorkspaceId(c)
}

func UserWorkspaceControllerInit(uws service.UserWorkspaceService) *UserWorkspaceControllerImpl {
	return &UserWorkspaceControllerImpl{
		Uws: uws,
	}
}
