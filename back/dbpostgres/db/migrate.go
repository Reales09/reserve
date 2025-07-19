package db

import (
	"log"

	"dbpostgres/db/models"

	"gorm.io/gorm"
)

// MigrateDB ejecuta todas las migraciones
func MigrateDB(db *gorm.DB) error {
	log.Println("Iniciando migraciones...")

	// Migrar modelos
	err := db.AutoMigrate(
		&models.Restaurant{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RestaurantStaff{},
		&models.Client{},
		&models.Table{},
		&models.Reservation{},
		&models.ReservationStatus{},
		&models.ReservationStatusHistory{},
	)
	if err != nil {
		return err
	}

	// Inicializar datos del sistema
	if err := initializeSystemData(db); err != nil {
		return err
	}

	log.Println("Migraciones completadas exitosamente")
	return nil
}

// initializeSystemData inicializa roles y permisos del sistema
func initializeSystemData(db *gorm.DB) error {
	log.Println("Inicializando datos del sistema...")

	// Crear permisos del sistema
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

	// Crear permisos si no existen
	for _, permission := range permissions {
		var existingPermission models.Permission
		if err := db.Where("code = ?", permission.Code).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permission).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// Crear roles del sistema
	roles := []models.Role{
		{
			Name: "Super Administrador", Code: models.ROLE_SUPER_ADMIN, Description: "Control total del sistema",
			Level: 1, IsSystem: true, Scope: models.SCOPE_PLATFORM,
		},
		{
			Name: "Administrador de Plataforma", Code: models.ROLE_PLATFORM_ADMIN, Description: "Administración de la plataforma",
			Level: 2, IsSystem: true, Scope: models.SCOPE_PLATFORM,
		},
		{
			Name: "Operador de Plataforma", Code: models.ROLE_PLATFORM_OPERATOR, Description: "Operaciones básicas de la plataforma",
			Level: 3, IsSystem: true, Scope: models.SCOPE_PLATFORM,
		},
		{
			Name: "Propietario de Restaurante", Code: models.ROLE_RESTAURANT_OWNER, Description: "Propietario del restaurante",
			Level: 1, IsSystem: true, Scope: models.SCOPE_RESTAURANT,
		},
		{
			Name: "Gerente de Restaurante", Code: models.ROLE_RESTAURANT_MANAGER, Description: "Gerente del restaurante",
			Level: 2, IsSystem: true, Scope: models.SCOPE_RESTAURANT,
		},
		{
			Name: "Personal de Restaurante", Code: models.ROLE_RESTAURANT_STAFF, Description: "Personal general del restaurante",
			Level: 3, IsSystem: true, Scope: models.SCOPE_RESTAURANT,
		},
		{
			Name: "Mesero", Code: models.ROLE_WAITER, Description: "Mesero del restaurante",
			Level: 4, IsSystem: true, Scope: models.SCOPE_RESTAURANT,
		},
		{
			Name: "Anfitrión", Code: models.ROLE_HOST, Description: "Anfitrión del restaurante",
			Level: 4, IsSystem: true, Scope: models.SCOPE_RESTAURANT,
		},
	}

	// Crear roles si no existen
	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("code = ?", role.Code).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&role).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	// Asignar permisos a roles
	if err := assignPermissionsToRoles(db); err != nil {
		return err
	}

	log.Println("Datos del sistema inicializados correctamente")
	return nil
}

// assignPermissionsToRoles asigna permisos a los roles del sistema
func assignPermissionsToRoles(db *gorm.DB) error {
	// Obtener todos los roles y permisos
	var roles []models.Role
	var permissions []models.Permission

	if err := db.Find(&roles).Error; err != nil {
		return err
	}
	if err := db.Find(&permissions).Error; err != nil {
		return err
	}

	// Crear mapa de permisos por código
	permissionMap := make(map[string]models.Permission)
	for _, p := range permissions {
		permissionMap[p.Code] = p
	}

	// Asignar permisos a cada rol
	for _, role := range roles {
		var rolePermissions []models.Permission

		switch role.Code {
		case models.ROLE_SUPER_ADMIN:
			// Super admin tiene todos los permisos
			rolePermissions = permissions

		case models.ROLE_PLATFORM_ADMIN:
			// Admin de plataforma tiene permisos de plataforma
			for _, p := range permissions {
				if p.Scope == models.SCOPE_PLATFORM {
					rolePermissions = append(rolePermissions, p)
				}
			}

		case models.ROLE_PLATFORM_OPERATOR:
			// Operador tiene permisos limitados de plataforma
			operatorPerms := []string{
				"restaurants:read", "users:read", "roles:read", "reports:read",
			}
			for _, permCode := range operatorPerms {
				if perm, exists := permissionMap[permCode]; exists {
					rolePermissions = append(rolePermissions, perm)
				}
			}

		case models.ROLE_RESTAURANT_OWNER:
			// Propietario tiene todos los permisos del restaurante
			for _, p := range permissions {
				if p.Scope == models.SCOPE_RESTAURANT || p.Scope == models.SCOPE_BOTH {
					rolePermissions = append(rolePermissions, p)
				}
			}

		case models.ROLE_RESTAURANT_MANAGER:
			// Gerente tiene permisos de gestión del restaurante
			managerPerms := []string{
				"reservations:manage", "reservations:read", "tables:manage", "tables:read",
				"clients:manage", "clients:read", "staff:manage", "staff:read",
				"restaurant:configure", "restaurant:reports", "delivery:manage", "delivery:read",
				"menu:manage", "menu:read",
			}
			for _, permCode := range managerPerms {
				if perm, exists := permissionMap[permCode]; exists {
					rolePermissions = append(rolePermissions, perm)
				}
			}

		case models.ROLE_RESTAURANT_STAFF:
			// Staff tiene permisos básicos
			staffPerms := []string{
				"reservations:read", "tables:read", "clients:read", "staff:read",
				"restaurant:reports", "delivery:read", "menu:read",
			}
			for _, permCode := range staffPerms {
				if perm, exists := permissionMap[permCode]; exists {
					rolePermissions = append(rolePermissions, perm)
				}
			}

		case models.ROLE_WAITER:
			// Mesero tiene permisos específicos
			waiterPerms := []string{
				"reservations:read", "tables:read", "clients:read", "menu:read",
			}
			for _, permCode := range waiterPerms {
				if perm, exists := permissionMap[permCode]; exists {
					rolePermissions = append(rolePermissions, perm)
				}
			}

		case models.ROLE_HOST:
			// Anfitrión tiene permisos de reservas
			hostPerms := []string{
				"reservations:manage", "reservations:read", "tables:read", "clients:manage", "clients:read",
			}
			for _, permCode := range hostPerms {
				if perm, exists := permissionMap[permCode]; exists {
					rolePermissions = append(rolePermissions, perm)
				}
			}
		}

		// Asignar permisos al rol
		if len(rolePermissions) > 0 {
			if err := db.Model(&role).Association("Permissions").Replace(rolePermissions); err != nil {
				return err
			}
		}
	}

	return nil
}
