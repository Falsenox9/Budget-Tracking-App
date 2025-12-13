package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID           uuid.UUID `json:"user_id" gorm:"type:varchar(36);primaryKey"`
	Name             string    `json:"name" gorm:"type:varchar(255);not null"`
	Email            string    `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password         string    `json:"password" gorm:"type:varchar(255);not null"`
	MonthlyAllowance int       `json:"monthly_allowance" gorm:"type:int;not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	
	// Transactions []Transaction `json:"transactions" gorm:"foreignKey:UserID;references:UserID"`
	// Categories   []Category    `json:"categories" gorm:"foreignKey:UserID;references:UserID"`
}
