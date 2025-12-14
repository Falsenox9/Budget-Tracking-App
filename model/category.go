package model

import "time"

type CategoryParam struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	Type           string `json:"type"`                            
	AllowanceLimit int    `json:"budget" binding:"required,min=0"` 
	Color          string `json:"color"`
}

type CategoryResponse struct {
	CategoryID     string    `json:"category_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	AllowanceLimit int       `json:"budget"` 
	Color          string    `json:"color"`
	TotalSpent     int       `json:"total_spent"`
	Remaining      int       `json:"remaining"`
	Progress       float64   `json:"progress"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CategoryWithTransactions struct {
	CategoryResponse
	Transactions []TransactionResponse `json:"transactions"`
}
