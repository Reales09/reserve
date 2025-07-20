package businesstypehandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/mapper"
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBusinessTypeByIDHandler maneja la solicitud de obtener un tipo de negocio por ID
// @Summary Obtener tipo de negocio por ID
// @Description Obtiene un tipo de negocio específico por su ID
// @Tags BusinessType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del tipo de negocio"
// @Success 200 {object} response.BusinessTypeSuccessResponse "Tipo de negocio obtenido exitosamente"
// @Failure 400 {object} response.BusinessTypeErrorResponse "ID inválido"
// @Failure 401 {object} response.BusinessTypeErrorResponse "Token de acceso requerido"
// @Failure 404 {object} response.BusinessTypeErrorResponse "Tipo de negocio no encontrado"
// @Failure 500 {object} response.BusinessTypeErrorResponse "Error interno del servidor"
// @Router /business-types/{id} [get]
func (h *BusinessTypeHandler) GetBusinessTypeByIDHandler(c *gin.Context) {
	// Obtener ID del path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, response.BusinessTypeErrorResponse{
			Success: false,
			Error:   "invalid_id",
			Message: "ID inválido",
		})
		return
	}

	ctx := c.Request.Context()

	// Ejecutar caso de uso
	businessType, err := h.usecase.GetBusinessTypeByID(ctx, uint(id))
	if err != nil {
		h.logger.Error().Err(err).Uint("id", uint(id)).Msg("Error al obtener tipo de negocio por ID")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "tipo de negocio no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Tipo de negocio no encontrado"
		}

		c.JSON(statusCode, response.BusinessTypeErrorResponse{
			Success: false,
			Error:   "not_found",
			Message: errorMessage,
		})
		return
	}

	// Convertir respuesta de dominio a response
	businessTypeResponse := mapper.ToBusinessTypeResponse(*businessType)

	h.logger.Info().Uint("id", uint(id)).Msg("Tipo de negocio obtenido exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.BusinessTypeSuccessResponse{
		Success: true,
		Message: "Tipo de negocio obtenido exitosamente",
		Data:    businessTypeResponse,
	})
}
