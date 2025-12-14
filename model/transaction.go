package model

import "time"

type TransactionParam struct {
	Title           string     `json:"title" binding:"required"`
	CategoryID      string     `json:"category_id"` 
	CategoryName    string     `json:"category"`    
	Amount          int        `json:"amount" binding:"required,min=1"`
	Description     string     `json:"description"`
	Status          string     `json:"status"`                  
	TransactionDate string     `json:"date" binding:"required"` 
	PaymentDueDate  *time.Time `json:"payment_due_date"`
}

type PlanParam struct {
	Title       string `json:"title" binding:"required"`
	Amount      int    `json:"amount" binding:"required,min=1"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate" binding:"required"` 
	Category    string `json:"category" binding:"required"` 
	Completed   bool   `json:"completed"`
}

type TransactionResponse struct {
	TransactionID   string     `json:"transaction_id"`
	Title           string     `json:"title"`
	CategoryID      string     `json:"category_id"`
	CategoryName    string     `json:"category_name"` 
	Amount          int        `json:"amount"`
	Description     string     `json:"description"`
	Status          string     `json:"status"`
	TransactionDate time.Time  `json:"transaction_date"`
	PaymentDueDate  *time.Time `json:"payment_due_date"` 
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type StatisticsResponse struct {
	Period     string              `json:"period"`      
	PeriodType string              `json:"period_type"` 
	TotalSpent int                 `json:"total_spent"`
	HighestDay *DailyStatistic     `json:"highest_day"` 
	DailyData  []DailyStatistic    `json:"daily_data"`
	ByCategory []CategoryStatistic `json:"by_category"`
}

type DailyStatistic struct {
	Date   string `json:"date"` 
	Amount int    `json:"amount"`
}

type CategoryStatistic struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalSpent   int     `json:"total_spent"`
	Percentage   float64 `json:"percentage"` 
}

type TopCategoryResponse struct {
	Rank         int     `json:"rank"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalSpent   int     `json:"total_spent"`
	Percentage   float64 `json:"percentage"`
	Color        string  `json:"color"` 
}
