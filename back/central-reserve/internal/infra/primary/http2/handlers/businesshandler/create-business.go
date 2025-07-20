package businesshandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/mapper"
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/request"
	"central_reserve/internal/infra/primary/http2/handlers/businesshandler/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateBusinessHandler maneja la solicitud de crear un nuevo negocio
// @Summary Crear negocio
// @Description Crea un nuevo negocio con los datos proporcionados
// @Tags Business
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param business body request.BusinessRequest true "Datos del negocio"
// @Success 201 {object} response.BusinessSuccessResponse "Negocio creado exitosamente"
// @Failure 400 {object} response.BusinessErrorResponse "Datos inválidos"
// @Failure 401 {object} response.BusinessErrorResponse "Token de acceso requerido"
// @Failure 409 {object} response.BusinessErrorResponse "Código o dominio ya existe"
// @Failure 500 {object} response.BusinessErrorResponse "Error interno del servidor"
// @Router /businesses [post]
func (h *BusinessHandler) CreateBusinessHandler(c *gin.Context) {
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
	business, err := h.usecase.CreateBusiness(ctx, businessRequest)
	if err != nil {
		h.logger.Error().Err(err).Str("name", req.Name).Str("code", req.Code).Msg("Error al crear negocio")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "el código '"+req.Code+"' ya existe" {
			statusCode = http.StatusConflict
			errorMessage = "El código ya existe"
		} else if err.Error() == "el dominio '"+req.CustomDomain+"' ya existe" {
			statusCode = http.StatusConflict
			errorMessage = "El dominio ya existe"
		}

		c.JSON(statusCode, response.BusinessErrorResponse{
			Success: false,
			Error:   "creation_failed",
			Message: errorMessage,
		})
		return
	}

	// Convertir respuesta de dominio a response
	businessResponse := mapper.ToBusinessResponse(*business)

	h.logger.Info().Uint("id", business.ID).Str("name", req.Name).Msg("Negocio creado exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusCreated, response.BusinessSuccessResponse{
		Success: true,
		Message: "Negocio creado exitosamente",
		Data:    businessResponse,
	})
}
