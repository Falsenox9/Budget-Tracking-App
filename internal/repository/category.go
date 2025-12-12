package repository

import (
	"budget-tracking-app/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICategoryRepository interface {
	Create(tx *gorm.DB, category *entity.Category) error
	FindByID(tx *gorm.DB, categoryID uuid.UUID) (*entity.Category, error)
	FindByUserID(tx *gorm.DB, userID uuid.UUID) ([]*entity.Category, error)
	Update(tx *gorm.DB, category *entity.Category) error
	Delete(tx *gorm.DB, categoryID uuid.UUID) error
	GetWithTransactions(tx *gorm.DB, categoryID uuid.UUID) (*entity.Category, error)
	CalculateTotalSpent(tx *gorm.DB, categoryID uuid.UUID) (int, error)
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ICategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(tx *gorm.DB, category *entity.Category) error {
	err := tx.Create(category).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) FindByID(tx *gorm.DB, categoryID uuid.UUID) (*entity.Category, error) {
	var category *entity.Category
	err := tx.Where("category_id = ?", categoryID).First(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) FindByUserID(tx *gorm.DB, userID uuid.UUID) ([]*entity.Category, error) {
	var categories []*entity.Category
	err := tx.Where("user_id = ?", userID).Order("created_at DESC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Update(tx *gorm.DB, category *entity.Category) error {
	err := tx.Updates(category).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) Delete(tx *gorm.DB, categoryID uuid.UUID) error {
	err := tx.Where("category_id = ?", categoryID).Delete(&entity.Category{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepository) GetWithTransactions(tx *gorm.DB, categoryID uuid.UUID) (*entity.Category, error) {
	var category *entity.Category
	err := tx.
		Preload("Transactions", func(db *gorm.DB) *gorm.DB {
			return db.Order("transaction_date DESC")
		}).
		Where("category_id = ?", categoryID).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *CategoryRepository) CalculateTotalSpent(tx *gorm.DB, categoryID uuid.UUID) (int, error) {
	var total int64
	err := tx.Model(&entity.Transaction{}).
		Where("category_id = ? AND status = ?", categoryID, "Sudah Terbayar").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return int(total), nil
}
