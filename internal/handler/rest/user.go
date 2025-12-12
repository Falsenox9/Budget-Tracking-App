package rest

import (
	"budget-tracking-app/entity"
	"budget-tracking-app/model"
	"budget-tracking-app/pkg/response"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (r *Rest) RegisterUser(c *gin.Context) {
	var param model.UserRegisterParam

	err := c.ShouldBindJSON(&param)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errosMessages := make([]string, 0)
			for _, e := range validationErrors {
				msg := fmt.Sprintf("Field '%s' error: %s", e.Field(), e.Tag())
				errosMessages = append(errosMessages, msg)
			}
			response.Error(c, http.StatusBadRequest, "invalid validation format", fmt.Errorf("%v", errosMessages))
			return
		}
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	reps, err := r.service.UserService.RegisterUser(param)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to register user", err)
		return
	}

	response.Success(c, http.StatusOK, "user registered successfully", reps)

}

func (r *Rest) LoginUser(c *gin.Context) {
	var param model.UserLoginParam

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	reps, err := r.service.UserService.LoginUser(param)
	if err != nil {
		if err.Error() == "email or password is wrong" {
			response.Error(c, http.StatusUnauthorized, "email or password is wrong", err)
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to login user", err)
		return
	}

	response.Success(c, http.StatusOK, "user logged in successfully", reps)
}

func (r *Rest) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	userEntity, ok := user.(*entity.User)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "failed to get user", nil)
		return
	}

	resp, err := r.service.UserService.GetUserProfile(model.UserParam{
		UserID: userEntity.UserID,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get user profile", err)
		return
	}

	response.Success(c, http.StatusOK, "user profile retrieved successfully", resp)
}

func (r *Rest) UpdateMonthlyBudget(c *gin.Context) {
	var param model.UserUpdateBudgetParam

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	resp, err := r.service.UserService.UpdateMonthlyBudget(userUUID, param.MonthlyAllowance)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update monthly budget", err)
		return
	}

	response.Success(c, http.StatusOK, "monthly allowance updated successfully", resp)
}
