package model

import "github.com/google/uuid"

type UserParam struct {
	UserID uuid.UUID `json:"-"`
	Name   string    `json:"-"`
	Email  string    `json:"-"`
}

type UserRegisterParam struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email,endswith=@gmail.com"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserRegisterResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserLoginParam struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	UserID           string `json:"user_id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	MonthlyAllowance int    `json:"monthly_allowance"` 
}

type UserUpdateBudgetParam struct {
	MonthlyAllowance int `json:"monthly_allowance" binding:"required,min=0"`
}
