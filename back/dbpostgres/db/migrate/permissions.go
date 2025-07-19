package migrate

import (
	"dbpostgres/db/models"
	"dbpostgres/db/repository"
	"dbpostgres/pkg/log"
)

// InitializePermissions inicializa los permisos del sistema
func InitializePermissions(systemRepo repository.SystemRepository, logger log.ILogger) error {
	// Verificar si todos los permisos ya existen
	permissions := []models.Permission{
		// Permisos de plataforma
		{Name: "Gestionar Restaurantes", Code: "restaurants:manage", Description: "Crear, editar y eliminar restaurantes", Resource: models.RESOURCE_RESTAURANTS, Action: models.ACTION_MANAGE, Scope: models.SCOPE_PLATFORM},
		{Name: "Ver Restaurantes", Code: "restaurants:read", Description: "Ver información de restaurantes", Resource: models.RESOURCE_RESTAURANTS, Action: models.ACTION_READ, Scope: models.SCOPE_PLATFORM},
		{Name: "Gestionar Usuarios", Code: "users:manage", Description: "Crear, editar y eliminar usuarios", Resource: models.RESOURCE_USERS, Action: models.ACTION_MANAGE, Scope: models.SCOPE_PLATFORM},
		{Name: "Ver Usuarios", Code: "users:read", Description: "Ver información de usuarios", Resource: models.RESOURCE_USERS, Action: models.ACTION_READ, Scope: models.SCOPE_PLATFORM},
		{Name: "Gestionar Roles", Code: "roles:manage", Description: "Crear, editar y eliminar roles", Resource: models.RESOURCE_ROLES, Action: models.ACTION_MANAGE, Scope: models.SCOPE_PLATFORM},
		{Name: "Ver Roles", Code: "roles:read", Description: "Ver información de roles", Resource: models.RESOURCE_ROLES, Action: models.ACTION_READ, Scope: models.SCOPE_PLATFORM},
		{Name: "Gestionar Permisos", Code: "permissions:manage", Description: "Crear, editar y eliminar permisos", Resource: models.RESOURCE_PERMISSIONS, Action: models.ACTION_MANAGE, Scope: models.SCOPE_PLATFORM},
		{Name: "Ver Reportes", Code: "reports:read", Description: "Ver reportes de la plataforma", Resource: models.RESOURCE_REPORTS, Action: models.ACTION_READ, Scope: models.SCOPE_PLATFORM},

		// Permisos de restaurante
		{Name: "Gestionar Reservas", Code: "reservations:manage", Description: "Crear, editar y eliminar reservas", Resource: models.RESOURCE_RESERVATIONS, Action: models.ACTION_MANAGE, Scope: models.SCOPE_RESTAURANT},
		{Name: "Ver Reservas", Code: "reservations:read", Description: "Ver información de reservas", Resource: models.RESOURCE_RESERVATIONS, Action: models.ACTION_READ, Scope: models.SCOPE_RESTAURANT},
		{Name: "Gestionar Mesas", Code: "tables:manage", Description: "Crear, editar y eliminar mesas", Resource: models.RESOURCE_TABLES, Action: models.ACTION_MANAGE, Scope: models.SCOPE_RESTAURANT},
		{Name: "Ver Mesas", Code: "tables:read", Description: "Ver información de mesas", Resource: models.RESOURCE_TABLES, Action: models.ACTION_READ, Scope: models.SCOPE_RESTAURANT},
		{Name: "Gestionar Clientes", Code: "clients:manage", Description: "Crear, editar y eliminar clientes", Resource: models.RESOURCE_CLIENTS, Action: models.ACTION_MANAGE, Scope: models.SCOPE_RESTAURANT},
		{Name: "Ver Clientes", Code: "clients:read", Description: "Ver información de clientes", Resource: models.RESOURCE_CLIENTS, Action: models.ACTION_READ, Scope: models.SCOPE_RESTAURANT},
		{Name: "Gestionar Personal", Code: "staff:manage", Description: "Gestionar personal del restaurante", Resource: models.RESOURCE_USERS, Action: models.ACTION_MANAGE, Scope: models.SCOPE_RESTAURANT},
		{Name: "Ver Personal", Code: "staff:read", Description: "Ver personal del restaurante", Resource: models.RESOURCE_USERS, Action: models.ACTION_READ, Scope: models.SCOPE_RESTAURANT},
		{Name: "Gestionar Configuración", Code: "restaurant:configure", Description: "Configurar restaurante (marca blanca)", Resource: models.RESOURCE_RESTAURANTS, Action: models.ACTION_UPDATE, Scope: models.SCOPE_RESTAURANT},
		{Name: "Ver Reportes del Restaurante", Code: "restaurant:reports", Description: "Ver reportes del restaurante", Resource: models.RESOURCE_REPORTS, Action: models.ACTION_READ, Scope: models.SCOPE_RESTAURANT},

		// Permisos de delivery (futuro)
		{Name: "Gestionar Delivery", Code: "delivery:manage", Description: "Gestionar pedidos de delivery", Resource: models.RESOURCE_DELIVERY, Action: models.ACTION_MANAGE, Scope: models.SCOPE_RESTAURANT},
		{Name: "Ver Delivery", Code: "delivery:read", Description: "Ver pedidos de delivery", Resource: models.RESOURCE_DELIVERY, Action: models.ACTION_READ, Scope: models.SCOPE_RESTAURANT},
		{Name: "Gestionar Menú", Code: "menu:manage", Description: "Gestionar menú del restaurante", Resource: models.RESOURCE_MENU, Action: models.ACTION_MANAGE, Scope: models.SCOPE_RESTAURANT},
		{Name: "Ver Menú", Code: "menu:read", Description: "Ver menú del restaurante", Resource: models.RESOURCE_MENU, Action: models.ACTION_READ, Scope: models.SCOPE_RESTAURANT},
	}

	// Verificar si todos los permisos ya existen
	allExist := true
	for _, permission := range permissions {
		exists, err := systemRepo.ExistsByCode(permission.Code)
		if err != nil {
			return err
		}
		if !exists {
			allExist = false
			break
		}
	}

	if allExist {
		logger.Info().Int("permissions_count", len(permissions)).Msg("✅ Permisos del sistema ya existen, saltando migración de permisos")
		return nil
	}

	// Inicializar permisos usando el repositorio
	if err := systemRepo.InitializePermissions(permissions); err != nil {
		return err
	}

	logger.Info().Int("permissions_count", len(permissions)).Msg("✅ Permisos del sistema inicializados correctamente")
	return nil
}
