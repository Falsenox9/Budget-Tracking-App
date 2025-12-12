package repository

import (
	"budget-tracking-app/entity"
	"budget-tracking-app/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(tx *gorm.DB, user *entity.User) error
	GetUser(param model.UserParam) (*entity.User, error)
	UpdateUser(tx *gorm.DB, user *entity.User) error
	DeleteUser(tx *gorm.DB, id string) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(tx *gorm.DB, user *entity.User) error {
	db := tx
	if db == nil {
		db = r.db
	}

	err := db.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUser(tx *gorm.DB, user *entity.User) error {
	db := tx
	if db == nil {
		db = r.db
	}

	err := db.Updates(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(param model.UserParam) (*entity.User, error) {
	var user entity.User

	query := r.db
	if param.UserID != uuid.Nil {
		query = query.Where("user_id = ?", param.UserID)
	}
	if param.Email != "" {
		query = query.Where("email = ?", param.Email)
	}

	err := query.First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) DeleteUser(tx *gorm.DB, id string) error {
	err := tx.Delete(&entity.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
