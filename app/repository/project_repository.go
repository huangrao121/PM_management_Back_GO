package repository

import (
	"pm_go_version/app/domain/entity"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	GetListofProjects(user_id uint, workspace_id uint) ([]entity.Project, error)
	CreateProject(user_id uint, project *entity.Project) (bool, error)
	GetProjectById(user_id uint, project_id uint) (entity.Project, error)
	UpdateProjectById(user_id uint, workspace_id uint, project_id uint, project *entity.Project) (bool, error)
}

type ProjectRepositoryImpl struct {
	db *gorm.DB
}

func (pr *ProjectRepositoryImpl) GetListofProjects(user_id uint, workspace_id uint) ([]entity.Project, error) {
	// isMember, err := CheckMember(pr.db, user_id, workspace_id)
	// if !isMember {
	// 	return []entity.Project{}, err
	// }
	// var projects []entity.Project
	// r := pr.db.Where("workspace_id=?", workspace_id).Find(&projects)
	// if r.Error != nil {
	// 	return []entity.Project{}, r.Error
	// }
	// return projects, nil
	var projects []entity.Project

	// 使用JOIN来一次性查询
	result := pr.db.Joins("JOIN user_workspaces ON user_workspaces.workspace_id = projects.workspace_id").
		Where("user_workspaces.user_id = ? AND projects.workspace_id = ?", user_id, workspace_id).
		Find(&projects)

	if result.Error != nil {
		return nil, result.Error
	}

	return projects, nil
}

func (pr *ProjectRepositoryImpl) CreateProject(user_id uint, project *entity.Project) (bool, error) {
	isMember, err := CheckMember(pr.db, user_id, uint(project.WorkspaceId))
	if !isMember {
		return false, err
	}
	r := pr.db.Create(project)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (pr *ProjectRepositoryImpl) GetProjectById(user_id uint, project_id uint) (entity.Project, error) {
	var project entity.Project
	r := pr.db.Joins("JOIN user_workspaces ON user_workspaces.workspace_id = projects.workspace_id").
		Where("user_workspaces.user_id = ? AND projects.id = ?", user_id, project_id).
		First(&project)
	if r.Error != nil {
		return entity.Project{}, r.Error
	}
	return project, nil
}

func (pr *ProjectRepositoryImpl) UpdateProjectById(user_id uint, workspace_id uint, project_id uint, project *entity.Project) (bool, error) {
	r1 := pr.db.Joins("JOIN user_workspaces ON user_workspaces.workspace_id = projects.workspace_id").
		Where("user_workspaces.user_id = ? AND projects.workspace_id = ? AND projects.id = ?", user_id, workspace_id, project_id).
		First(&project)
	if r1.Error != nil {
		return false, r1.Error
	}
	r2 := pr.db.Model(&entity.Project{}).Where("id = ?", project_id).Updates(project)
	if r2.Error != nil {
		return false, r2.Error
	}
	return true, nil
}

func ProjectRepositoryInit(db *gorm.DB) *ProjectRepositoryImpl {
	return &ProjectRepositoryImpl{db: db}
}
