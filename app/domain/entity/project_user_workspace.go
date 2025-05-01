package entity

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Task status constants
const (
	TaskStatusBacklog    = "BACKLOG"
	TaskStatusTodo       = "TODO"
	TaskStatusInProgress = "IN_PROGRESS"
	TaskStatusInReview   = "IN_REVIEW"
	TaskStatusDone       = "DONE"
	TaskStatusBlocked    = "BLOCKED"
)

// ValidTaskStatuses contains all valid task statuses
var ValidTaskStatuses = []string{
	TaskStatusBacklog,
	TaskStatusTodo,
	TaskStatusInProgress,
	TaskStatusInReview,
	TaskStatusDone,
	TaskStatusBlocked,
}

type BaseModel struct {
	CreatedAt time.Time `gorm:"column:create_at;autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"column:update_at;autoUpdateTime" json:"-"`
}

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;unique;column:id" json:"-"`
	UserName string `gorm:"column:username; unique; not null" json:"username"`
	Email    string `gorm:"column:email;unique;not null" json:"email"`
	Password string `gorm:"column:password;->:false;<-:create" json:"password"`
	//Workspaces []UserWorkspace `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Workspaces []Workspace `gorm:"many2many:user_workspaces;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	BaseModel
}

type Workspace struct {
	ID   uint   `gorm:"column:id;primaryKey;autoIncrement;unique" json:"-"`
	Name string `gorm:"column:name; not null" json:"name" form:"name" binding:"required"`
	//CreaterID   int    `gorm:"column:creater_id" json:"creater_id" form:"creater_id" binding:"required"`
	CreaterID uint `gorm:"column:creater_id" json:"creater_id" form:"creater_id"`
	//CreaterName string `gorm:"column:creater_user_name" json:"creater_user_name" form:"creater_user_name" binding:"required"`
	CreaterName string `gorm:"column:creater_user_name" json:"creater_user_name" form:"creater_user_name"`
	ImageUrl    string `gorm:"column:image_url; unique" json:"image_url"`
	InviteCode  string `gorm:"column:invite_code;unique;not null" json:"invite_code" form:"invite_code"`
	//Users       []UserWorkspace `gorm:"foreignKey:WorkspaceID;references:ID"`
	Projects []Project
	BaseModel
}

type UserWorkspace struct {
	UserID      uint   `gorm:"column:user_id;primaryKey"`
	WorkspaceID uint   `gorm:"column:workspace_id;primaryKey"`
	UserMember  string `gorm:"column:user_member"`
}

type Project struct {
	ID          uint   `gorm:"column:id;primaryKey;autoIncrement;unique" json:"id"`
	Name        string `gorm:"column:name;not null" json:"name" form:"name"`
	ImageUrl    string `gorm:"column:image_url; unique" json:"image_url" form:"image_url"`
	WorkspaceId uint   `gorm:"column:workspace_id" json:"workspace_id" form:"workspace_id"`
	Tasks       []Task `gorm:"foreignKey:ProjectId;references:ID"`
	BaseModel
}

type Task struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement;unique" json:"id"`
	Name        string    `gorm:"column:name;not null" json:"name" form:"name"`
	ProjectId   uint      `gorm:"column:project_id;not null" json:"project_id" form:"project_id"`
	WorkspaceId uint      `gorm:"column:workspace_id;not null" json:"workspace_id" form:"workspace_id"`
	AssigneeId  uint      `gorm:"column:assignee_id;not null" json:"assignee_id" form:"assignee_id"`
	Description string    `gorm:"column:description" json:"description" form:"description"`
	DueDate     time.Time `gorm:"column:due_date;not null" json:"due_date" form:"due_date"`
	Status      string    `gorm:"column:status;type:varchar(20);check:status IN ('BACKLOG','TODO','IN_PROGRESS','IN_REVIEW','DONE','BLOCKED')" json:"status" form:"status"`
	Position    int       `gorm:"column:position;not null;default:0" json:"position" form:"position"`
	BaseModel
}

type UpdateTask struct {
	Name        *string    `json:"name"`
	ProjectId   *uint      `json:"project_id"`
	WorkspaceId *uint      `json:"workspace_id"`
	AssigneeId  *uint      `json:"assignee_id"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Status      *string    `json:"status"`
}

type TaskInfo struct {
	Task
	ProjectName   string `json:"project_name"`
	AssigneeName  string `json:"assignee_name"`
	AssigneeEmail string `json:"assignee_email"`
	ProjectImage  string `json:"project_image"`
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	user.UpdatedAt = time.Now()
	return
}

func (uw *UserWorkspace) BeforeDelete(db *gorm.DB) (err error) {
	if uw.UserMember == "member" {
		return errors.New("unauthorized delete")
	}
	return
}

// ValidateStatus checks if the given status is valid
func (t *Task) ValidateStatus() error {
	for _, validStatus := range ValidTaskStatuses {
		if t.Status == validStatus {
			return nil
		}
	}
	return fmt.Errorf("invalid task status: %s", t.Status)
}

// BeforeSave validates the task status before saving to database
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	return t.ValidateStatus()
}

func (ut *UpdateTask) BeforeUpdate(tx *gorm.DB) error {
	for _, validStatus := range ValidTaskStatuses {
		if ut.Status != nil && *ut.Status == validStatus {
			return nil
		}
	}
	return fmt.Errorf("invalid task status")
}
