package service

import (
	"budget-tracking-app/entity"
	"budget-tracking-app/internal/repository"
	"budget-tracking-app/model"
	"budget-tracking-app/pkg/bcrypt"
	"budget-tracking-app/pkg/jwt"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IUserService interface {
	RegisterUser(param model.UserRegisterParam) (*model.UserRegisterResponse, error)
	LoginUser(param model.UserLoginParam) (*model.UserLoginResponse, error)
	GetUserProfile(param model.UserParam) (*model.UserResponse, error)
	UpdateMonthlyBudget(userID uuid.UUID, amount int) (*model.UserResponse, error)
	VerifyPassword(userID uuid.UUID, password string) error
	UpdatePassword(userID uuid.UUID, newPassword string) error
	UpdateProfile(userID uuid.UUID, param model.UpdateProfileParam) (*model.UserResponse, error)
	UpdateBalance(userID uuid.UUID, param model.UpdateBalanceParam) (*model.UserResponse, error)
	UploadAvatar(userID uuid.UUID, file *multipart.FileHeader) (string, error)
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
	hashedPassword, err := s.bcrypt.GenerateFromPassword(param.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user := &entity.User{
		UserID:           uuid.New(),
		Name:             param.Name,
		Email:            param.Email,
		Password:         hashedPassword,
		MonthlyAllowance: 0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err = s.userRepo.CreateUser(nil, user)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("email already registered")
		}
		return nil, err
	}

	token, err := s.jwtAuth.CreateJWTToken(user.UserID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.UserRegisterResponse{
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}

func (s *UserService) LoginUser(param model.UserLoginParam) (*model.UserLoginResponse, error) {
	user, err := s.userRepo.GetUser(model.UserParam{
		Email: param.Email,
	})
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	err = s.bcrypt.CompareAndHashPassword(user.Password, param.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtAuth.CreateJWTToken(user.UserID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.UserLoginResponse{
		Token: token,
		User: model.UserResponse{
			UserID:           user.UserID.String(),
			Name:             user.Name,
			Email:            user.Email,
			MonthlyAllowance: user.MonthlyAllowance,
			Avatar:           user.Avatar,
		},
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
	user, err := s.userRepo.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.MonthlyAllowance = amount
	user.UpdatedAt = time.Now()
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

func (s *UserService) VerifyPassword(userID uuid.UUID, password string) error {
	user, err := s.userRepo.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return errors.New("user not found")
	}

	err = s.bcrypt.CompareAndHashPassword(user.Password, password)
	if err != nil {
		return errors.New("invalid password")
	}

	return nil
}

func (s *UserService) UpdatePassword(userID uuid.UUID, newPassword string) error {
	hashedPassword, err := s.bcrypt.GenerateFromPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user, err := s.userRepo.GetUser(model.UserParam{
		UserID: userID,
	})
	if err != nil {
		return errors.New("user not found")
	}

	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	err = s.userRepo.UpdateUser(nil, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, param model.UpdateProfileParam) (*model.UserResponse, error) {
	user, err := s.userRepo.GetUser(model.UserParam{UserID: userID})
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = param.FullName
	user.Email = param.Email
	user.UpdatedAt = time.Now()

	err = s.userRepo.UpdateUser(nil, user)
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		UserID:           user.UserID.String(),
		Name:             user.Name,
		Email:            user.Email,
		MonthlyAllowance: user.MonthlyAllowance,
		Avatar:           user.Avatar,
	}, nil
}

func (s *UserService) UpdateBalance(userID uuid.UUID, param model.UpdateBalanceParam) (*model.UserResponse, error) {
	return s.UpdateMonthlyBudget(userID, param.Amount)
}

func (s *UserService) UploadAvatar(userID uuid.UUID, file *multipart.FileHeader) (string, error) {
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	urlPath := strings.ReplaceAll(filePath, "\\", "/")
	fullURL := fmt.Sprintf("/%s", urlPath)

	user, err := s.userRepo.GetUser(model.UserParam{UserID: userID})
	if err != nil {
		return "", errors.New("user not found")
	}

	user.Avatar = fullURL
	user.UpdatedAt = time.Now()

	if err := s.userRepo.UpdateUser(nil, user); err != nil {
		return "", err
	}

	return fullURL, nil
}
