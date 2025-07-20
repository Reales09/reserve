package businesstypehandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/mapper"
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetBusinessTypesHandler maneja la solicitud de obtener todos los tipos de negocio
// @Summary Obtener todos los tipos de negocio
// @Description Obtiene la lista completa de tipos de negocio disponibles
// @Tags BusinessType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.BusinessTypeListSuccessResponse "Tipos de negocio obtenidos exitosamente"
// @Failure 401 {object} response.BusinessTypeErrorResponse "Token de acceso requerido"
// @Failure 500 {object} response.BusinessTypeErrorResponse "Error interno del servidor"
// @Router /business-types [get]
func (h *BusinessTypeHandler) GetBusinessTypesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// Ejecutar caso de uso
	businessTypes, err := h.usecase.GetBusinessTypes(ctx)
	if err != nil {
		h.logger.Error().Err(err).Msg("Error al obtener tipos de negocio")
		c.JSON(http.StatusInternalServerError, response.BusinessTypeErrorResponse{
			Success: false,
			Error:   "internal_error",
			Message: "Error interno del servidor",
		})
		return
	}

	// Convertir respuesta de dominio a response
	businessTypesResponse := mapper.ToBusinessTypeListResponse(businessTypes)

	h.logger.Info().Int("count", len(businessTypes)).Msg("Tipos de negocio obtenidos exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.BusinessTypeListSuccessResponse{
		Success: true,
		Message: "Tipos de negocio obtenidos exitosamente",
		Data:    businessTypesResponse,
	})
}
