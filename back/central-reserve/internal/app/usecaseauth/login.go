package usecaseauth

import (
	"central_reserve/internal/domain/dtos"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Login maneja la lógica de autenticación del usuario (simplificado)
func (uc *AuthUseCase) Login(ctx context.Context, request dtos.LoginRequest) (*dtos.LoginResponse, error) {
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

	// Obtener roles del usuario para el token
	roles, err := uc.repository.GetUserRoles(ctx, user.ID)
	if err != nil {
		uc.log.Error().Err(err).Uint("user_id", user.ID).Msg("Error al obtener roles del usuario")
		return nil, fmt.Errorf("error interno del servidor")
	}

	// Generar token JWT
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

	// Construir respuesta simplificada
	response := &dtos.LoginResponse{
		User: dtos.UserInfo{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Phone:       user.Phone,
			AvatarURL:   user.AvatarURL,
			IsActive:    user.IsActive,
			LastLoginAt: user.LastLoginAt,
		},
		Token: token,
	}

	uc.log.Info().
		Str("email", user.Email).
		Uint("user_id", user.ID).
		Msg("Login exitoso")

	return response, nil
}
