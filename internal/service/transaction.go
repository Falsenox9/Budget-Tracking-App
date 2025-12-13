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

type ITransactionService interface {
	CreateTransaction(param model.TransactionParam, userID uuid.UUID) (*model.TransactionResponse, error)
	GetTransaction(transactionID string, userID uuid.UUID) (*model.TransactionResponse, error)
	GetAllTransactions(userID uuid.UUID) ([]model.TransactionResponse, error)
	GetTransactionsByCategory(categoryID string, userID uuid.UUID) ([]model.TransactionResponse, error)
	UpdateTransaction(transactionID string, param model.TransactionParam, userID uuid.UUID) (*model.TransactionResponse, error)
	DeleteTransaction(transactionID string, userID uuid.UUID) error
	GetStatistics7Days(userID uuid.UUID) (*model.StatisticsResponse, error)
	GetStatisticsThisWeek(userID uuid.UUID) (*model.StatisticsResponse, error)
	GetTopCategories(userID uuid.UUID, limit int) ([]model.TopCategoryResponse, error)
}

type TransactionService struct {
	db              *gorm.DB
	transactionRepo repository.ITransactionRepository
	categoryRepo    repository.ICategoryRepository
}

func NewTransactionService(db *gorm.DB, transactionRepo repository.ITransactionRepository, categoryRepo repository.ICategoryRepository) ITransactionService {
	return &TransactionService{
		db:              db,
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
	}
}

func (s *TransactionService) CreateTransaction(param model.TransactionParam, userID uuid.UUID) (*model.TransactionResponse, error) {
	catUUID, err := uuid.Parse(param.CategoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	category, err := s.categoryRepo.FindByID(s.db, catUUID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if category.UserID != userID {
		return nil, errors.New("unauthorized access to category")
	}

	transaction := &entity.Transaction{
		TransactionID:   uuid.New(),
		Title:           param.Title,
		Amount:          param.Amount,
		Description:     param.Description,
		Status:          param.Status,
		TransactionDate: param.TransactionDate,
		PaymentDueDate:  param.PaymentDueDate,
		CategoryID:      catUUID,
		UserID:          userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = s.transactionRepo.Create(s.db, transaction)
	if err != nil {
		return nil, err
	}

	return &model.TransactionResponse{
		TransactionID:   transaction.TransactionID.String(),
		Title:           transaction.Title,
		CategoryID:      transaction.CategoryID.String(),
		CategoryName:    category.Name,
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		Status:          transaction.Status,
		TransactionDate: transaction.TransactionDate,
		PaymentDueDate:  transaction.PaymentDueDate,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}, nil
}

func (s *TransactionService) GetTransaction(transactionID string, userID uuid.UUID) (*model.TransactionResponse, error) {
	txUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	transaction, err := s.transactionRepo.FindByID(s.db, txUUID)
	if err != nil {
		return nil, err
	}

	if transaction.UserID != userID {
		return nil, errors.New("unauthorized access to transaction")
	}

	return &model.TransactionResponse{
		TransactionID:   transaction.TransactionID.String(),
		Title:           transaction.Title,
		CategoryID:      transaction.CategoryID.String(),
		CategoryName:    transaction.Category.Name,
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		Status:          transaction.Status,
		TransactionDate: transaction.TransactionDate,
		PaymentDueDate:  transaction.PaymentDueDate,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}, nil
}

func (s *TransactionService) GetAllTransactions(userID uuid.UUID) ([]model.TransactionResponse, error) {
	transactions, err := s.transactionRepo.FindByUserID(s.db, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = model.TransactionResponse{
			TransactionID:   tx.TransactionID.String(),
			Title:           tx.Title,
			CategoryID:      tx.CategoryID.String(),
			CategoryName:    tx.Category.Name,
			Amount:          tx.Amount,
			Description:     tx.Description,
			Status:          tx.Status,
			TransactionDate: tx.TransactionDate,
			PaymentDueDate:  tx.PaymentDueDate,
			CreatedAt:       tx.CreatedAt,
			UpdatedAt:       tx.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *TransactionService) GetTransactionsByCategory(categoryID string, userID uuid.UUID) ([]model.TransactionResponse, error) {
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	category, err := s.categoryRepo.FindByID(s.db, catUUID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if category.UserID != userID {
		return nil, errors.New("unauthorized access to category")
	}

	transactions, err := s.transactionRepo.FindByCategoryID(s.db, catUUID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = model.TransactionResponse{
			TransactionID:   tx.TransactionID.String(),
			Title:           tx.Title,
			CategoryID:      tx.CategoryID.String(),
			CategoryName:    category.Name,
			Amount:          tx.Amount,
			Description:     tx.Description,
			Status:          tx.Status,
			TransactionDate: tx.TransactionDate,
			PaymentDueDate:  tx.PaymentDueDate,
			CreatedAt:       tx.CreatedAt,
			UpdatedAt:       tx.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *TransactionService) UpdateTransaction(transactionID string, param model.TransactionParam, userID uuid.UUID) (*model.TransactionResponse, error) {
	txUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return nil, errors.New("invalid transaction ID")
	}

	transaction, err := s.transactionRepo.FindByID(s.db, txUUID)
	if err != nil {
		return nil, err
	}

	if transaction.UserID != userID {
		return nil, errors.New("unauthorized access to transaction")
	}

	catUUID, err := uuid.Parse(param.CategoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}

	category, err := s.categoryRepo.FindByID(s.db, catUUID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	if category.UserID != userID {
		return nil, errors.New("unauthorized access to category")
	}

	transaction.Title = param.Title
	transaction.CategoryID = catUUID
	transaction.Amount = param.Amount
	transaction.Description = param.Description
	transaction.Status = param.Status
	transaction.TransactionDate = param.TransactionDate
	transaction.PaymentDueDate = param.PaymentDueDate
	transaction.UpdatedAt = time.Now()

	err = s.transactionRepo.Update(s.db, transaction)
	if err != nil {
		return nil, err
	}

	return &model.TransactionResponse{
		TransactionID:   transaction.TransactionID.String(),
		Title:           transaction.Title,
		CategoryID:      transaction.CategoryID.String(),
		CategoryName:    category.Name,
		Amount:          transaction.Amount,
		Description:     transaction.Description,
		Status:          transaction.Status,
		TransactionDate: transaction.TransactionDate,
		PaymentDueDate:  transaction.PaymentDueDate,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}, nil
}

func (s *TransactionService) DeleteTransaction(transactionID string, userID uuid.UUID) error {
	txUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return errors.New("invalid transaction ID")
	}

	transaction, err := s.transactionRepo.FindByID(s.db, txUUID)
	if err != nil {
		return err
	}

	if transaction.UserID != userID {
		return errors.New("unauthorized access to transaction")
	}

	return s.transactionRepo.Delete(s.db, txUUID)
}

func (s *TransactionService) GetStatistics7Days(userID uuid.UUID) (*model.StatisticsResponse, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)

	transactions, err := s.transactionRepo.GetStatisticsByDateRange(s.db, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	totalSpent := 0
	for _, tx := range transactions {
		totalSpent += tx.Amount
	}

	dailyMap := make(map[string]int)
	for _, tx := range transactions {
		dateKey := tx.TransactionDate.Format("2006-01-02")
		dailyMap[dateKey] += tx.Amount
	}

	dailyData := make([]model.DailyStatistic, 0)
	var highestDay *model.DailyStatistic
	maxAmount := 0

	for date, amount := range dailyMap {
		daily := model.DailyStatistic{
			Date:   date,
			Amount: amount,
		}
		dailyData = append(dailyData, daily)

		if amount > maxAmount {
			maxAmount = amount
			highestDay = &daily
		}
	}

	categoryMap := make(map[uuid.UUID]*model.CategoryStatistic)
	for _, tx := range transactions {
		if stat, exists := categoryMap[tx.CategoryID]; exists {
			stat.TotalSpent += tx.Amount
		} else {
			categoryMap[tx.CategoryID] = &model.CategoryStatistic{
				CategoryID:   tx.CategoryID.String(),
				CategoryName: tx.Category.Name,
				TotalSpent:   tx.Amount,
			}
		}
	}

	byCategory := make([]model.CategoryStatistic, 0)
	for _, stat := range categoryMap {
		if totalSpent > 0 {
			stat.Percentage = (float64(stat.TotalSpent) / float64(totalSpent)) * 100
		}
		byCategory = append(byCategory, *stat)
	}

	return &model.StatisticsResponse{
		Period:     "Last 7 Days",
		PeriodType: "daily",
		TotalSpent: totalSpent,
		HighestDay: highestDay,
		DailyData:  dailyData,
		ByCategory: byCategory,
	}, nil
}

func (s *TransactionService) GetStatisticsThisWeek(userID uuid.UUID) (*model.StatisticsResponse, error) {
	transactions, err := s.transactionRepo.GetThisWeekTransactions(s.db, userID)
	if err != nil {
		return nil, err
	}

	totalSpent := 0
	for _, tx := range transactions {
		totalSpent += tx.Amount
	}
	weekdayNames := map[time.Weekday]string{
		time.Monday:    "Sen",
		time.Tuesday:   "Sel",
		time.Wednesday: "Rab",
		time.Thursday:  "Kam",
		time.Friday:    "Jum",
		time.Saturday:  "Sab",
		time.Sunday:    "Min",
	}

	dailyMap := make(map[string]int)
	for _, tx := range transactions {
		weekday := tx.TransactionDate.Weekday()
		dayName := weekdayNames[weekday]
		dailyMap[dayName] += tx.Amount
	}

	days := []string{"Sen", "Sel", "Rab", "Kam", "Jum", "Sab", "Min"}
	dailyData := make([]model.DailyStatistic, 0)
	var highestDay *model.DailyStatistic
	maxAmount := 0

	for _, day := range days {
		amount := dailyMap[day]
		daily := model.DailyStatistic{
			Date:   day,
			Amount: amount,
		}
		dailyData = append(dailyData, daily)

		if amount > maxAmount {
			maxAmount = amount
			highestDay = &daily
		}
	}

	categoryMap := make(map[uuid.UUID]*model.CategoryStatistic)
	for _, tx := range transactions {
		if stat, exists := categoryMap[tx.CategoryID]; exists {
			stat.TotalSpent += tx.Amount
		} else {
			categoryMap[tx.CategoryID] = &model.CategoryStatistic{
				CategoryID:   tx.CategoryID.String(),
				CategoryName: tx.Category.Name,
				TotalSpent:   tx.Amount,
			}
		}
	}

	byCategory := make([]model.CategoryStatistic, 0)
	for _, stat := range categoryMap {
		if totalSpent > 0 {
			stat.Percentage = (float64(stat.TotalSpent) / float64(totalSpent)) * 100
		}
		byCategory = append(byCategory, *stat)
	}

	return &model.StatisticsResponse{
		Period:     "This Week",
		PeriodType: "weekly",
		TotalSpent: totalSpent,
		HighestDay: highestDay,
		DailyData:  dailyData,
		ByCategory: byCategory,
	}, nil
}

func (s *TransactionService) GetTopCategories(userID uuid.UUID, limit int) ([]model.TopCategoryResponse, error) {
	topCats, err := s.transactionRepo.GetTopCategoriesBySpending(s.db, userID, limit)
	if err != nil {
		return nil, err
	}

	totalSpent := 0
	for _, cat := range topCats {
		totalSpent += cat["total_spent"].(int)
	}
	colors := []string{"#FF6B35", "#F24C88", "#4E9FE5", "#00D9A3", "#8E44AD"}

	responses := make([]model.TopCategoryResponse, len(topCats))
	for i, cat := range topCats {
		spent := cat["total_spent"].(int)
		percentage := 0.0
		if totalSpent > 0 {
			percentage = (float64(spent) / float64(totalSpent)) * 100
		}

		color := colors[i%len(colors)] 

		responses[i] = model.TopCategoryResponse{
			Rank:         i + 1,
			CategoryID:   cat["category_id"].(string),
			CategoryName: cat["category_name"].(string),
			TotalSpent:   spent,
			Percentage:   percentage,
			Color:        color,
		}
	}

	return responses, nil
}
