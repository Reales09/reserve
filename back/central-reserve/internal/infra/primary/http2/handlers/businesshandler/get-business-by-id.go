package businesshandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/mapper"
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetBusinessByIDHandler maneja la solicitud de obtener un negocio por ID
// @Summary Obtener negocio por ID
// @Description Obtiene un negocio específico por su ID
// @Tags Business
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del negocio"
// @Success 200 {object} response.BusinessSuccessResponse "Negocio obtenido exitosamente"
// @Failure 400 {object} response.BusinessErrorResponse "ID inválido"
// @Failure 401 {object} response.BusinessErrorResponse "Token de acceso requerido"
// @Failure 404 {object} response.BusinessErrorResponse "Negocio no encontrado"
// @Failure 500 {object} response.BusinessErrorResponse "Error interno del servidor"
// @Router /businesses/{id} [get]
func (h *BusinessHandler) GetBusinessByIDHandler(c *gin.Context) {
	// Obtener ID del path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("id", idStr).Msg("ID inválido")
		c.JSON(http.StatusBadRequest, response.BusinessErrorResponse{
			Success: false,
			Error:   "invalid_id",
			Message: "ID inválido",
		})
		return
	}

	ctx := c.Request.Context()

	// Ejecutar caso de uso
	business, err := h.usecase.GetBusinessByID(ctx, uint(id))
	if err != nil {
		h.logger.Error().Err(err).Uint("id", uint(id)).Msg("Error al obtener negocio por ID")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "negocio no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Negocio no encontrado"
		}

		c.JSON(statusCode, response.BusinessErrorResponse{
			Success: false,
			Error:   "not_found",
			Message: errorMessage,
		})
		return
	}

	// Convertir respuesta de dominio a response
	businessResponse := mapper.ToBusinessResponse(*business)

	h.logger.Info().Uint("id", uint(id)).Msg("Negocio obtenido exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.BusinessSuccessResponse{
		Success: true,
		Message: "Negocio obtenido exitosamente",
		Data:    businessResponse,
	})
}
