package businesstypehandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/mapper"
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/request"
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateBusinessTypeHandler maneja la solicitud de actualizar un tipo de negocio
// @Summary Actualizar tipo de negocio
// @Description Actualiza un tipo de negocio existente con los datos proporcionados
// @Tags BusinessType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del tipo de negocio"
// @Param businessType body request.BusinessTypeRequest true "Datos del tipo de negocio"
// @Success 200 {object} response.BusinessTypeSuccessResponse "Tipo de negocio actualizado exitosamente"
// @Failure 400 {object} response.BusinessTypeErrorResponse "Datos inválidos o ID inválido"
// @Failure 401 {object} response.BusinessTypeErrorResponse "Token de acceso requerido"
// @Failure 404 {object} response.BusinessTypeErrorResponse "Tipo de negocio no encontrado"
// @Failure 409 {object} response.BusinessTypeErrorResponse "Código ya existe"
// @Failure 500 {object} response.BusinessTypeErrorResponse "Error interno del servidor"
// @Router /business-types/{id} [put]
func (h *BusinessTypeHandler) UpdateBusinessTypeHandler(c *gin.Context) {
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

	// Obtener datos del request
	var req request.BusinessTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Error al validar datos del request")
		c.JSON(http.StatusBadRequest, response.BusinessTypeErrorResponse{
			Success: false,
			Error:   "invalid_data",
			Message: "Datos inválidos: " + err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// Convertir request a DTO de dominio
	businessTypeRequest := mapper.ToBusinessTypeRequest(req)

	// Ejecutar caso de uso
	businessType, err := h.usecase.UpdateBusinessType(ctx, uint(id), businessTypeRequest)
	if err != nil {
		h.logger.Error().Err(err).Uint("id", uint(id)).Str("name", req.Name).Msg("Error al actualizar tipo de negocio")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "tipo de negocio no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Tipo de negocio no encontrado"
		} else if err.Error() == "el código '"+req.Code+"' ya existe" {
			statusCode = http.StatusConflict
			errorMessage = "El código ya existe"
		}

		c.JSON(statusCode, response.BusinessTypeErrorResponse{
			Success: false,
			Error:   "update_failed",
			Message: errorMessage,
		})
		return
	}

	// Convertir respuesta de dominio a response
	businessTypeResponse := mapper.ToBusinessTypeResponse(*businessType)

	h.logger.Info().Uint("id", uint(id)).Str("name", req.Name).Msg("Tipo de negocio actualizado exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.BusinessTypeSuccessResponse{
		Success: true,
		Message: "Tipo de negocio actualizado exitosamente",
		Data:    businessTypeResponse,
	})
}
