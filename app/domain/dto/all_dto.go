package dto

type UserDTO struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type WorkspaceDTO struct {
	ID          uint   `json:"id" form:"id"`
	Name        string `json:"name" form:"name"`
	CreaterID   string `json:"creater_id" form:"creater_id"`
	CreaterName string `json:"creater_user_name" form:"creater_user_name"`
	ImageUrl    string `json:"url"`
	InviteCode  string `json:"invite_code"`
}

type MembersDTO struct {
	UserId     uint   `json:"user_id"`
	UserMember string `json:"user_member"`
	UserName   string `json:"username" gorm:"column:username"`
	Email      string `json:"email"`
}
