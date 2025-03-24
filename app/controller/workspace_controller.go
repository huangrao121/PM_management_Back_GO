package controller

import (
	"pm_go_version/app/service"

	"github.com/gin-gonic/gin"
)

type WorkspaceController interface {
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

type WorkspaceControllerImpl struct {
	Ws service.WorkspaceService
}

func (wc *WorkspaceControllerImpl) GetListofWorkspaces(c *gin.Context) {
	wc.Ws.GetListofWorkspaces(c)
}

func (wc *WorkspaceControllerImpl) CreateWorkspace(c *gin.Context) {
	wc.Ws.CreateWorkspace(c)
}

func (wc *WorkspaceControllerImpl) GetWorkspacesById(c *gin.Context) {
	wc.Ws.GetWorkspacesById(c)
}

func (wc *WorkspaceControllerImpl) UpdateWorkspaceById(c *gin.Context) {
	wc.Ws.UpdateWorkspaceById(c)
}

func (wc *WorkspaceControllerImpl) DeleteWorkspaceById(c *gin.Context) {
	wc.Ws.DeleteWorkspaceById(c)
}

func (wc *WorkspaceControllerImpl) GetSingleWorkspaceById(c *gin.Context) {
	wc.Ws.GetSingleWorkspaceById(c)
}

func (wc *WorkspaceControllerImpl) ResetInvite(c *gin.Context) {
	wc.Ws.ResetInvite(c)
}

func (wc *WorkspaceControllerImpl) JoinWorkspace(c *gin.Context) {
	wc.Ws.JoinWorkspace(c)
}

func (wc *WorkspaceControllerImpl) GetWorkspaceInfo(c *gin.Context) {
	wc.Ws.GetWorkspaceInfo(c)
}

func WorkspaceControllerInit(ws service.WorkspaceService) *WorkspaceControllerImpl {
	return &WorkspaceControllerImpl{Ws: ws}
}
