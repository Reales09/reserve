package reservehandler

import (
	"central_reserve/internal/infra/primary/http2/middleware"
	"central_reserve/internal/pkg/jwt"
	"central_reserve/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registra las rutas del handler de reservas
func RegisterRoutes(v1Group *gin.RouterGroup, handler IReserveHandler, jwtService *jwt.JWTService, logger log.ILogger) {
	// Crear el subgrupo /reserves dentro de /api/v1
	reserves := v1Group.Group("/reserves")
	// Aplicar middleware de autenticación a todas las rutas de reservas
	reserves.Use(middleware.AuthMiddleware(jwtService, logger))
	{
		reserves.GET("", handler.GetReservesHandler)
		reserves.GET("/:id", handler.GetReserveByIDHandler)
		reserves.POST("", handler.CreateReserveHandler)
		reserves.PUT("/:id", handler.UpdateReservationHandler)
		reserves.PATCH("/:id/cancel", handler.CancelReservationHandler)
	}
}
