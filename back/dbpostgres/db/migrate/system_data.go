package migrate

import (
	"dbpostgres/db/repository"
	"dbpostgres/pkg/env"
	"dbpostgres/pkg/log"
)

// InitializeSystemData inicializa roles, permisos y usuario administrador del sistema
func InitializeSystemData(systemRepo repository.SystemRepository, logger log.ILogger, config env.IConfig) error {
	logger.Info().Msg("Inicializando datos del sistema...")

	// Inicializar permisos
	if err := InitializePermissions(systemRepo, logger); err != nil {
		return err
	}

	// Inicializar roles
	if err := InitializeRoles(systemRepo, logger); err != nil {
		return err
	}

	// Asignar permisos a roles
	if err := AssignPermissionsToRoles(systemRepo, logger); err != nil {
		return err
	}

	// Inicializar estados de reserva
	if err := InitializeReservationStatuses(systemRepo, logger); err != nil {
		return err
	}

	// Crear usuario administrador súper por defecto
	if err := InitializeSuperAdminUser(systemRepo, logger, config); err != nil {
		return err
	}

	// Crear restaurante de pruebas por defecto
	if err := InitializeTestRestaurant(systemRepo, logger); err != nil {
		return err
	}

	logger.Info().Msg("Datos del sistema inicializados correctamente")
	return nil
}
