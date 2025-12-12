package service

import (
	"budget-tracking-app/internal/repository"
	"budget-tracking-app/pkg/bcrypt"
	"budget-tracking-app/pkg/jwt"

	"gorm.io/gorm"
)

type Service struct {
	UserService        IUserService
	CategoryService    ICategoryService
	TransactionService ITransactionService
}

func NewService(db *gorm.DB, repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) *Service {
	return &Service{
		UserService:        NewUserService(repository.UserRepository, bcrypt, jwtAuth),
		CategoryService:    NewCategoryService(db, repository.CategoryRepository, repository.TransactionRepository),
		TransactionService: NewTransactionService(db, repository.TransactionRepository, repository.CategoryRepository),
	}
}
