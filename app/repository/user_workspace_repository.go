package repository

import (
	"pm_go_version/app/domain/dto"
	"pm_go_version/app/domain/entity"

	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

type UserWorkspaceRepository interface {
	GetListofMembersByWorkspaceId(user_id uint, workspace_id uint) ([]dto.MembersDTO, error)
	DeleteMemberByWorkspaceId(user_id uint, workspace_id uint, member_id uint) (bool, error)
}

type UserWorkspaceRepositoryImpl struct {
	db *gorm.DB
}

func (uwr *UserWorkspaceRepositoryImpl) GetListofMembersByWorkspaceId(user_id uint, workspace_id uint) ([]dto.MembersDTO, error) {
	var members []dto.MembersDTO
	//var uw entity.UserWorkspace
	// isOwner, _ := checkOwnerE(uwr.db, uw, user_id, workspace_id)
	// if !isOwner {
	// 	log.Error("User is not owner to delete workspace. Error: ")
	// 	return []dto.MembersDTO{}, errors.New("user is not owner")
	// }
	// r := uwr.db.
	// 	Debug().
	// 	Table("(?) as uw", uwr.db.Model(&entity.UserWorkspace{}).
	// 		Select("user_id", "user_member").
	// 		Where("workspace_id=?", workspace_id)).
	// 	Select("user_id, user_member, username, email").
	// 	Joins("left join users on uw.user_id=users.id").
	// 	Scan(&members)
	r := uwr.db.Model(&entity.UserWorkspace{}).
		Joins("join workspaces on user_workspaces.workspace_id = workspaces.id").
		Joins("join users on user_workspaces.user_id = users.id").
		Where("user_workspaces.user_id = ? AND user_workspaces.workspace_id = ?", user_id, workspace_id).
		Select("user_workspaces.user_id, user_workspaces.user_member, users.username, users.email").
		Scan(&members)
	if r.Error != nil {
		log.Error("Got and error when get list of workspaces by id. Error: ", r.Error)
		return nil, r.Error
	}
	return members, nil
}

func (uwr *UserWorkspaceRepositoryImpl) DeleteMemberByWorkspaceId(user_id uint, workspace_id uint, member_id uint) (bool, error) {
	uw := entity.UserWorkspace{}
	//uw2 := entity.UserWorkspace{}
	isOwner, err := checkOwnerE(uwr.db, user_id, workspace_id)
	if !isOwner {
		log.Error("User doesn't have authority to delete. Error: ", err)
		return false, err
	}
	result := uwr.db.Where("user_id=? and workspace_id=?", member_id, workspace_id).Delete(&uw)
	//log.Info("rows affected is: ", result.RowsAffected)
	if result.Error != nil || result.RowsAffected == 0 {
		log.Error("Failed to delete member from workspaces. Error: ", err)
		return false, err
	}
	return true, nil
}

func UserWorkspaceRepositoryInit(db *gorm.DB) *UserWorkspaceRepositoryImpl {
	return &UserWorkspaceRepositoryImpl{db: db}
}

func CheckMember(db *gorm.DB, user_id uint, workspace_id uint) (bool, error) {
	var uw entity.UserWorkspace
	r1 := db.Where("user_id=? AND workspace_id=?", user_id, workspace_id).First(&uw)
	if r1.Error != nil {
		log.Error("The user is not member of workspace. Error: ", r1.Error)
		return false, r1.Error
	}
	return true, nil
}
