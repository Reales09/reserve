package usecaseauth

import (
	"central_reserve/internal/domain"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Login maneja la lógica de autenticación del usuario
func (uc *AuthUseCase) Login(ctx context.Context, request domain.LoginRequest) (*domain.LoginResponse, error) {
	uc.log.Info().Str("email", request.Email).Msg("Iniciando proceso de login")

	// Validar que el email y contraseña no estén vacíos
	if request.Email == "" || request.Password == "" {
		uc.log.Error().Msg("Email o contraseña vacíos")
		return nil, fmt.Errorf("email y contraseña son requeridos")
	}

	// Obtener usuario por email
	user, err := uc.repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		uc.log.Error().Err(err).Str("email", request.Email).Msg("Error al obtener usuario por email")
		return nil, fmt.Errorf("credenciales inválidas")
	}

	if user == nil {
		uc.log.Error().Str("email", request.Email).Msg("Usuario no encontrado")
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// Verificar que el usuario esté activo
	if !user.IsActive {
		uc.log.Error().Str("email", request.Email).Msg("Usuario inactivo")
		return nil, fmt.Errorf("usuario inactivo")
	}

	// Validar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		uc.log.Error().Err(err).Str("email", request.Email).Msg("Contraseña inválida")
		return nil, fmt.Errorf("credenciales inválidas")
	}

	// Obtener roles del usuario
	roles, err := uc.repository.GetUserRoles(ctx, user.ID)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.ID).Msg("Error al obtener roles del usuario")
		return nil, fmt.Errorf("error interno del servidor")
	}

	// Obtener permisos de todos los roles del usuario
	var allPermissions []domain.Permission
	for _, role := range roles {
		permissions, err := uc.repository.GetRolePermissions(ctx, role.ID)
		if err != nil {
			uc.log.Error().Err(err).Uint("role_id", role.ID).Msg("Error al obtener permisos del rol")
			continue
		}
		allPermissions = append(allPermissions, permissions...)
	}

	// Verificar si es super admin (tiene rol super_admin)
	isSuper := false
	for _, role := range roles {
		if role.Code == "super_admin" {
			isSuper = true
			break
		}
	}

	// Generar token JWT real
	roleCodes := make([]string, len(roles))
	for i, role := range roles {
		roleCodes[i] = role.Code
	}

	token, err := uc.jwtService.GenerateToken(user.ID, user.Email, roleCodes)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.ID).Msg("Error al generar token JWT")
		return nil, fmt.Errorf("error interno del servidor")
	}

	// Actualizar último login
	if err := uc.repository.UpdateLastLogin(ctx, user.ID); err != nil {
		uc.log.Warn().Err(err).Uint("user_id", user.ID).Msg("Error al actualizar último login")
		// No retornamos error aquí porque el login ya fue exitoso
	}

	// Construir respuesta
	response := &domain.LoginResponse{
		User: domain.UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Phone:       user.Phone,
			AvatarURL:   user.AvatarURL,
			IsActive:    user.IsActive,
			LastLoginAt: user.LastLoginAt,
		},
		Token:       token,
		IsSuper:     isSuper,
		Roles:       make([]domain.RoleInfo, len(roles)),
		Permissions: make([]domain.PermissionInfo, len(allPermissions)),
	}

	// Mapear roles
	for i, role := range roles {
		response.Roles[i] = domain.RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Level:       role.Level,
			Scope:       role.Scope,
		}
	}

	// Mapear permisos (eliminar duplicados)
	permissionMap := make(map[string]domain.PermissionInfo)
	for _, permission := range allPermissions {
		key := permission.Code
		if _, exists := permissionMap[key]; !exists {
			permissionMap[key] = domain.PermissionInfo{
				ID:          permission.ID,
				Name:        permission.Name,
				Code:        permission.Code,
				Description: permission.Description,
				Resource:    permission.Resource,
				Action:      permission.Action,
				Scope:       permission.Scope,
			}
		}
	}

	// Convertir map a slice
	permissionIndex := 0
	for _, permission := range permissionMap {
		response.Permissions[permissionIndex] = permission
		permissionIndex++
	}
	response.Permissions = response.Permissions[:permissionIndex]

	uc.log.Info().
		Str("email", user.Email).
		Uint("user_id", user.ID).
		Int("roles_count", len(roles)).
		Int("permissions_count", len(response.Permissions)).
		Bool("is_super", isSuper).
		Msg("Login exitoso")

	return response, nil
}
