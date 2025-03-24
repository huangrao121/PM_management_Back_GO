package repository

import (
	//"pm_go_version/app/domain/dto"
	"pm_go_version/app/domain/entity"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Login(*entity.User) (entity.User, error)
	Save(*entity.User) (entity.User, error)
	//GetUserInfo() (entity.User, error)
}

type UserRepositoryImpl struct{ db *gorm.DB }

func (u *UserRepositoryImpl) Login(user *entity.User) (entity.User, error) {
	var result entity.User

	r := u.db.Where("email=? AND password=?", user.Email, user.Password).First(&result)
	if r.Error != nil {
		log.Error("Got and error when log user by email and password. Error: ", r.Error)
		return entity.User{}, r.Error
	}

	return result, nil
}

func (u *UserRepositoryImpl) Save(user *entity.User) (entity.User, error) {
	r := u.db.Create(user)
	if r.Error != nil {
		log.Error("Got and error when create new account. Error: ", r.Error)
		return entity.User{}, r.Error
	}

	return *user, nil
}

// func (u *UserRepositoryImpl) GetListUser() (entity.User, error) {
// 	var user entity.User
// 	r := u.db.Find(&user)
// 	if r.Error != nil {
// 		log.Error("Got and error when get list of users. Error: ", r.Error)
// 		return entity.User{}, r.Error
// 	}
// 	return user, nil
// }

func UserRepositoryInit(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}
