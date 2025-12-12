package repository

import "gorm.io/gorm"

type Repository struct {
	UserRepository        IUserRepository
	CategoryRepository    ICategoryRepository
	TransactionRepository ITransactionRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:        NewUserRepository(db),
		CategoryRepository:    NewCategoryRepository(db),
		TransactionRepository: NewTransactionRepository(db),
	}
}
