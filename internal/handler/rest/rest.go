package rest

import (
	"budget-tracking-app/internal/service"
	"budget-tracking-app/pkg/middleware"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router     *gin.Engine
	service    *service.Service
	middleware middleware.Interface
}

func NewRest(service *service.Service, middleware middleware.Interface) *Rest {
	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Serve Static Files
	router.Static("/uploads", "./uploads")

	return &Rest{
		router:     router,
		service:    service,
		middleware: middleware,
	}
}

func (r *Rest) MountEndPoint() {
	baseURL := r.router.Group("/api/v2")

	// User Endpoint
	auth := baseURL.Group("/auth")
	auth.POST("/register", r.RegisterUser)
	auth.POST("/login", r.LoginUser)

	user := baseURL.Group("/user")
	user.Use(r.middleware.AuthenticateUser)
	user.PUT("/budget", r.UpdateMonthlyBudget)

	// Profile Endpoint (At root as requested)
	baseURL.GET("/profile", r.middleware.AuthenticateUser, r.GetProfile)
	// user.GET("/get-all-users", r.GetAllUsers)
	// user.GET("/get-user-by-id/:id", r.GetUserByID)
	// user.PUT("/update-user/:id", r.UpdateUser)
	// user.DELETE("/delete-user/:id", r.DeleteUser)

	// Category Endpoint
	category := baseURL.Group("/categories")
	category.Use(r.middleware.AuthenticateUser)
	category.POST("", r.CreateCategory)
	category.GET("", r.GetAllCategories)
	category.GET("/:id", r.GetCategory)
	category.GET("/:id/transactions", r.GetCategoryWithTransactions)
	category.PUT("/:id", r.UpdateCategory)
	// category.DELETE("/:id", r.DeleteCategory)

	// Transaction Endpoint
	transaction := baseURL.Group("/transactions")
	transaction.Use(r.middleware.AuthenticateUser)
	transaction.POST("", r.CreateTransaction)
	transaction.GET("", r.GetAllTransactions)
	transaction.GET("/:id", r.GetTransaction)
	transaction.PUT("/:id", r.UpdateTransaction)
	// transaction.DELETE("/:id", r.DeleteTransaction)

	// Statistics Endpoint
	statistics := baseURL.Group("/statistics")
	statistics.Use(r.middleware.AuthenticateUser)
	statistics.GET("/7days", r.GetStatistics7Days)
	statistics.GET("/weekly", r.GetStatisticsThisWeek)
	statistics.GET("/top-categories", r.GetTopCategories)

	// Plans Endpoint (Khusus Rencana)
	plans := baseURL.Group("/plans")
	plans.Use(r.middleware.AuthenticateUser)
	plans.POST("", r.CreatePlan)
	plans.GET("", r.GetPlans) // NEW

	// Verify Password Endpoint
	baseURL.POST("/verify-password", r.middleware.AuthenticateUser, r.VerifyPassword)
	baseURL.PUT("/change-password", r.middleware.AuthenticateUser, r.ChangePassword)

	// Profile & Balance & Notifications (New Requirements)
	baseURL.PUT("/profile", r.middleware.AuthenticateUser, r.UpdateProfile)
	baseURL.POST("/profile/photo", r.middleware.AuthenticateUser, r.UploadAvatar)
	baseURL.POST("/balance", r.middleware.AuthenticateUser, r.UpdateBalance)
	baseURL.GET("/balance", r.middleware.AuthenticateUser, r.GetBalance) // NEW
	baseURL.GET("/notifications", r.middleware.AuthenticateUser, r.GetNotifications)
}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
