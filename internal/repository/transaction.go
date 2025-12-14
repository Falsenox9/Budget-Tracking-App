package repository

import (
	"budget-tracking-app/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	Create(tx *gorm.DB, transaction *entity.Transaction) error
	FindByID(tx *gorm.DB, transactionID uuid.UUID) (*entity.Transaction, error)
	FindByUserID(tx *gorm.DB, userID uuid.UUID) ([]*entity.Transaction, error)
	FindByCategoryID(tx *gorm.DB, categoryID uuid.UUID) ([]*entity.Transaction, error)
	Update(tx *gorm.DB, transaction *entity.Transaction) error
	Delete(tx *gorm.DB, transactionID uuid.UUID) error
	GetLast7Days(tx *gorm.DB, userID uuid.UUID) ([]*entity.Transaction, error)
	GetThisWeekTransactions(tx *gorm.DB, userID uuid.UUID) ([]*entity.Transaction, error)
	GetStatisticsByDateRange(tx *gorm.DB, userID uuid.UUID, startDate, endDate time.Time) ([]*entity.Transaction, error)
	SumByCategoryInPeriod(tx *gorm.DB, userID uuid.UUID, startDate, endDate time.Time) (map[uuid.UUID]int, error)
	GetTopCategoriesBySpending(tx *gorm.DB, userID uuid.UUID, limit int) ([]map[string]interface{}, error)
}

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *gorm.DB, transaction *entity.Transaction) error {
	err := tx.Create(transaction).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) FindByID(tx *gorm.DB, transactionID uuid.UUID) (*entity.Transaction, error) {
	var transaction *entity.Transaction
	err := tx.
		Preload("Category").
		Where("transaction_id = ?", transactionID).
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *TransactionRepository) FindByUserID(tx *gorm.DB, userID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := tx.
		Preload("Category").
		Where("user_id = ?", userID).
		Order("transaction_date DESC, created_at DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) FindByCategoryID(tx *gorm.DB, categoryID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := tx.
		Where("category_id = ?", categoryID).
		Order("transaction_date DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) Update(tx *gorm.DB, transaction *entity.Transaction) error {
	err := tx.Updates(transaction).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) Delete(tx *gorm.DB, transactionID uuid.UUID) error {
	err := tx.Where("transaction_id = ?", transactionID).Delete(&entity.Transaction{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) GetLast7Days(tx *gorm.DB, userID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	err := tx.
		Preload("Category").
		Where("user_id = ? AND transaction_date >= ?", userID, sevenDaysAgo).
		Order("transaction_date DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetThisWeekTransactions(tx *gorm.DB, userID uuid.UUID) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction

	now := time.Now()
	weekday := now.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -int(weekday)+1)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	endOfWeek := startOfWeek.AddDate(0, 0, 6)
	endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 999999999, endOfWeek.Location())

	err := tx.
		Preload("Category").
		Where("user_id = ? AND transaction_date BETWEEN ? AND ? AND status = ?",
			userID, startOfWeek, endOfWeek, "Sudah Terbayar").
		Order("transaction_date ASC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetStatisticsByDateRange(tx *gorm.DB, userID uuid.UUID, startDate, endDate time.Time) ([]*entity.Transaction, error) {
	var transactions []*entity.Transaction
	err := tx.
		Preload("Category").
		Where("user_id = ? AND transaction_date BETWEEN ? AND  ? AND status = ?",
			userID, startDate, endDate, "Sudah Terbayar").
		Order("transaction_date DESC").
		Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) SumByCategoryInPeriod(tx *gorm.DB, userID uuid.UUID, startDate, endDate time.Time) (map[uuid.UUID]int, error) {
	type Result struct {
		CategoryID uuid.UUID
		Total      int64
	}

	var results []Result
	err := tx.Model(&entity.Transaction{}).
		Select("category_id, SUM(amount) as total").
		Where("user_id = ? AND transaction_date BETWEEN ? AND ? AND status = ?",
			userID, startDate, endDate, "Sudah Terbayar").
		Group("category_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	sumMap := make(map[uuid.UUID]int)
	for _, result := range results {
		sumMap[result.CategoryID] = int(result.Total)
	}

	return sumMap, nil
}

func (r *TransactionRepository) GetTopCategoriesBySpending(tx *gorm.DB, userID uuid.UUID, limit int) ([]map[string]interface{}, error) {
	type Result struct {
		CategoryID   uuid.UUID
		CategoryName string
		TotalSpent   int64
	}

	var results []Result
	err := tx.Table("transactions").
		Select("transactions.category_id, categories.name as category_name, SUM(transactions.amount) as total_spent").
		Joins("JOIN categories ON transactions.category_id = categories.category_id").
		Where("transactions.user_id = ? AND transactions.status = ?", userID, "Sudah Terbayar").
		Group("transactions.category_id, categories.name").
		Order("total_spent DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	topCategories := make([]map[string]interface{}, len(results))
	for i, result := range results {
		topCategories[i] = map[string]interface{}{
			"category_id":   result.CategoryID.String(),
			"category_name": result.CategoryName,
			"total_spent":   int(result.TotalSpent),
		}
	}

	return topCategories, nil
}
