package middleware

import (
	"budget-tracking-app/internal/service"
	"budget-tracking-app/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type Interface interface {
	AuthenticateUser(c *gin.Context)
}

type middleware struct {
	service *service.Service
	jwtAuth jwt.Interface
}

func Init(service *service.Service, jwtAuth jwt.Interface) Interface {
	return &middleware{
		service: service,
		jwtAuth: jwtAuth,
	}
}