package migrate

import (
	"dbpostgres/db/models"
	"dbpostgres/db/repository"
	"dbpostgres/pkg/env"
	"dbpostgres/pkg/log"

	"golang.org/x/crypto/bcrypt"
)

// InitializeSuperAdminUser crea el usuario administrador súper por defecto
func InitializeSuperAdminUser(systemRepo repository.SystemRepository, logger log.ILogger, config env.IConfig) error {
	logger.Info().Msg("Inicializando usuario administrador súper...")

	// Obtener email del usuario por defecto desde variables de entorno
	userEmail := config.Get("EMAIL_USER_DEFAULT")
	if userEmail == "" {
		userEmail = "admin@reserve.com" // Email por defecto si no está configurado
		logger.Warn().Msg("EMAIL_USER_DEFAULT no configurado, usando email por defecto")
	}

	// Obtener contraseña del usuario por defecto desde variables de entorno
	userPass := config.Get("USER_PASS_DEFAULT")
	if userPass == "" {
		userPass = "admin123" // Contraseña por defecto si no está configurada
		logger.Warn().Msg("USER_PASS_DEFAULT no configurada, usando contraseña por defecto")
	}

	logger.Info().Str("email", userEmail).Msg("Configuración de usuario super admin cargada")

	// Verificar si ya existe ALGÚN usuario en la tabla (no solo este usuario específico)
	logger.Debug().Msg("Verificando si existe algún usuario en la tabla...")
	userExists, err := systemRepo.UserExists()
	if err != nil {
		logger.Error().Err(err).Msg("Error al verificar si existen usuarios en la tabla")
		return err
	}

	logger.Debug().Bool("user_exists", userExists).Msg("Resultado de verificación de usuarios existentes")

	if userExists {
		logger.Info().Msg("✅ Ya existen usuarios en el sistema, saltando creación de usuario super admin")
		return nil
	}

	logger.Info().Str("email", userEmail).Msg("No existen usuarios en el sistema, creando primer usuario super admin...")

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPass), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(err).Msg("Error al hashear contraseña del super admin")
		return err
	}

	logger.Debug().Msg("Contraseña hasheada correctamente")

	// Crear usuario super admin
	superAdmin := models.User{
		Name:     "Super Administrador",
		Email:    userEmail,
		Password: string(hashedPassword),
		Phone:    "",
		IsActive: true,
	}

	logger.Debug().Str("email", userEmail).Msg("Intentando crear primer usuario en la base de datos...")

	// Crear el usuario
	if err := systemRepo.Create(&superAdmin); err != nil {
		logger.Error().Err(err).Str("email", userEmail).Msg("Error al crear usuario super admin")
		return err
	}

	logger.Debug().Uint("user_id", superAdmin.ID).Str("email", userEmail).Msg("Primer usuario creado exitosamente en la base de datos")

	// Asignar rol de super admin
	logger.Debug().Str("email", userEmail).Msg("Asignando rol super_admin...")
	if err := systemRepo.AssignRolesToUser(superAdmin.Email, []string{"super_admin"}); err != nil {
		logger.Error().Err(err).Str("email", userEmail).Msg("Error al asignar rol super_admin al usuario")
		return err
	}

	logger.Info().Str("email", superAdmin.Email).Msg("✅ Primer usuario super admin creado exitosamente")
	return nil
}
