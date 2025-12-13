package rest

import (
	"budget-tracking-app/model"
	"budget-tracking-app/pkg/response"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) CreateTransaction(c *gin.Context) {
	var param model.TransactionParam

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

	resp, err := r.service.TransactionService.CreateTransaction(param, userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create transaction", err)
		return
	}

	response.Success(c, http.StatusCreated, "transaction created successfully", resp)
}

func (r *Rest) GetTransaction(c *gin.Context) {
	transactionID := c.Param("id")

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

	resp, err := r.service.TransactionService.GetTransaction(transactionID, userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get transaction", err)
		return
	}

	response.Success(c, http.StatusOK, "transaction retrieved successfully", resp)
}

func (r *Rest) GetAllTransactions(c *gin.Context) {
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

	resp, err := r.service.TransactionService.GetAllTransactions(userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get transactions", err)
		return
	}

	response.Success(c, http.StatusOK, "transactions retrieved successfully", resp)
}

func (r *Rest) UpdateTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	var param model.TransactionParam

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

	resp, err := r.service.TransactionService.UpdateTransaction(transactionID, param, userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update transaction", err)
		return
	}

	response.Success(c, http.StatusOK, "transaction updated successfully", resp)
}

// func (r *Rest) DeleteTransaction(c *gin.Context) {
// 	transactionID := c.Param("id")
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
// 	err = r.service.TransactionService.DeleteTransaction(transactionID, userUUID)
// 	if err != nil {
// 		response.Error(c, http.StatusInternalServerError, "failed to delete transaction", err)
// 		return
// 	}
//
// 	response.Success(c, http.StatusOK, "transaction deleted successfully", nil)
// }

func (r *Rest) GetStatistics7Days(c *gin.Context) {
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

	resp, err := r.service.TransactionService.GetStatistics7Days(userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get statistics", err)
		return
	}

	response.Success(c, http.StatusOK, "statistics retrieved successfully", resp)
}

func (r *Rest) GetStatisticsThisWeek(c *gin.Context) {
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

	resp, err := r.service.TransactionService.GetStatisticsThisWeek(userUUID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get weekly statistics", err)
		return
	}

	response.Success(c, http.StatusOK, "weekly statistics retrieved successfully", resp)
}

func (r *Rest) GetTopCategories(c *gin.Context) {
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

	limit := 5

	resp, err := r.service.TransactionService.GetTopCategories(userUUID, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get top categories", err)
		return
	}

	response.Success(c, http.StatusOK, "top categories retrieved successfully", resp)
}
