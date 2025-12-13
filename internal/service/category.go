package service

import (
	"budget-tracking-app/entity"
	"budget-tracking-app/internal/repository"
	"budget-tracking-app/model"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICategoryService interface {
	CreateCategory(param model.CategoryParam, userID uuid.UUID) (*model.CategoryResponse, error)
	GetCategory(categoryID string, userID uuid.UUID) (*model.CategoryResponse, error)
	GetCategoryWithTransactions(categoryID string, userID uuid.UUID) (*model.CategoryWithTransactions, error)
	GetAllCategories(userID uuid.UUID) ([]model.CategoryResponse, error)
	UpdateCategory(categoryID string, param model.CategoryParam, userID uuid.UUID) (*model.CategoryResponse, error)
	DeleteCategory(categoryID string, userID uuid.UUID) error
}

type CategoryService struct {
	db              *gorm.DB
	categoryRepo    repository.ICategoryRepository
	transactionRepo repository.ITransactionRepository
}

func NewCategoryService(db *gorm.DB, categoryRepo repository.ICategoryRepository, transactionRepo repository.ITransactionRepository) ICategoryService {
	return &CategoryService{
		db:              db,
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *CategoryService) CreateCategory(param model.CategoryParam, userID uuid.UUID) (*model.CategoryResponse, error) {
	category := &entity.Category{
		CategoryID:     uuid.New(),
		Name:           param.Name,
		Description:    param.Description,
		Type:           param.Type,
		AllowanceLimit: param.AllowanceLimit,
		UserID:         userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := s.categoryRepo.Create(s.db, category)
	if err != nil {
		return nil, err
	}

	// Calculate progress (initially 0 spent)
	return s.calculateCategoryProgress(category, 0), nil
}

func (s *CategoryService) GetCategory(categoryID string, userID uuid.UUID) (*model.CategoryResponse, error) {
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	category, err := s.categoryRepo.FindByID(s.db, catUUID)
	if err != nil {
		return nil, err
	}

	if category.UserID != userID {
		return nil, errors.New("unauthorized access to category")
	}

	totalSpent, err := s.categoryRepo.CalculateTotalSpent(s.db, catUUID)
	if err != nil {
		return nil, err
	}

	return s.calculateCategoryProgress(category, totalSpent), nil
}

func (s *CategoryService) GetCategoryWithTransactions(categoryID string, userID uuid.UUID) (*model.CategoryWithTransactions, error) {
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	category, err := s.categoryRepo.FindByID(s.db, catUUID)
	if err != nil {
		return nil, err
	}

	if category.UserID != userID {
		return nil, errors.New("unauthorized access to category")
	}

	totalSpent, err := s.categoryRepo.CalculateTotalSpent(s.db, catUUID)
	if err != nil {
		return nil, err
	}

	txList, err := s.transactionRepo.FindByCategoryID(s.db, catUUID)
	if err != nil {
		return nil, err
	}

	transactions := make([]model.TransactionResponse, len(txList))
	for i, tx := range txList {
		transactions[i] = model.TransactionResponse{
			TransactionID:   tx.TransactionID.String(),
			Title:           tx.Title,
			CategoryID:      tx.CategoryID.String(),
			CategoryName:    category.Name,
			Amount:          tx.Amount,
			Description:     tx.Description,
			Status:          tx.Status,
			TransactionDate: tx.TransactionDate,
			CreatedAt:       tx.CreatedAt,
			UpdatedAt:       tx.UpdatedAt,
		}
	}

	return &model.CategoryWithTransactions{
		CategoryResponse: *s.calculateCategoryProgress(category, totalSpent),
		Transactions:     transactions,
	}, nil
}

func (s *CategoryService) GetAllCategories(userID uuid.UUID) ([]model.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindByUserID(s.db, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.CategoryResponse, len(categories))
	for i, category := range categories {
		totalSpent, err := s.categoryRepo.CalculateTotalSpent(s.db, category.CategoryID)
		if err != nil {
			return nil, err
		}
		responses[i] = *s.calculateCategoryProgress(category, totalSpent)
	}

	return responses, nil
}

func (s *CategoryService) UpdateCategory(categoryID string, param model.CategoryParam, userID uuid.UUID) (*model.CategoryResponse, error) {
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	category, err := s.categoryRepo.FindByID(s.db, catUUID)
	if err != nil {
		return nil, err
	}

	if category.UserID != userID {
		return nil, errors.New("unauthorized access to category")
	}

	category.Name = param.Name
	category.Description = param.Description
	category.Type = param.Type
	category.AllowanceLimit = param.AllowanceLimit
	category.UpdatedAt = time.Now()

	err = s.categoryRepo.Update(s.db, category)
	if err != nil {
		return nil, err
	}

	totalSpent, err := s.categoryRepo.CalculateTotalSpent(s.db, catUUID)
	if err != nil {
		return nil, err
	}

	return s.calculateCategoryProgress(category, totalSpent), nil
}

func (s *CategoryService) DeleteCategory(categoryID string, userID uuid.UUID) error {
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return errors.New("invalid category ID")
	}

	category, err := s.categoryRepo.FindByID(s.db, catUUID)
	if err != nil {
		return err
	}

	if category.UserID != userID {
		return errors.New("unauthorized access to category")
	}

	return s.categoryRepo.Delete(s.db, catUUID)
}

func (s *CategoryService) calculateCategoryProgress(category *entity.Category, totalSpent int) *model.CategoryResponse {
	remaining := category.AllowanceLimit - totalSpent
	if remaining < 0 {
		remaining = 0
	}

	var progress float64
	if category.AllowanceLimit > 0 {
		progress = (float64(totalSpent) / float64(category.AllowanceLimit)) * 100
		if progress > 100 {
			progress = 100
		}
	}

	return &model.CategoryResponse{
		CategoryID:     category.CategoryID.String(),
		Name:           category.Name,
		Description:    category.Description,
		Type:           category.Type,
		AllowanceLimit: category.AllowanceLimit,
		TotalSpent:     totalSpent,
		Remaining:      remaining,
		Progress:       progress,
		CreatedAt:      category.CreatedAt,
		UpdatedAt:      category.UpdatedAt,
	}
}
