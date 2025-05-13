package service

import (
	"net/http"
	"pm_go_version/app/constant"
	"pm_go_version/app/domain/dto"
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
	DeleteTaskById(c *gin.Context)
	UpdateTaskById(c *gin.Context)
	GetTaskById(c *gin.Context)
	BatchUpdateTask(c *gin.Context)
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
	if err := c.ShouldBind(&task); err != nil {
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
	workspaceId, err := strconv.Atoi(c.Param("workspaceId"))
	if err != nil {
		log.Error("Invalid workspace ID parameter value: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}
	var taskQuery entity.Task
	if err := c.ShouldBindQuery(&taskQuery); err != nil {
		log.Error("Failed to bind task: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}

	userIdValue, _ := c.Get("parse_id")
	userId := int(userIdValue.(uint))
	log.Debug("Task Query: ", taskQuery)
	result, err := ts.Tr.GetListofTasks(&taskQuery, uint(workspaceId), userId)
	if err != nil {
		log.Error("Failed to get list of tasks: ", err)
		pkg.PanicException(constant.UnknownError)
	}

	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, result))
}

func (ts *TaskServiceImpl) DeleteTaskById(c *gin.Context) {
	defer pkg.PanicHandler(c)
	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		log.Error("Invalid task ID format: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	userIdValue, _ := c.Get("parse_id")
	userId := int(userIdValue.(uint))

	result, err := ts.Tr.DeleteTaskById(taskId, userId)
	if !result {
		log.Error("Failed to delete task: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, taskId))
}

func (ts *TaskServiceImpl) UpdateTaskById(c *gin.Context) {
	defer pkg.PanicHandler(c)
	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		log.Error("Invalid task ID format: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	userIdValue, _ := c.Get("parse_id")
	userId := int(userIdValue.(uint))

	var updateTask entity.UpdateTask
	if err := c.ShouldBind(&updateTask); err != nil {
		log.Error("Failed to bind task: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}
	// updateTask := map[string]any{}
	// if task.Name != nil {
	// 	updateTask["name"] = *task.Name
	// }
	// if task.ProjectId != nil {
	// 	updateTask["project_id"] = uint(*task.ProjectId)
	// }
	// if task.AssigneeId != nil {
	// 	updateTask["assignee_id"] = uint(*task.AssigneeId)
	// }
	// if task.Description != nil {
	// 	updateTask["description"] = *task.Description
	// }
	// if task.DueDate != nil {
	// 	updateTask["due_date"] = *task.DueDate
	// }
	// if task.Status != nil {
	// 	updateTask["status"] = *task.Status
	// }
	// if task.Position != nil {
	// 	updateTask["position"] = *task.Position
	// }

	result, err := ts.Tr.UpdateTaskById(taskId, updateTask, userId)
	if !result {
		log.Error("Failed to update task: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, taskId))
}

func (ts *TaskServiceImpl) GetTaskById(c *gin.Context) {
	defer pkg.PanicHandler(c)
	taskId, err := strconv.Atoi(c.Param("taskId"))
	if err != nil {
		log.Error("Invalid task ID format: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	userIdValue, _ := c.Get("parse_id")
	userId := int(userIdValue.(uint))

	result, err := ts.Tr.GetTaskById(taskId, userId)
	if err != nil {
		log.Error("Failed to get task information: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, *result))
}

func (ts *TaskServiceImpl) BatchUpdateTask(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var batchTasks dto.BatchUpdateTaskDTO
	if err := c.ShouldBindJSON(&batchTasks); err != nil {
		log.Error("Failed to bind task: ", err)
		pkg.PanicException(constant.InvalidRequest)
	}
	result, err := ts.Tr.BatchUpdateTask(batchTasks)
	if !result {
		log.Error("Failed to batch update task: ", err)
		pkg.PanicException(constant.UnknownError)
	}
	c.JSON(http.StatusOK, pkg.BuildResponse(constant.Success, batchTasks))
}

func TaskServiceInit(taskRepository repository.TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{
		Tr: taskRepository,
	}
}
