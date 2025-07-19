package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"dbpostgres/db/models"

	"gorm.io/gorm"
)

// UserContext contiene la información del usuario autenticado
type UserContext struct {
	UserID        uint
	Email         string
	RestaurantIDs []uint // Lista de IDs de restaurantes a los que pertenece
	Roles         []string
	Permissions   []string
}

// AuthService maneja la autenticación y autorización
type AuthService struct {
	db *gorm.DB
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// AuthenticateUser autentica un usuario y retorna su contexto
func (as *AuthService) AuthenticateUser(email, password string) (*UserContext, error) {
	var user models.User
	err := as.db.Preload("Roles.Permissions").Preload("Restaurants").Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("usuario no encontrado o inactivo")
		}
		return nil, err
	}

	// TODO: Implementar verificación de contraseña con bcrypt
	// if !CheckPasswordHash(password, user.Password) {
	//     return nil, errors.New("contraseña incorrecta")
	// }

	// Construir contexto del usuario
	userContext := &UserContext{
		UserID:        user.ID,
		Email:         user.Email,
		RestaurantIDs: make([]uint, 0),
		Roles:         make([]string, 0),
		Permissions:   make([]string, 0),
	}

	// Extraer IDs de restaurantes
	for _, restaurant := range user.Restaurants {
		userContext.RestaurantIDs = append(userContext.RestaurantIDs, restaurant.ID)
	}

	// Extraer roles y permisos
	permissionSet := make(map[string]bool)
	for _, role := range user.Roles {
		userContext.Roles = append(userContext.Roles, role.Code)
		for _, permission := range role.Permissions {
			permissionKey := fmt.Sprintf("%s:%s", permission.Resource, permission.Action)
			if !permissionSet[permissionKey] {
				permissionSet[permissionKey] = true
				userContext.Permissions = append(userContext.Permissions, permissionKey)
			}
		}
	}

	return userContext, nil
}

// HasPermission verifica si un usuario tiene un permiso específico
func (as *AuthService) HasPermission(userCtx *UserContext, resource, action string) bool {
	permissionKey := fmt.Sprintf("%s:%s", resource, action)
	for _, permission := range userCtx.Permissions {
		if permission == permissionKey {
			return true
		}
	}
	return false
}

// HasRole verifica si un usuario tiene un rol específico
func (as *AuthService) HasRole(userCtx *UserContext, roleCode string) bool {
	for _, role := range userCtx.Roles {
		if role == roleCode {
			return true
		}
	}
	return false
}

// HasAnyRole verifica si un usuario tiene al menos uno de los roles especificados
func (as *AuthService) HasAnyRole(userCtx *UserContext, roleCodes ...string) bool {
	for _, roleCode := range roleCodes {
		if as.HasRole(userCtx, roleCode) {
			return true
		}
	}
	return false
}

// IsSuperUser verifica si el usuario es un super usuario
func (as *AuthService) IsSuperUser(userCtx *UserContext) bool {
	return as.HasAnyRole(userCtx, models.ROLE_SUPER_ADMIN, models.ROLE_PLATFORM_ADMIN)
}

// IsRestaurantUser verifica si el usuario pertenece a algún restaurante
func (as *AuthService) IsRestaurantUser(userCtx *UserContext) bool {
	return len(userCtx.RestaurantIDs) > 0
}

// CanAccessRestaurant verifica si el usuario puede acceder a un restaurante específico
func (as *AuthService) CanAccessRestaurant(userCtx *UserContext, restaurantID uint) bool {
	// Super usuarios pueden acceder a cualquier restaurante
	if as.IsSuperUser(userCtx) {
		return true
	}

	// Usuarios de restaurante solo pueden acceder a sus restaurantes
	if as.IsRestaurantUser(userCtx) {
		for _, id := range userCtx.RestaurantIDs {
			if id == restaurantID {
				return true
			}
		}
	}

	return false
}

// GetUserRestaurantIDs retorna los IDs de los restaurantes del usuario
func (as *AuthService) GetUserRestaurantIDs(userCtx *UserContext) []uint {
	return userCtx.RestaurantIDs
}

// ContextKey es la clave para almacenar el contexto del usuario
type ContextKey string

const UserContextKey ContextKey = "user_context"

// GetUserFromContext obtiene el contexto del usuario desde el contexto HTTP
func GetUserFromContext(ctx context.Context) (*UserContext, bool) {
	userCtx, ok := ctx.Value(UserContextKey).(*UserContext)
	return userCtx, ok
}

// SetUserInContext establece el contexto del usuario en el contexto HTTP
func SetUserInContext(ctx context.Context, userCtx *UserContext) context.Context {
	return context.WithValue(ctx, UserContextKey, userCtx)
}

// RequirePermission middleware que requiere un permiso específico
func (as *AuthService) RequirePermission(resource, action string) func(*UserContext) error {
	return func(userCtx *UserContext) error {
		if !as.HasPermission(userCtx, resource, action) {
			return fmt.Errorf("permiso requerido: %s:%s", resource, action)
		}
		return nil
	}
}

// RequireRole middleware que requiere un rol específico
func (as *AuthService) RequireRole(roleCode string) func(*UserContext) error {
	return func(userCtx *UserContext) error {
		if !as.HasRole(userCtx, roleCode) {
			return fmt.Errorf("rol requerido: %s", roleCode)
		}
		return nil
	}
}

// RequireAnyRole middleware que requiere al menos uno de los roles especificados
func (as *AuthService) RequireAnyRole(roleCodes ...string) func(*UserContext) error {
	return func(userCtx *UserContext) error {
		if !as.HasAnyRole(userCtx, roleCodes...) {
			return fmt.Errorf("se requiere al menos uno de los roles: %s", strings.Join(roleCodes, ", "))
		}
		return nil
	}
}

// RequireSuperUser middleware que requiere ser super usuario
func (as *AuthService) RequireSuperUser() func(*UserContext) error {
	return func(userCtx *UserContext) error {
		if !as.IsSuperUser(userCtx) {
			return errors.New("se requiere ser super usuario")
		}
		return nil
	}
}

// RequireRestaurantAccess middleware que requiere acceso a un restaurante específico
func (as *AuthService) RequireRestaurantAccess(restaurantID uint) func(*UserContext) error {
	return func(userCtx *UserContext) error {
		if !as.CanAccessRestaurant(userCtx, restaurantID) {
			return fmt.Errorf("sin acceso al restaurante %d", restaurantID)
		}
		return nil
	}
}
