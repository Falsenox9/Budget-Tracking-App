package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	CategoryID     uuid.UUID `json:"category_id" gorm:"type:varchar(36);primaryKey"`
	Name           string    `json:"name" gorm:"type:varchar(100);not null"`
	Description    string    `json:"description" gorm:"type:text"`
	Type           string    `json:"type" gorm:"type:enum('Pengeluaran','Rencana');not null"`
	AllowanceLimit int       `json:"allowance_limit" gorm:"type:int;default:0"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	UserID uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null"`

	User User `json:"user" gorm:"foreignKey:UserID;references:UserID"`
}
