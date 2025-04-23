package repository

import (
	"github.com/AronditFire/todo-app/entity"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (r *AuthRepo) CreateUser(user entity.UserRegisterRequest) error {
	var newUser entity.User
	newUser.Username = user.Username
	newUser.Password = user.Password
	newUser.IsAdmin = false

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&newUser).Error; err != nil { // error maybe
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *AuthRepo) GetUser(username string) (entity.User, error) {
	var user entity.User

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entity.User{}, err
	}

	if err := tx.First(&user, "username = ?", username).Error; err != nil {
		tx.Rollback()
		return entity.User{}, err
	}

	return user, tx.Commit().Error
}

func (r *AuthRepo) LoginUser(userLogin entity.UserAuthRequest) (entity.User, error) {
	var user entity.User

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return entity.User{}, err
	}

	if err := tx.First(&user, "username = ?", userLogin.Username).Error; err != nil {
		tx.Rollback()
		return entity.User{}, err
	}

	return user, tx.Commit().Error
}
