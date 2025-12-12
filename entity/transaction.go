package entity

import (
	"github.com/google/uuid"

	"time"
)

type Transaction struct {
	TransactionID   uuid.UUID `json:"transaction_id" gorm:"type:varchar(36);primaryKey"`
	Title           string    `json:"title" gorm:"type:varchar(200);not null"`
	Amount          int       `json:"amount" gorm:"type:int;not null"`
	Description     string    `json:"description" gorm:"type:text"`
	Status          string    `json:"status" gorm:"type:enum('Belum Terbayar','Sudah Terbayar');default:'Belum Terbayar'"`
	TransactionDate time.Time `json:"transaction_date" gorm:"type:date;not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	CategoryID uuid.UUID `json:"category_id" gorm:"type:varchar(36);not null"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null"`

	Category Category `json:"category" gorm:"foreignKey:CategoryID;references:CategoryID"`
	User     User     `json:"user" gorm:"foreignKey:UserID;references:UserID"`
}
