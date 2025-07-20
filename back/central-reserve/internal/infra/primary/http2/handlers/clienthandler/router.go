package clienthandler

import (
	"central_reserve/internal/infra/primary/http2/middleware"
	"central_reserve/internal/pkg/jwt"
	"central_reserve/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra las rutas del handler de clientes
func RegisterRoutes(v1Group *gin.RouterGroup, handler IClientHandler, jwtService *jwt.JWTService, logger log.ILogger) {
	// Crear el subgrupo /clients dentro de /api/v1
	clients := v1Group.Group("/clients")
	// Aplicar middleware de autenticación a todas las rutas de clientes
	clients.Use(middleware.AuthMiddleware(jwtService, logger))
	{
		clients.GET("", handler.GetClientsHandler)
		clients.GET("/:id", handler.GetClientByIDHandler)
		clients.POST("", handler.CreateClientHandler)
		clients.PUT("/:id", handler.UpdateClientHandler)
		clients.DELETE("/:id", handler.DeleteClientHandler)
	}
}
