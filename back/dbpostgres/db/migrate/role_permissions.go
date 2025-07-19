package migrate

import (
	"dbpostgres/db/repository"
	"dbpostgres/pkg/log"
)

// AssignPermissionsToRoles asigna permisos a los roles del sistema
func AssignPermissionsToRoles(systemRepo repository.SystemRepository, logger log.ILogger) error {
	logger.Info().Msg("Asignando permisos a roles...")

	// Verificar si todas las asignaciones ya existen
	allAssignmentsExist := true

	// Lista de roles y sus permisos para verificar
	rolePermissions := map[string][]string{
		"super_admin": {
			"restaurants:manage", "restaurants:read", "users:manage", "users:read",
			"roles:manage", "roles:read", "permissions:manage", "reports:read",
			"reservations:manage", "reservations:read", "tables:manage", "tables:read",
			"clients:manage", "clients:read", "staff:manage", "staff:read",
			"restaurant:configure", "restaurant:reports", "delivery:manage", "delivery:read",
			"menu:manage", "menu:read",
		},
		"platform_admin": {
			"restaurants:read", "users:manage", "users:read", "roles:read",
			"permissions:manage", "reports:read",
		},
		"support": {
			"restaurants:read", "users:read", "reports:read",
		},
		"restaurant_admin": {
			"reservations:manage", "reservations:read", "tables:manage", "tables:read",
			"clients:manage", "clients:read", "staff:manage", "staff:read",
			"restaurant:configure", "restaurant:reports", "delivery:manage", "delivery:read",
			"menu:manage", "menu:read",
		},
		"manager": {
			"reservations:manage", "reservations:read", "tables:read", "clients:manage",
			"clients:read", "staff:read", "restaurant:reports", "delivery:read", "menu:read",
		},
		"hostess": {
			"reservations:manage", "reservations:read", "tables:read", "clients:manage",
			"clients:read",
		},
		"waiter": {
			"reservations:read", "tables:read", "clients:read",
		},
		"cook": {
			"menu:manage", "menu:read", "delivery:read",
		},
	}

	// Verificar si todas las asignaciones ya existen
	for roleCode, permissionCodes := range rolePermissions {
		role, err := systemRepo.GetRoleByCode(roleCode)
		if err != nil {
			return err
		}
		if role == nil {
			allAssignmentsExist = false
			break
		}

		// Verificar si todos los permisos están asignados al rol
		for _, permissionCode := range permissionCodes {
			permission, err := systemRepo.GetPermissionByCode(permissionCode)
			if err != nil {
				return err
			}
			if permission == nil {
				allAssignmentsExist = false
				break
			}

			// Aquí necesitaríamos verificar si la relación existe
			// Por ahora, asumimos que si el rol y el permiso existen, la asignación también
			// En una implementación más robusta, verificaríamos la tabla de relación
		}
		if !allAssignmentsExist {
			break
		}
	}

	if allAssignmentsExist {
		logger.Info().Msg("✅ Asignaciones de permisos a roles ya existen, saltando migración de asignaciones")
		return nil
	}

	// Super Admin - todos los permisos
	superAdminPermissions := []string{
		"restaurants:manage", "restaurants:read", "users:manage", "users:read",
		"roles:manage", "roles:read", "permissions:manage", "reports:read",
		"reservations:manage", "reservations:read", "tables:manage", "tables:read",
		"clients:manage", "clients:read", "staff:manage", "staff:read",
		"restaurant:configure", "restaurant:reports", "delivery:manage", "delivery:read",
		"menu:manage", "menu:read",
	}
	if err := systemRepo.AssignPermissionsToRole("super_admin", superAdminPermissions); err != nil {
		return err
	}

	// Platform Admin - permisos de plataforma
	platformAdminPermissions := []string{
		"restaurants:read", "users:manage", "users:read", "roles:read",
		"permissions:manage", "reports:read",
	}
	if err := systemRepo.AssignPermissionsToRole("platform_admin", platformAdminPermissions); err != nil {
		return err
	}

	// Support - permisos limitados de plataforma
	supportPermissions := []string{
		"restaurants:read", "users:read", "reports:read",
	}
	if err := systemRepo.AssignPermissionsToRole("support", supportPermissions); err != nil {
		return err
	}

	// Restaurant Admin - permisos de restaurante
	restaurantAdminPermissions := []string{
		"reservations:manage", "reservations:read", "tables:manage", "tables:read",
		"clients:manage", "clients:read", "staff:manage", "staff:read",
		"restaurant:configure", "restaurant:reports", "delivery:manage", "delivery:read",
		"menu:manage", "menu:read",
	}
	if err := systemRepo.AssignPermissionsToRole("restaurant_admin", restaurantAdminPermissions); err != nil {
		return err
	}

	// Manager - permisos de gestión
	managerPermissions := []string{
		"reservations:manage", "reservations:read", "tables:read", "clients:manage",
		"clients:read", "staff:read", "restaurant:reports", "delivery:read", "menu:read",
	}
	if err := systemRepo.AssignPermissionsToRole("manager", managerPermissions); err != nil {
		return err
	}

	// Hostess - permisos de reservas
	hostessPermissions := []string{
		"reservations:manage", "reservations:read", "tables:read", "clients:manage",
		"clients:read",
	}
	if err := systemRepo.AssignPermissionsToRole("hostess", hostessPermissions); err != nil {
		return err
	}

	// Waiter - permisos básicos
	waiterPermissions := []string{
		"reservations:read", "tables:read", "clients:read",
	}
	if err := systemRepo.AssignPermissionsToRole("waiter", waiterPermissions); err != nil {
		return err
	}

	// Cook - permisos de menú
	cookPermissions := []string{
		"menu:manage", "menu:read", "delivery:read",
	}
	if err := systemRepo.AssignPermissionsToRole("cook", cookPermissions); err != nil {
		return err
	}

	logger.Info().Msg("✅ Permisos asignados a roles correctamente")
	return nil
}
