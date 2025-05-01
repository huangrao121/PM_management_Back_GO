package repository

import (
	"pm_go_version/app/domain/dto"
	"pm_go_version/app/domain/entity"

	"errors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WorkspaceRepository interface {
	GetWorkspaces() ([]entity.Workspace, error)
	CreateWorkspace(request *entity.Workspace) (uint, error)
	GetWorkspacesById(id uint) ([]dto.WorkspaceDTO, error)
	DeleteWorkspaceById(user_id uint, workspace_id uint) (bool, error)
	GetImageName(user_id uint, workspace_id uint) (string, error)
	UpdateWorkspaceById(user_id uint, request *dto.WorkspaceDTO) (bool, error)
	GetSingleWorkspaceById(user_id uint, workspace_id uint) (entity.Workspace, error)
	ResetInvite(user_id uint, workspace_id uint, invite_code string) (bool, error)
	JoinWorkspace(user_id uint, workspace_id uint, code string) (*entity.UserWorkspace, error)
	GetWorkspaceInfo(workspace_id uint) (*entity.Workspace, error)
}

type WorkspaceRepositoryImpl struct {
	db *gorm.DB
}

func (wr *WorkspaceRepositoryImpl) CreateWorkspace(request *entity.Workspace) (uint, error) {
	//var workspace entity.Workspace
	r := wr.db.Create(request)

	if r.Error != nil {
		log.Error("Got and error when create a new workspace. Error: ", r.Error)
		return 0, r.Error
	}
	return request.ID, nil
}

func (wr *WorkspaceRepositoryImpl) GetWorkspaces() ([]entity.Workspace, error) {
	var workspaces []entity.Workspace
	r := wr.db.Find(&workspaces)
	if r.Error != nil {
		log.Error("Got and error when get list of all workspaces. Error: ", r.Error)
		return nil, r.Error
	}
	return workspaces, nil
}

func (wr *WorkspaceRepositoryImpl) GetWorkspacesById(id uint) ([]dto.WorkspaceDTO, error) {
	var result []dto.WorkspaceDTO
	r := wr.db.
		Debug().
		Table("user_workspaces").
		Select("id, name, creater_id, creater_user_name, image_url").
		Joins("left join workspaces on user_workspaces.workspace_id = workspaces.id").
		Where(
			"user_workspaces.user_id = ?", id,
		).Scan(&result)
	if r.Error != nil {
		log.Error("Got and error when get list of workspaces by id. Error: ", r.Error)
		return nil, r.Error
	}
	return result, nil
}

func (wr *WorkspaceRepositoryImpl) GetImageName(user_id uint, workspace_id uint) (string, error) {
	var uw entity.UserWorkspace
	var workspace entity.Workspace
	err := wr.db.Transaction(func(db *gorm.DB) error {
		r1 := db.Where("user_id=? AND workspace_id=?", user_id, workspace_id).First(&uw)
		if r1.Error != nil {
			log.Error("Got and error when check if user have workspaces. Error: ", r1.Error)
			return r1.Error
		}

		if uw.UserMember != "Owner" {
			log.Error("User is not owner to delete workspace. Error: ")
			return errors.New("user is not owner")
		}
		r2 := db.Where("id=?", workspace_id).First(&workspace)
		if r2.Error != nil {
			log.Error("No such workspace in database. Error: ", r1.Error)
			return r2.Error
		}
		return nil
	})

	if err != nil {
		return "", errors.New("user is not owner")
	}
	return workspace.ImageUrl, nil
}

func (wr *WorkspaceRepositoryImpl) DeleteWorkspaceById(user_id uint, workspace_id uint) (bool, error) {
	var uw entity.UserWorkspace
	var workspace entity.Workspace
	err := wr.db.Transaction(func(db *gorm.DB) error {

		if isOwner, _ := wr.checkOwner(uw, user_id, workspace_id); !isOwner {
			log.Error("User is not owner to delete workspace. Error: ")
			return errors.New("user is not owner")
		} else {
			r2 := db.Where("id=?", workspace_id).Delete(&workspace)
			if r2.Error != nil {
				log.Error("Got and error when user try to delete the workspace. Error: ", r2.Error)
				return r2.Error
			}
			return nil
		}
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (wr *WorkspaceRepositoryImpl) UpdateWorkspaceById(user_id uint, request *dto.WorkspaceDTO) (bool, error) {
	var uw entity.UserWorkspace
	err := wr.db.Transaction(func(db *gorm.DB) error {

		// r1 := db.Where("user_id=? AND workspace_id=?", user_id, request.ID).First(&uw)
		// if r1.Error != nil {
		// 	log.Error("Got and error when check if user have workspaces. Error: ", r1.Error)
		// 	return r1.Error
		// }
		if isOwner, _ := wr.checkOwner(uw, user_id, request.ID); !isOwner {
			log.Error("User is not owner to delete workspace. Error: ")
			return errors.New("user is not owner")
		} else {
			r2 := db.Model(&entity.Workspace{}).Where("id=?", request.ID).Select("name", "image_url").Updates(map[string]interface{}{"name": request.Name, "image_url": request.ImageUrl})
			if r2.Error != nil {
				log.Error("Got and error when user try to delete the workspace. Error: ", r2.Error)
				return r2.Error
			}
			return nil
		}
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (wr *WorkspaceRepositoryImpl) GetSingleWorkspaceById(user_id uint, workspace_id uint) (entity.Workspace, error) {
	// var uw entity.UserWorkspace
	var workspace entity.Workspace
	err := wr.db.Transaction(func(db *gorm.DB) error {
		// if isOwner, _ := wr.checkOwner(uw, user_id, workspace_id); !isOwner {
		// 	log.Error("User is not owner to delete workspace. Error: ")
		// 	return errors.New("user is not owner")
		// } else {
			r2 := db.Table("workspaces").
				Where("id = ?", workspace_id).
				Scan(&workspace)
			if r2.Error != nil {
				log.Error("Got and error when user try to delete the workspace. Error: ", r2.Error)
				return r2.Error
			}
			return nil
		// }
	})
	if err != nil {
		return entity.Workspace{}, err
	}
	return workspace, nil
}

func (wr *WorkspaceRepositoryImpl) ResetInvite(user_id uint, workspace_id uint, invite_code string) (bool, error) {
	var uw entity.UserWorkspace
	//var workspace dto.WorkspaceDTO
	err := wr.db.Transaction(func(db *gorm.DB) error {
		if isOwner, _ := wr.checkOwner(uw, user_id, workspace_id); !isOwner {
			log.Error("User is not owner to delete workspace. Error: ")
			return errors.New("user is not owner")
		} else {
			r2 := db.Model(&entity.Workspace{}).Where("id=?", workspace_id).Update("invite_code", invite_code)
			if r2.Error != nil {
				log.Error("Got and error when user try to delete the workspace. Error: ", r2.Error)
				return r2.Error
			}
			return nil
		}
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (wr *WorkspaceRepositoryImpl) JoinWorkspace(user_id uint, workspace_id uint, code string) (*entity.UserWorkspace, error) {
	userWorkspace := entity.UserWorkspace{UserID: uint(user_id), WorkspaceID: uint(workspace_id)}
	err := wr.db.Transaction(func(db *gorm.DB) error {
		//var uw entity.UserWorkspace
		var ws entity.Workspace
		r1 := db.First(&userWorkspace)
		if errors.Is(r1.Error, gorm.ErrRecordNotFound) {
			// log.Error("The member have already joined. Error: ", r1.Error)
			// return r1.Error

			r2 := db.Where("id=?", workspace_id).First(&ws)
			if r2.Error != nil {
				log.Error("Failed to get the related workspace by this id. Error: ", r2.Error)
				return r2.Error
			}
			log.Info("retrieve from the database:: ", ws.InviteCode)
			log.Info("get from the front page: ", code)
			if ws.InviteCode != code {
				return errors.New("invite code is not correct")
			}
			userWorkspace.UserMember = "member"
			r3 := db.Create(&userWorkspace)
			if r3.Error != nil {
				log.Error("Failed to create entry in user_workspace database. Error: ", r3.Error)
				return r3.Error
			}
			return nil
		} else {
			log.Error("The member have already joined. Error: ", r1.Error)
			return r1.Error
		}
		// r2 := db.Where("id=?", workspace_id).First(&ws)
		// if r2.Error != nil {
		// 	log.Error("Failed to get the related workspace by this id. Error: ", r2.Error)
		// 	return r2.Error
		// }

		// if ws.InviteCode != code {
		// 	return errors.New("invite code is not correct")
		// }
		// userWorkspace.UserMember = "member"
		// r3 := db.Create(&userWorkspace)
		// if r3.Error != nil {
		// 	log.Error("Failed to create entry in user_workspace database. Error: ", r3.Error)
		// 	return r3.Error
		// }
		// return nil
	})
	if err != nil {
		return &entity.UserWorkspace{}, err
	}
	return &userWorkspace, nil
}

func (wr *WorkspaceRepositoryImpl) GetWorkspaceInfo(workspace_id uint) (*entity.Workspace, error) {
	var ws entity.Workspace
	r1 := wr.db.Select("name").Where("id=?", workspace_id).First(&ws)
	if r1.Error != nil {
		log.Error("Got and error when check if user have workspaces. Error: ", r1.Error)
		return &entity.Workspace{}, r1.Error
	}
	return &ws, nil
}

func (wr *WorkspaceRepositoryImpl) checkOwner(uw entity.UserWorkspace, user_id uint, workspace_id uint) (bool, error) {
	r1 := wr.db.Where("user_id=? AND workspace_id=?", user_id, workspace_id).First(&uw)
	if r1.Error != nil {
		log.Error("Got and error when check if user have workspaces. Error: ", r1.Error)
		return false, r1.Error
	}
	if uw.UserMember != "Owner" {
		log.Error("User is not owner to delete workspace. Error: ")
		return false, errors.New("user is not owner")
	} else {
		return true, nil
	}
}

func checkOwnerE(db *gorm.DB, user_id uint, workspace_id uint) (bool, error) {
	var uw entity.UserWorkspace
	r1 := db.Where("user_id=? AND workspace_id=?", user_id, workspace_id).First(&uw)
	if r1.Error != nil {
		log.Error("Got and error when check if user have workspaces. Error: ", r1.Error)
		return false, r1.Error
	}
	if uw.UserMember != "Owner" {
		log.Error("User is not owner to delete workspace. Error: ")
		return false, errors.New("user is not owner")
	} else {
		return true, nil
	}
}

func WorkspaceRepositoryInit(db *gorm.DB) *WorkspaceRepositoryImpl {
	return &WorkspaceRepositoryImpl{
		db: db,
	}
}
