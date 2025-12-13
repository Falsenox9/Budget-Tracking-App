package rest

import (
	"budget-tracking-app/model"
	"budget-tracking-app/pkg/response"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) CreateCategory(c *gin.Context) {
	var param model.CategoryParam

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	resp, err := r.service.CategoryService.CreateCategory(param, userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create category", err)
		return
	}

	response.Success(c, http.StatusCreated, "category created successfully", resp)
}

func (r *Rest) GetCategory(c *gin.Context) {
	categoryID := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	resp, err := r.service.CategoryService.GetCategory(categoryID, userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get category", err)
		return
	}

	response.Success(c, http.StatusOK, "category retrieved successfully", resp)
}

func (r *Rest) GetAllCategories(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	resp, err := r.service.CategoryService.GetAllCategories(userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get categories", err)
		return
	}

	response.Success(c, http.StatusOK, "categories retrieved successfully", resp)
}

func (r *Rest) GetCategoryWithTransactions(c *gin.Context) {
	categoryID := c.Param("id")

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	resp, err := r.service.CategoryService.GetCategoryWithTransactions(categoryID, userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get category with transactions", err)
		return
	}

	response.Success(c, http.StatusOK, "category with transactions retrieved successfully", resp)
}

func (r *Rest) UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	var param model.CategoryParam

	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated", nil)
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
		return
	}

	resp, err := r.service.CategoryService.UpdateCategory(categoryID, param, userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update category", err)
		return
	}

	response.Success(c, http.StatusOK, "category updated successfully", resp)
}

// func (r *Rest) DeleteCategory(c *gin.Context) {
// 	categoryID := c.Param("id")
//
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		response.Error(c, http.StatusUnauthorized, "user not authenticated", nil)
// 		return
// 	}
//
// 	userUUID, err := uuid.Parse(userID.(string))
// 	if err != nil {
// 		response.Error(c, http.StatusBadRequest, "invalid user ID", err)
// 		return
// 	}
//
// 	err = r.service.CategoryService.DeleteCategory(categoryID, userUUID)
// 	if err != nil {
// 		response.Error(c, http.StatusInternalServerError, "failed to delete category", err)
// 		return
// 	}
//
// 	response.Success(c, http.StatusOK, "category deleted successfully", nil)
// }
