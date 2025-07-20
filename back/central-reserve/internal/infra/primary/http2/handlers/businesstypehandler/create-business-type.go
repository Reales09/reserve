package businesstypehandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/mapper"
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/request"
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateBusinessTypeHandler maneja la solicitud de crear un nuevo tipo de negocio
// @Summary Crear tipo de negocio
// @Description Crea un nuevo tipo de negocio con los datos proporcionados
// @Tags BusinessType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param businessType body request.BusinessTypeRequest true "Datos del tipo de negocio"
// @Success 201 {object} response.BusinessTypeSuccessResponse "Tipo de negocio creado exitosamente"
// @Failure 400 {object} response.BusinessTypeErrorResponse "Datos inválidos"
// @Failure 401 {object} response.BusinessTypeErrorResponse "Token de acceso requerido"
// @Failure 409 {object} response.BusinessTypeErrorResponse "Código ya existe"
// @Failure 500 {object} response.BusinessTypeErrorResponse "Error interno del servidor"
// @Router /business-types [post]
func (h *BusinessTypeHandler) CreateBusinessTypeHandler(c *gin.Context) {
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
	businessType, err := h.usecase.CreateBusinessType(ctx, businessTypeRequest)
	if err != nil {
		h.logger.Error().Err(err).Str("name", req.Name).Str("code", req.Code).Msg("Error al crear tipo de negocio")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "el código '"+req.Code+"' ya existe" {
			statusCode = http.StatusConflict
			errorMessage = "El código ya existe"
		}

		c.JSON(statusCode, response.BusinessTypeErrorResponse{
			Success: false,
			Error:   "creation_failed",
			Message: errorMessage,
		})
		return
	}

	// Convertir respuesta de dominio a response
	businessTypeResponse := mapper.ToBusinessTypeResponse(*businessType)

	h.logger.Info().Uint("id", businessType.ID).Str("name", req.Name).Msg("Tipo de negocio creado exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusCreated, response.BusinessTypeSuccessResponse{
		Success: true,
		Message: "Tipo de negocio creado exitosamente",
		Data:    businessTypeResponse,
	})
}
