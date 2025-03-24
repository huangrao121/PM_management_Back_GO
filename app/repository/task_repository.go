package repository

import (
	"pm_go_version/app/domain/entity"

	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(user_id uint, task *entity.Task) (bool, error)
	CheckMember(user_id uint, workspace_id uint) (bool, error)
	GetListofTasks(taskQuery entity.Task, workspace_id uint) ([]entity.Task, error)
}

type TaskRepositoryImpl struct {
	db *gorm.DB
}

func (tr *TaskRepositoryImpl) CreateTask(user_id uint, task *entity.Task) (bool, error) {
	r1 := tr.db.First(&entity.UserWorkspace{}, "user_id = ? AND workspace_id = ?", user_id, uint(task.WorkspaceId))
	if r1.Error != nil {
		return false, r1.Error
	}
	r2 := tr.db.Create(task)
	if r2.Error != nil {
		return false, r2.Error
	}
	return true, nil
}

func (tr *TaskRepositoryImpl) CheckMember(user_id uint, workspace_id uint) (bool, error) {
	r1 := tr.db.First(&entity.UserWorkspace{}, "user_id = ? AND workspace_id = ?", user_id, workspace_id)
	if r1.Error != nil {
		return false, r1.Error
	}
	return true, nil
}

func (tr *TaskRepositoryImpl) GetListofTasks(taskQuery entity.Task, workspace_id uint) ([]entity.Task, error) {
	var tasks []entity.Task
	r1 := tr.db.Model(&entity.Task{}).Where("workspace_id = ?", workspace_id)
	if taskQuery.Status != "" {
		r1 = r1.Where("status = ?", taskQuery.Status)
	}
	if taskQuery.AssigneeId != 0 {
		r1 = r1.Where("assignee_id = ?", taskQuery.AssigneeId)
	}
	if taskQuery.ProjectId != 0 {
		r1 = r1.Where("project_id = ?", taskQuery.ProjectId)
	}
	if !taskQuery.DueDate.IsZero() {
		r1 = r1.Where("due_date = ?", taskQuery.DueDate)
	}
	if taskQuery.Name != "" {
		r1 = r1.Where("name LIKE ?", "%"+taskQuery.Name+"%")
	}
	r1.Order("position ASC").Find(&tasks)
	if r1.Error != nil {
		return nil, r1.Error
	}
	return tasks, nil
}
func TaskRepositoryInit(db *gorm.DB) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{db: db}
}
