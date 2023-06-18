package services

import (
	"github.com/gin-gonic/gin"
)

type ClientService interface {
	RegisterRoutes(*gin.Engine) error
}
