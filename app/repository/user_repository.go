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
	GetUserByEmail(string) (*entity.User, error)
	GetOAuthIdentity(string, string, string) (*entity.OAuthIdentity, error)
	CreateUserOAuthIdentity(email, name, providerID, provider string) (*entity.User, error)
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

func (ur *UserRepositoryImpl) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	r := ur.db.Where("email=?", email).First(&user)
	if r.Error != nil {
		return &entity.User{}, r.Error
	}
	return &user, nil
}

func (ur *UserRepositoryImpl) GetOAuthIdentity(email, sub, provider string) (*entity.OAuthIdentity, error) {
	var identity entity.OAuthIdentity
	r := ur.db.Where("email=? AND provider_user_id=? AND provider=?", email, sub, provider).First(&identity)
	if r.Error != nil {
		return &entity.OAuthIdentity{}, r.Error
	}
	return &identity, nil
}

func (ur *UserRepositoryImpl) CreateUserOAuthIdentity(email, name, providerID, provider string) (*entity.User, error) {
	var user entity.User
	err := ur.db.Transaction(func(db *gorm.DB) error {
		user.Email = email
		user.UserName = name
		if err := db.Create(&user).Error; err != nil {
			return err
		}
		var identity entity.OAuthIdentity
		identity.Email = email
		identity.Provider = provider
		identity.ProviderUserID = providerID
		identity.UserID = user.ID
		if err := db.Create(&identity).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Error("Failed to create user and oauth identity in db transaction", err)
		return nil, err
	}
	return &user, nil
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
