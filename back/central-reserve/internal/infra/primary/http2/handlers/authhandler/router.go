package authhandler

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configura las rutas de autenticación
func SetupRoutes(router *gin.Engine, handler IAuthHandler) {
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login", handler.LoginHandler)
	}
}
