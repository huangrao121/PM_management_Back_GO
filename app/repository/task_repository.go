package repository

import (
	"pm_go_version/app/domain/dto"
	"pm_go_version/app/domain/entity"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaskRepository interface {
	CreateTask(user_id uint, task *entity.Task) (bool, error)
	CheckMember(user_id uint, workspace_id uint) (bool, error)
	GetListofTasks(taskQuery *entity.Task, workspace_id uint, user_id int) ([]*entity.TaskInfo, error)
	DeleteTaskById(taskId int, user_id int) (bool, error)
	GetTaskById(taskId int, user_id int) (*entity.TaskInfo, error)
	UpdateTaskById(taskId int, task entity.UpdateTask, user_id int) (bool, error)
	BatchUpdateTask(batchTasks dto.BatchUpdateTaskDTO) (bool, error)
}

type TaskRepositoryImpl struct {
	db *gorm.DB
}

func (tr *TaskRepositoryImpl) CreateTask(user_id uint, task *entity.Task) (bool, error) {
	r := tr.db.Transaction(func(db *gorm.DB) error {
		r1 := db.First(&entity.UserWorkspace{}, "user_id = ? AND workspace_id = ?", user_id, uint(task.WorkspaceId))
		if r1.Error != nil {
			return r1.Error
		}
		r2 := db.Create(task)
		if r2.Error != nil {
			return r2.Error
		}
		return nil
	})
	if r != nil {
		return false, r
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

func (tr *TaskRepositoryImpl) GetListofTasks(taskQuery *entity.Task, workspace_id uint, user_id int) ([]*entity.TaskInfo, error) {
	//var results []interface{}
	var task []*entity.TaskInfo
	//r1 := tr.db.Model(&entity.Task{}).Where("workspace_id = ?", workspace_id)
	subQuery := tr.db.Model(&entity.Task{}).Where("workspace_id = ?", workspace_id)
	r1 := tr.db.Debug().Table("(?) u", subQuery).Select(`u.id,
		u.name,
		u.project_id,
		u.workspace_id,
		u.assignee_id,
		u.description,
		u.due_date,
		u.status,
		u.position,
		p.name project_name,
		p.image_url project_image,
		us.username assignee_name,
		us.email assignee_email`).
		Joins("JOIN projects p ON u.project_id = p.id").
		Joins("JOIN users us ON u.assignee_id = us.id").
		Joins("JOIN user_workspaces uw ON u.workspace_id = uw.workspace_id AND uw.user_id=?", user_id)
	if (*taskQuery).Status != "" {
		r1 = r1.Where("u.status = ?", (*taskQuery).Status)
	}
	if (*taskQuery).AssigneeId != 0 {
		r1 = r1.Where("u.assignee_id = ?", (*taskQuery).AssigneeId)
	}
	if (*taskQuery).ProjectId != 0 {
		r1 = r1.Where("u.project_id = ?", (*taskQuery).ProjectId)
	}
	if !(*taskQuery).DueDate.IsZero() {
		r1 = r1.Where("u.due_date = ?", (*taskQuery).DueDate)
	}
	if (*taskQuery).Name != "" {
		r1 = r1.Where("u.name LIKE ?", "%"+(*taskQuery).Name+"%")
	}
	log.Debug("Task Query project_id: ", (*taskQuery).ProjectId)
	r1.Order("u.position ASC").Scan(&task)
	if r1.Error != nil {
		return nil, r1.Error
	}
	return task, nil
}

func (tr *TaskRepositoryImpl) DeleteTaskById(taskId int, user_id int) (bool, error) {
	r1 := tr.db.Debug().Joins("JOIN user_workspace uw on u.workspace_id = uw.workspace_id AND uw.user_id=?", user_id).Delete(&entity.Task{}, taskId)
	if r1.Error != nil {
		return false, r1.Error
	}
	return true, nil
}

func (tr *TaskRepositoryImpl) GetTaskById(taskId int, user_id int) (*entity.TaskInfo, error) {
	subQuery := tr.db.Model(&entity.Task{}).Where("id=?", taskId)
	r1 := tr.db.Debug().Table("(?) u", subQuery).Select(`u.id,
		u.name,
		u.project_id,
		u.workspace_id,
		u.assignee_id,
		u.description,
		u.due_date,
		u.status,
		u.position,
		p.name project_name,
		p.image_url project_image,
		us.username assignee_name,
		us.email assignee_email`).
		Joins("JOIN projects p ON u.project_id = p.id").
		Joins("JOIN users us ON u.assignee_id = us.id").
		Joins("JOIN user_workspaces uw ON u.workspace_id = uw.workspace_id AND uw.user_id=?", user_id)
	var singleTask entity.TaskInfo
	r1.Scan(&singleTask)
	if r1.Error != nil {
		return nil, r1.Error
	}
	return &singleTask, nil
}

func (tr *TaskRepositoryImpl) UpdateTaskById(taskId int, task entity.UpdateTask, user_id int) (bool, error) {
	// 获取要更新的字段
	filteredTask := tr.db.Model(&entity.Task{}).Where("id = ?", taskId).
		Joins("JOIN user_workspaces uw ON u.workspace_id = uw.workspace_id AND uw.user_id=?", user_id)
	r1 := filteredTask.Updates(task)

	if r1.Error != nil {
		log.Error("Failed to update task: ", r1.Error)
		return false, r1.Error
	}
	return true, nil
}

func (tr *TaskRepositoryImpl) BatchUpdateTask(batchTasks dto.BatchUpdateTaskDTO) (bool, error) {
	// 准备批量数据
	var tasks []entity.Task
	for _, task := range batchTasks.Tasks {
		tasks = append(tasks, entity.Task{
			ID:       *task.Id,
			Status:   *task.Status,
			Position: *task.Position,
		})
	}

	err := tr.db.Transaction(func(tx *gorm.DB) error {
		// 使用 CreateInBatches 批量处理
		r1 := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"status", "position"}),
		}).CreateInBatches(tasks, 100) // 每批100条记录

		return r1.Error
	})

	if err != nil {
		log.Error("Failed to batch update tasks: ", err)
		return false, err
	}
	return true, nil
}

func TaskRepositoryInit(db *gorm.DB) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{db: db}
}
