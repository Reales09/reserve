package businesshandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteBusinessHandler maneja la solicitud de eliminar un negocio
// @Summary Eliminar negocio
// @Description Elimina un negocio existente por su ID
// @Tags Business
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del negocio"
// @Success 200 {object} response.BusinessErrorResponse "Negocio eliminado exitosamente"
// @Failure 400 {object} response.BusinessErrorResponse "ID inválido"
// @Failure 401 {object} response.BusinessErrorResponse "Token de acceso requerido"
// @Failure 404 {object} response.BusinessErrorResponse "Negocio no encontrado"
// @Failure 500 {object} response.BusinessErrorResponse "Error interno del servidor"
// @Router /businesses/{id} [delete]
func (h *BusinessHandler) DeleteBusinessHandler(c *gin.Context) {
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
	err = h.usecase.DeleteBusiness(ctx, uint(id))
	if err != nil {
		h.logger.Error().Err(err).Uint("id", uint(id)).Msg("Error al eliminar negocio")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "negocio no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Negocio no encontrado"
		}

		c.JSON(statusCode, response.BusinessErrorResponse{
			Success: false,
			Error:   "deletion_failed",
			Message: errorMessage,
		})
		return
	}

	h.logger.Info().Uint("id", uint(id)).Msg("Negocio eliminado exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.BusinessErrorResponse{
		Success: true,
		Error:   "",
		Message: "Negocio eliminado exitosamente",
	})
}
