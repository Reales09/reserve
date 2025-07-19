package migrate

import (
	"dbpostgres/db/models"
	"dbpostgres/db/repository"
	"dbpostgres/pkg/log"
)

// InitializeRoles inicializa los roles del sistema
func InitializeRoles(systemRepo repository.SystemRepository, logger log.ILogger) error {
	// Verificar si todos los roles ya existen
	roles := []models.Role{
		// Roles de plataforma
		{Name: "Super Administrador", Code: "super_admin", Description: "Acceso completo a toda la plataforma", Level: 1, IsSystem: true, Scope: models.SCOPE_PLATFORM},
		{Name: "Administrador de Plataforma", Code: "platform_admin", Description: "Administrador de la plataforma", Level: 2, IsSystem: true, Scope: models.SCOPE_PLATFORM},
		{Name: "Soporte Técnico", Code: "support", Description: "Soporte técnico de la plataforma", Level: 3, IsSystem: true, Scope: models.SCOPE_PLATFORM},

		// Roles de restaurante
		{Name: "Administrador de Restaurante", Code: "restaurant_admin", Description: "Administrador del restaurante", Level: 2, IsSystem: true, Scope: models.SCOPE_RESTAURANT},
		{Name: "Gerente", Code: "manager", Description: "Gerente del restaurante", Level: 3, IsSystem: true, Scope: models.SCOPE_RESTAURANT},
		{Name: "Hostess", Code: "hostess", Description: "Hostess del restaurante", Level: 4, IsSystem: true, Scope: models.SCOPE_RESTAURANT},
		{Name: "Mesero", Code: "waiter", Description: "Mesero del restaurante", Level: 4, IsSystem: true, Scope: models.SCOPE_RESTAURANT},
		{Name: "Cocinero", Code: "cook", Description: "Cocinero del restaurante", Level: 4, IsSystem: true, Scope: models.SCOPE_RESTAURANT},
	}

	// Verificar si todos los roles ya existen
	allExist := true
	for _, role := range roles {
		exists, err := systemRepo.ExistsByCode(role.Code)
		if err != nil {
			return err
		}
		if !exists {
			allExist = false
			break
		}
	}

	if allExist {
		logger.Info().Int("roles_count", len(roles)).Msg("✅ Roles del sistema ya existen, saltando migración de roles")
		return nil
	}

	// Inicializar roles usando el repositorio
	if err := systemRepo.InitializeRoles(roles); err != nil {
		return err
	}

	logger.Info().Int("roles_count", len(roles)).Msg("✅ Roles del sistema inicializados correctamente")
	return nil
}
