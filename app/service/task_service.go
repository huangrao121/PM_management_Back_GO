package service

import (
	"net/http"
	"pm_go_version/app/constant"
	"pm_go_version/app/domain/entity"
	"pm_go_version/app/pkg"
	"pm_go_version/app/repository"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type TaskService interface {
	CreateTask(c *gin.Context)
	GetListofTasks(c *gin.Context)
}

type TaskServiceImpl struct {
	Tr repository.TaskRepository
}

func (ts *TaskServiceImpl) CreateTask(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userIDValue, exists := c.Get("parse_id")
	if !exists {
		log.Error("User ID not found in context")
		pkg.PanicException(constant.UnknownError)
	}

	userId, isInt := ConvertAnyToInt(userIDValue)
	if !isInt {
		log.Error("Failed to convert user ID: ", userId)
		pkg.PanicException(constant.UnknownError)
	}

	var task entity.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Error("Failed to bind task: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}
	log.Info("Task: ", task)
	isTask, err := ts.Tr.CreateTask(userId, &task)
	if !isTask {
		log.Error("Failed to create task: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, task))
}

func (ts *TaskServiceImpl) GetListofTasks(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var taskQuery entity.Task
	if err := c.ShouldBindQuery(&taskQuery); err != nil {
		log.Error("Failed to bind task: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	userIDValue, exists := c.Get("parse_id")
	if !exists {
		log.Error("User ID not found in context")
		pkg.PanicException(constant.UnknownError)
	}

	userId, isInt := ConvertAnyToInt(userIDValue)
	if !isInt {
		log.Error("Failed to convert user ID: ", userId)
	}

	workspaceId, err := strconv.ParseUint(c.Param("workspaceId"), 10, 32)
	if err != nil {
		log.Error("Invalid workspace ID format: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	checkMember, err := ts.Tr.CheckMember(userId, uint(workspaceId))
	if !checkMember {
		log.Error("User is not a member of the workspace")
		pkg.PanicException(constant.UnknownError)
	}

	result, err := ts.Tr.GetListofTasks(taskQuery, uint(workspaceId))
	if err != nil {
		log.Error("Failed to get list of tasks: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, result))
}

func TaskServiceInit(taskRepository repository.TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{
		Tr: taskRepository,
	}
}
