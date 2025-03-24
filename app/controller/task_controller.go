package controller

import (
	"pm_go_version/app/service"

	"github.com/gin-gonic/gin"
)

type TaskController interface {
	CreateTask(c *gin.Context)
	GetListofTasks(c *gin.Context)
}

type TaskControllerImpl struct {
	Ts service.TaskService
}

func (tc *TaskControllerImpl) CreateTask(c *gin.Context) {
	tc.Ts.CreateTask(c)
}

func (tc *TaskControllerImpl) GetListofTasks(c *gin.Context) {
	tc.Ts.GetListofTasks(c)
}

func TaskControllerInit(taskService service.TaskService) *TaskControllerImpl {
	return &TaskControllerImpl{
		Ts: taskService,
	}
}
