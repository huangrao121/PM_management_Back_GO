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
	GetWorkspacesById(userId uint) ([]dto.WorkspaceDTO, error)
	DeleteWorkspaceById(userId uint, workspaceId uint) (bool, error)
	GetImageName(userId uint, workspaceId uint) (string, error)
	UpdateWorkspaceById(userId uint, workspaceId int, request *entity.UpdateWorkspace) (bool, error)
	GetSingleWorkspaceById(userId uint, workspaceId uint) (*dto.WorkspaceDTO, error)
	ResetInvite(userId uint, workspaceId uint, inviteCode string) (bool, error)
	JoinWorkspace(userId uint, workspaceId uint, code string) (*entity.UserWorkspace, error)
	GetWorkspaceInfo(workspaceId uint) (*entity.Workspace, error)
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

func (wr *WorkspaceRepositoryImpl) GetWorkspacesById(userId uint) ([]dto.WorkspaceDTO, error) {
	var result []dto.WorkspaceDTO
	r := wr.db.
		Debug().
		Table("user_workspaces uw").
		Select("w.id, w.name, w.creater_id, w.creater_user_name, w.image_url, w.invite_code").
		Joins("join workspaces w on uw.workspace_id = w.id").
		Where(
			"uw.user_id = ?", userId,
		).Scan(&result)
	if r.Error != nil {
		log.Error("Got and error when get list of workspaces by id. Error: ", r.Error)
		return nil, r.Error
	}
	log.Debug("workspace result: ", result)
	return result, nil
}

func (wr *WorkspaceRepositoryImpl) GetImageName(userId uint, workspaceId uint) (string, error) {
	var uw entity.UserWorkspace
	var workspace entity.Workspace
	err := wr.db.Transaction(func(db *gorm.DB) error {
		r1 := db.Where("user_id=? AND workspace_id=?", userId, workspaceId).First(&uw)
		if r1.Error != nil {
			log.Error("Got and error when check if user have workspaces. Error: ", r1.Error)
			return r1.Error
		}

		if uw.UserMember != "Owner" {
			log.Error("User is not owner to delete workspace. Error: ")
			return errors.New("user is not owner")
		}
		r2 := db.Where("id=?", workspaceId).First(&workspace)
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

func (wr *WorkspaceRepositoryImpl) DeleteWorkspaceById(userId uint, workspaceId uint) (bool, error) {
	var uw entity.UserWorkspace
	var workspace entity.Workspace
	err := wr.db.Transaction(func(db *gorm.DB) error {

		if isOwner, _ := wr.checkOwner(uw, userId, workspaceId); !isOwner {
			log.Error("User is not owner to delete workspace. Error: ")
			return errors.New("user is not owner")
		} else {
			r2 := db.Where("id=?", workspaceId).Delete(&workspace)
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

func (wr *WorkspaceRepositoryImpl) UpdateWorkspaceById(userId uint, workspaceId int, request *entity.UpdateWorkspace) (bool, error) {

	r1 := wr.db.Debug().Model(&entity.Workspace{}).
		Where("id=? and creater_id=?", workspaceId, userId).
		Updates(*request)
	if r1.Error != nil {
		log.Error("Failed to update workspace: ", r1.Error)
		return false, r1.Error
	}
	return true, nil
}

func (wr *WorkspaceRepositoryImpl) GetSingleWorkspaceById(userId uint, workspaceId uint) (*dto.WorkspaceDTO, error) {
	// var uw entity.UserWorkspace
	var result dto.WorkspaceDTO
	r := wr.db.Model(&entity.Workspace{}).
		Select("id, name, creater_id, creater_user_name, image_url, invite_code").
		Joins("join user_workspaces uw on uw.workspace_id=workspaces.id").
		Where("workspaces.id = ? and uw.user_id=?", workspaceId, userId).
		Scan(&result)
	if r.Error != nil {
		log.Error("Got and error when user try to delete the workspace. Error: ", r.Error)
		return &dto.WorkspaceDTO{}, r.Error
	}
	return &result, nil
}

func (wr *WorkspaceRepositoryImpl) ResetInvite(userId uint, workspaceId uint, inviteCode string) (bool, error) {
	r := wr.db.Debug().Model(&entity.Workspace{}).Where("creater_id=? and id=?", userId, workspaceId).
		Update("invite_code", inviteCode)

	if r.Error != nil {
		log.Error("Failed to reset invite code: ", r.Error)
		return false, r.Error
	}
	return true, nil
}

func (wr *WorkspaceRepositoryImpl) JoinWorkspace(userId uint, workspaceId uint, code string) (*entity.UserWorkspace, error) {
	userWorkspace := entity.UserWorkspace{UserID: uint(userId), WorkspaceID: uint(workspaceId)}
	err := wr.db.Transaction(func(db *gorm.DB) error {
		//var uw entity.UserWorkspace
		var ws entity.Workspace
		r1 := db.First(&userWorkspace)
		if errors.Is(r1.Error, gorm.ErrRecordNotFound) {
			// log.Error("The member have already joined. Error: ", r1.Error)
			// return r1.Error

			r2 := db.Where("id=?", workspaceId).First(&ws)
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
		// r2 := db.Where("id=?", workspaceId).First(&ws)
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

func (wr *WorkspaceRepositoryImpl) GetWorkspaceInfo(workspaceId uint) (*entity.Workspace, error) {
	var ws entity.Workspace
	r1 := wr.db.Select("name").Where("id=?", workspaceId).First(&ws)
	if r1.Error != nil {
		log.Error("Got and error when check if user have workspaces. Error: ", r1.Error)
		return &entity.Workspace{}, r1.Error
	}
	return &ws, nil
}

func (wr *WorkspaceRepositoryImpl) checkOwner(uw entity.UserWorkspace, userId uint, workspaceId uint) (bool, error) {
	r1 := wr.db.Where("user_id=? AND workspace_id=?", userId, workspaceId).First(&uw)
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

func checkOwnerE(db *gorm.DB, userId uint, workspaceId uint) (bool, error) {
	var uw entity.UserWorkspace
	r1 := db.Where("user_id=? AND workspace_id=?", userId, workspaceId).First(&uw)
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
