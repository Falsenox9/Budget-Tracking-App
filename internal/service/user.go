package service

import (
	"budget-tracking-app/entity"
	"budget-tracking-app/internal/repository"
	"budget-tracking-app/model"
	"budget-tracking-app/pkg/bcrypt"
	"budget-tracking-app/pkg/jwt"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IUserService interface {
	RegisterUser(param model.UserRegisterParam) (*model.UserRegisterResponse, error)
	LoginUser(param model.UserLoginParam) (*model.UserLoginResponse, error)
	GetUserProfile(param model.UserParam) (*model.UserResponse, error)
	UpdateMonthlyBudget(userID uuid.UUID, amount int) (*model.UserResponse, error)
}

type UserService struct {
	userRepo repository.IUserRepository
	bcrypt   bcrypt.Interface
	jwtAuth  jwt.Interface
}

func NewUserService(userRepo repository.IUserRepository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) IUserService {
	return &UserService{
		userRepo: userRepo,
		bcrypt:   bcrypt,
		jwtAuth:  jwtAuth,
	}
}

func (s *UserService) RegisterUser(param model.UserRegisterParam) (*model.UserRegisterResponse, error) {
	// Hash password
	hashedPassword, err := s.bcrypt.GenerateFromPassword(param.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user entity
	user := &entity.User{
		UserID:           uuid.New(),
		Name:             param.Name,
		Email:            param.Email,
		Password:         hashedPassword,
		MonthlyAllowance: 0, // Default value
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Save to database
	err = s.userRepo.CreateUser(nil, user)
	if err != nil {
		// Check for duplicate email error (MySQL Error 1062)
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("email already registered")
		}
		return nil, err
	}

	return &model.UserRegisterResponse{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) LoginUser(param model.UserLoginParam) (*model.UserLoginResponse, error) {
	// Find user by email (using name field as email for now based on UserLoginParam)
	user, err := s.userRepo.GetUser(model.UserParam{
		Email: param.Email,
	})
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare password
	err = s.bcrypt.CompareAndHashPassword(user.Password, param.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.jwtAuth.CreateJWTToken(user.UserID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.UserLoginResponse{
		Token: token,
	}, nil
}

func (s *UserService) GetUserProfile(param model.UserParam) (*model.UserResponse, error) {
	user, err := s.userRepo.GetUser(param)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		UserID:           user.UserID.String(),
		Name:             user.Name,
		Email:            user.Email,
		MonthlyAllowance: user.MonthlyAllowance,
	}, nil
}

func (s *UserService) UpdateMonthlyBudget(userID uuid.UUID, amount int) (*model.UserResponse, error) {
	// Get user
	user, err := s.userRepo.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update monthly allowance
	user.MonthlyAllowance = amount
	user.UpdatedAt = time.Now()

	// Save to database
	err = s.userRepo.UpdateUser(nil, user)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		UserID:           user.UserID.String(),
		Name:             user.Name,
		Email:            user.Email,
		MonthlyAllowance: user.MonthlyAllowance,
	}, nil
}
