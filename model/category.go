package model

import "time"

// CategoryParam - untuk create/update category
type CategoryParam struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	Type           string `json:"type" binding:"required,oneof=Pengeluaran Rencana"`
	AllowanceLimit int    `json:"allowance_limit" binding:"required,min=0"`
}

// CategoryResponse - response dengan computed fields untuk UI
type CategoryResponse struct {
	CategoryID     string    `json:"category_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	AllowanceLimit int       `json:"allowance_limit"`
	TotalSpent     int       `json:"total_spent"` // Computed: total dari transactions
	Remaining      int       `json:"remaining"`   // Computed: limit - spent
	Progress       float64   `json:"progress"`    // Computed: (spent/limit) * 100
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CategoryWithTransactions - untuk detail dengan list transaksi
type CategoryWithTransactions struct {
	CategoryResponse
	Transactions []TransactionResponse `json:"transactions"`
}
