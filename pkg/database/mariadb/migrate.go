package mariadb

import (
	"budget-tracking-app/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Category{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Transaction{}); err != nil {
		return err
	}

	return nil
}
