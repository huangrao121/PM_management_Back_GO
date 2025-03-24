package controller

import (
	"pm_go_version/app/service"

	"github.com/gin-gonic/gin"
)

type ProjectController interface {
	GetListofProjects(c *gin.Context)
	CreateProject(c *gin.Context)
	GetProjectById(c *gin.Context)
	UpdateProjectById(c *gin.Context)
}

type ProjectControllerImpl struct {
	Ps service.ProjectService
}

func (pc *ProjectControllerImpl) GetListofProjects(c *gin.Context) {
	pc.Ps.GetListofProjects(c)
}

func (pc *ProjectControllerImpl) CreateProject(c *gin.Context) {
	pc.Ps.CreateProject(c)
}

func (pc *ProjectControllerImpl) GetProjectById(c *gin.Context) {
	pc.Ps.GetProjectById(c)
}

func (pc *ProjectControllerImpl) UpdateProjectById(c *gin.Context) {
	pc.Ps.UpdateProjectById(c)
}

func ProjectControllerInit(ps service.ProjectService) *ProjectControllerImpl {
	return &ProjectControllerImpl{Ps: ps}
}
