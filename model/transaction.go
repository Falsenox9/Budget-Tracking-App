package model

import "time"

// TransactionParam - untuk create/update transaction dari form UI
// TransactionParam - untuk create/update transaction dari form UI
type TransactionParam struct {
	Title           string     `json:"title" binding:"required"`
	CategoryID      string     `json:"category_id" binding:"required"`
	Amount          int        `json:"amount" binding:"required,min=1"`
	Description     string     `json:"description"`
	Status          string     `json:"status" binding:"required,oneof=Belum\ Terbayar Sudah\ Terbayar"`
	TransactionDate time.Time  `json:"transaction_date" binding:"required"`
	PaymentDueDate  *time.Time `json:"payment_due_date"` // Optional: kapan harus dibayar
}

// TransactionResponse - response model untuk list & detail
type TransactionResponse struct {
	TransactionID   string     `json:"transaction_id"`
	Title           string     `json:"title"`
	CategoryID      string     `json:"category_id"`
	CategoryName    string     `json:"category_name"` // Computed from relation
	Amount          int        `json:"amount"`
	Description     string     `json:"description"`
	Status          string     `json:"status"`
	TransactionDate time.Time  `json:"transaction_date"`
	PaymentDueDate  *time.Time `json:"payment_due_date"` // Tanggal jatuh tempo
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// StatisticsResponse - untuk tracking periode tertentu
type StatisticsResponse struct {
	Period     string              `json:"period"`      // "Last 7 Days", "This Week", etc
	PeriodType string              `json:"period_type"` // "weekly", "daily", "monthly"
	TotalSpent int                 `json:"total_spent"`
	HighestDay *DailyStatistic     `json:"highest_day"` // Hari dengan pengeluaran terbanyak
	DailyData  []DailyStatistic    `json:"daily_data"`
	ByCategory []CategoryStatistic `json:"by_category"`
}

// DailyStatistic - breakdown per hari
type DailyStatistic struct {
	Date   string `json:"date"` // Format: "2025-12-08" or "Sen", "Sel", etc
	Amount int    `json:"amount"`
}

// CategoryStatistic - breakdown per kategori
type CategoryStatistic struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalSpent   int     `json:"total_spent"`
	Percentage   float64 `json:"percentage"` // Percentage of total
}

// TopCategoryResponse - untuk ranking top categories
type TopCategoryResponse struct {
	Rank         int     `json:"rank"`
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalSpent   int     `json:"total_spent"`
	Percentage   float64 `json:"percentage"`
	Color        string  `json:"color"` // For UI display
}
