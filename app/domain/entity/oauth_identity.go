package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OAuthIdentity struct {
	ID             uuid.UUID `gorm:"column:id;primaryKey;unique;not null;default:uuid_generate_v4()"`
	UserID         uint      `gorm:"column:user_id;not null"`
	Provider       string    `gorm:"column:provider;not null"`
	ProviderUserID string    `gorm:"column:provider_user_id;not null"`
	Email          string    `gorm:"column:email;not null"`
	User           User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	BaseModel
}

func (oi *OAuthIdentity) BeforeCreate(tx *gorm.DB) (err error) {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return
}
