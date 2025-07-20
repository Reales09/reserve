package businesshandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/mapper"
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetBusinessesHandler maneja la solicitud de obtener todos los negocios
// @Summary Obtener todos los negocios
// @Description Obtiene la lista completa de negocios disponibles
// @Tags Business
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BusinessListSuccessResponse "Negocios obtenidos exitosamente"
// @Failure 401 {object} response.BusinessErrorResponse "Token de acceso requerido"
// @Failure 500 {object} response.BusinessErrorResponse "Error interno del servidor"
// @Router /businesses [get]
func (h *BusinessHandler) GetBusinessesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// Ejecutar caso de uso
	businesses, err := h.usecase.GetBusinesses(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al obtener negocios")
		c.JSON(http.StatusInternalServerError, response.BusinessErrorResponse{
			Success: false,
			Error:   "internal_error",
			Message: "Error interno del servidor",
		})
		return
	}

	// Convertir respuesta de dominio a response
	businessesResponse := mapper.ToBusinessListResponse(businesses)

	h.logger.Info().Int("count", len(businesses)).Msg("Negocios obtenidos exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.BusinessListSuccessResponse{
		Success: true,
		Message: "Negocios obtenidos exitosamente",
		Data:    businessesResponse,
	})
}
