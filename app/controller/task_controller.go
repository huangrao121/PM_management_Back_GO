package controller

import (
	"pm_go_version/app/service"

	"github.com/gin-gonic/gin"
)

type TaskController interface {
	CreateTask(c *gin.Context)
	GetListofTasks(c *gin.Context)
	DeleteTaskById(c *gin.Context)
	UpdateTaskById(c *gin.Context)
	GetTaskById(c *gin.Context)
	BatchUpdateTask(c *gin.Context)
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

func (tc *TaskControllerImpl) DeleteTaskById(c *gin.Context) {
	tc.Ts.DeleteTaskById(c)
}

func (tc *TaskControllerImpl) UpdateTaskById(c *gin.Context) {
	tc.Ts.UpdateTaskById(c)
}

func (tc *TaskControllerImpl) GetTaskById(c *gin.Context) {
	tc.Ts.GetTaskById(c)
}

func (tc *TaskControllerImpl) BatchUpdateTask(c *gin.Context) {
	tc.Ts.BatchUpdateTask(c)
}

func TaskControllerInit(taskService service.TaskService) *TaskControllerImpl {
	return &TaskControllerImpl{
		Ts: taskService,
	}
}
