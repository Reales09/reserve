package businesshandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/mapper"
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/request"
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateBusinessHandler maneja la solicitud de actualizar un negocio
// @Summary Actualizar negocio
// @Description Actualiza un negocio existente con los datos proporcionados
// @Tags Business
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del negocio"
// @Param business body request.BusinessRequest true "Datos del negocio"
// @Success 200 {object} response.BusinessSuccessResponse "Negocio actualizado exitosamente"
// @Failure 400 {object} response.BusinessErrorResponse "Datos inválidos o ID inválido"
// @Failure 401 {object} response.BusinessErrorResponse "Token de acceso requerido"
// @Failure 404 {object} response.BusinessErrorResponse "Negocio no encontrado"
// @Failure 409 {object} response.BusinessErrorResponse "Código o dominio ya existe"
// @Failure 500 {object} response.BusinessErrorResponse "Error interno del servidor"
// @Router /businesses/{id} [put]
func (h *BusinessHandler) UpdateBusinessHandler(c *gin.Context) {
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

	// Obtener datos del request
	var req request.BusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar datos del request")
		c.JSON(http.StatusBadRequest, response.BusinessErrorResponse{
			Success: false,
			Error:   "invalid_data",
			Message: "Datos inválidos: " + err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Convertir request a DTO de dominio
	businessRequest := mapper.ToBusinessRequest(req)

	// Ejecutar caso de uso
	business, err := h.usecase.UpdateBusiness(ctx, uint(id), businessRequest)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", uint(id)).Str("name", req.Name).Msg("Error al actualizar negocio")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "negocio no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Negocio no encontrado"
		} else if err.Error() == "el código '"+req.Code+"' ya existe" {
			statusCode = http.StatusConflict
			errorMessage = "El código ya existe"
		} else if err.Error() == "el dominio '"+req.CustomDomain+"' ya existe" {
			statusCode = http.StatusConflict
			errorMessage = "El dominio ya existe"
		}

		c.JSON(statusCode, response.BusinessErrorResponse{
			Success: false,
			Error:   "update_failed",
			Message: errorMessage,
		})
		return
	}

	// Convertir respuesta de dominio a response
	businessResponse := mapper.ToBusinessResponse(*business)

	h.logger.Info().Uint("id", uint(id)).Str("name", req.Name).Msg("Negocio actualizado exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.BusinessSuccessResponse{
		Success: true,
		Message: "Negocio actualizado exitosamente",
		Data:    businessResponse,
	})
}
