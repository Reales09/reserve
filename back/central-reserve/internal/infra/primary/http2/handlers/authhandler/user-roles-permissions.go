package authhandler

import (
	"central_reserve/internal/infra/primary/http2/handlers/authhandler/mapper"
	"central_reserve/internal/infra/primary/http2/handlers/authhandler/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserRolesPermissionsHandler maneja la solicitud de obtener roles y permisos del usuario
// @Summary Obtener roles y permisos del usuario
// @Description Obtiene los roles y permisos de un usuario específico
// @Tags Auth
// @Accept json
// @Produce json
// @Param user_id path int true "ID del usuario"
// @Security BearerAuth
// @Success 200 {object} response.UserRolesPermissionsSuccessResponse "Roles y permisos obtenidos exitosamente"
// @Failure 400 {object} response.LoginErrorResponse "ID de usuario inválido"
// @Failure 401 {object} response.LoginErrorResponse "Token de acceso requerido"
// @Failure 403 {object} response.LoginErrorResponse "Acceso denegado"
// @Failure 404 {object} response.LoginErrorResponse "Usuario no encontrado"
// @Failure 500 {object} response.LoginErrorResponse "Error interno del servidor"
// @Router /auth/users/{user_id}/roles-permissions [get]
func (h *AuthHandler) GetUserRolesPermissionsHandler(c *gin.Context) {
	// Obtener el ID del usuario desde la URL
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userIDStr).Msg("ID de usuario inválido")
		c.JSON(http.StatusBadRequest, response.LoginErrorResponse{
			Error: "ID de usuario inválido",
		})
		return
	}

	// Obtener el token del header de autorización
	token := c.GetHeader("Authorization")
	if token == "" {
		h.logger.Error().Msg("Token de autorización requerido")
		c.JSON(http.StatusUnauthorized, response.LoginErrorResponse{
			Error: "Token de acceso requerido",
		})
		return
	}

	// Remover el prefijo "Bearer " si existe
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Ejecutar caso de uso para obtener roles y permisos
	rolesPermissions, err := h.usecase.GetUserRolesPermissions(c.Request.Context(), uint(userID), token)
	if err != nil {
		h.logger.Error().Err(err).Uint("user_id", uint(userID)).Msg("Error al obtener roles y permisos del usuario")

		// Determinar el código de estado HTTP apropiado
		statusCode := http.StatusInternalServerError
		errorMessage := "Error interno del servidor"

		if err.Error() == "usuario no encontrado" {
			statusCode = http.StatusNotFound
			errorMessage = "Usuario no encontrado"
		} else if err.Error() == "token inválido" {
			statusCode = http.StatusUnauthorized
			errorMessage = "Token inválido"
		} else if err.Error() == "acceso denegado" {
			statusCode = http.StatusForbidden
			errorMessage = "Acceso denegado"
		}

		c.JSON(statusCode, response.LoginErrorResponse{
			Error: errorMessage,
		})
		return
	}

	// Convertir respuesta de dominio a response
	rolesPermissionsResponse := mapper.ToUserRolesPermissionsResponse(rolesPermissions)

	h.logger.Info().
		Uint("user_id", uint(userID)).
		Bool("is_super", rolesPermissions.IsSuper).
		Int("roles_count", len(rolesPermissions.Roles)).
		Int("permissions_count", len(rolesPermissions.Permissions)).
		Msg("Roles y permisos obtenidos exitosamente")

	// Retornar respuesta exitosa
	c.JSON(http.StatusOK, response.UserRolesPermissionsSuccessResponse{
		Success: true,
		Data:    rolesPermissionsResponse,
	})
}
