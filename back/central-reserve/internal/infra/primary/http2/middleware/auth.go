package middleware

import (
	"central_reserve/internal/pkg/jwt"
	"central_reserve/internal/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware crea un middleware de autenticación JWT
func AuthMiddleware(jwtService *jwt.JWTService, logger log.ILogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del header de autorización
		token := c.GetHeader("Authorization")
		if token == "" {
			logger.Error().Msg("Token de autorización requerido")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token de acceso requerido",
			})
			c.Abort()
			return
		}

		// Remover el prefijo "Bearer " si existe
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		// Validar y decodificar el token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			logger.Error().Err(err).Msg("Token inválido")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido",
			})
			c.Abort()
			return
		}

		// Almacenar los claims en el contexto para que los handlers puedan acceder
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Set("jwt_claims", claims)

		// Agregar información del usuario al logger para trazabilidad
		logger.Debug().
			Uint("user_id", claims.UserID).
			Str("user_email", claims.Email).
			Strs("user_roles", claims.Roles).
			Msg("Usuario autenticado")

		c.Next()
	}
}

// GetUserID obtiene el ID del usuario desde el contexto de Gin
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	if id, ok := userID.(uint); ok {
		return id, true
	}
	return 0, false
}

// GetUserEmail obtiene el email del usuario desde el contexto de Gin
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}
	if e, ok := email.(string); ok {
		return e, true
	}
	return "", false
}

// GetUserRoles obtiene los roles del usuario desde el contexto de Gin
func GetUserRoles(c *gin.Context) ([]string, bool) {
	roles, exists := c.Get("user_roles")
	if !exists {
		return nil, false
	}
	if r, ok := roles.([]string); ok {
		return r, true
	}
	return nil, false
}

// GetJWTClaims obtiene los claims completos del JWT desde el contexto de Gin
func GetJWTClaims(c *gin.Context) (*jwt.Claims, bool) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		return nil, false
	}
	if c, ok := claims.(*jwt.Claims); ok {
		return c, true
	}
	return nil, false
}

// RequireRole crea un middleware que requiere un rol específico
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := GetUserRoles(c)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: roles no disponibles",
			})
			c.Abort()
			return
		}

		// Verificar si el usuario tiene el rol requerido
		hasRole := false
		for _, role := range roles {
			if role == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: rol requerido",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole crea un middleware que requiere al menos uno de los roles especificados
func RequireAnyRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := GetUserRoles(c)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: roles no disponibles",
			})
			c.Abort()
			return
		}

		// Verificar si el usuario tiene al menos uno de los roles requeridos
		hasRole := false
		for _, userRole := range roles {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: rol requerido",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
