package mariadb

import (
	"budget-tracking-app/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Migrate in order: base tables first, then tables with foreign keys
	// 1. User table (no foreign keys)
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return err
	}

	// 2. Category table (has UserID foreign key)
	if err := db.AutoMigrate(&entity.Category{}); err != nil {
		return err
	}

	// 3. Transaction table (has both UserID and CategoryID foreign keys)
	if err := db.AutoMigrate(&entity.Transaction{}); err != nil {
		return err
	}

	return nil
}
