package repository

import (
	"dbpostgres/db/models"
	"dbpostgres/pkg/log"

	"gorm.io/gorm"
)

// systemRepository implementa SystemRepository
type systemRepository struct {
	permissionRepo PermissionRepository
	roleRepo       RoleRepository
	userRepo       UserRepository
	restaurantRepo RestaurantRepository
	tableRepo      TableRepository
	statusRepo     ReservationStatusRepository
	logger         log.ILogger
}

// NewSystemRepository crea una nueva instancia del repositorio del sistema
func NewSystemRepository(db *gorm.DB, logger log.ILogger) SystemRepository {
	return &systemRepository{
		permissionRepo: NewPermissionRepository(db, logger),
		roleRepo:       NewRoleRepository(db, logger),
		userRepo:       NewUserRepository(db, logger),
		restaurantRepo: NewRestaurantRepository(db, logger),
		tableRepo:      NewTableRepository(db, logger),
		statusRepo:     NewReservationStatusRepository(db, logger),
		logger:         logger,
	}
}

// InitializePermissions inicializa los permisos del sistema
func (r *systemRepository) InitializePermissions(permissions []models.Permission) error {
	r.logger.Info().Msg("Inicializando permisos del sistema...")

	// Verificar si todos los permisos ya existen
	allExist := true
	for _, permission := range permissions {
		exists, err := r.permissionRepo.ExistsByCode(permission.Code)
		if err != nil {
			return err
		}
		if !exists {
			allExist = false
			break
		}
	}

	if allExist {
		r.logger.Info().Int("permissions_count", len(permissions)).Msg("✅ Permisos del sistema ya existen, saltando migración de permisos")
		return nil
	}

	for _, permission := range permissions {
		exists, err := r.permissionRepo.ExistsByCode(permission.Code)
		if err != nil {
			return err
		}

		if !exists {
			if err := r.permissionRepo.Create(&permission); err != nil {
				return err
			}
		}
	}

	r.logger.Info().Int("permissions_count", len(permissions)).Msg("✅ Permisos del sistema inicializados correctamente")
	return nil
}

// InitializeRoles inicializa los roles del sistema
func (r *systemRepository) InitializeRoles(roles []models.Role) error {
	r.logger.Info().Msg("Inicializando roles del sistema...")

	// Verificar si todos los roles ya existen
	allExist := true
	for _, role := range roles {
		exists, err := r.roleRepo.ExistsByCode(role.Code)
		if err != nil {
			return err
		}
		if !exists {
			allExist = false
			break
		}
	}

	if allExist {
		r.logger.Info().Int("roles_count", len(roles)).Msg("✅ Roles del sistema ya existen, saltando migración de roles")
		return nil
	}

	for _, role := range roles {
		exists, err := r.roleRepo.ExistsByCode(role.Code)
		if err != nil {
			return err
		}

		if !exists {
			if err := r.roleRepo.Create(&role); err != nil {
				return err
			}
		}
	}

	r.logger.Info().Int("roles_count", len(roles)).Msg("✅ Roles del sistema inicializados correctamente")
	return nil
}

// InitializeUsers inicializa los usuarios del sistema
func (r *systemRepository) InitializeUsers(users []models.User) error {
	r.logger.Info().Msg("Inicializando usuarios del sistema...")

	for _, user := range users {
		exists, err := r.userRepo.ExistsByEmail(user.Email)
		if err != nil {
			return err
		}

		if !exists {
			if err := r.userRepo.Create(&user); err != nil {
				return err
			}
		}
	}

	r.logger.Info().Int("users_count", len(users)).Msg("Usuarios del sistema inicializados correctamente")
	return nil
}

// InitializeRestaurants inicializa los restaurantes del sistema
func (r *systemRepository) InitializeRestaurants(restaurants []models.Restaurant) error {
	r.logger.Info().Msg("Inicializando restaurantes del sistema...")

	for _, restaurant := range restaurants {
		exists, err := r.restaurantRepo.ExistsByCode(restaurant.Code)
		if err != nil {
			return err
		}

		if !exists {
			if err := r.restaurantRepo.Create(&restaurant); err != nil {
				return err
			}
		}
	}

	r.logger.Info().Int("restaurants_count", len(restaurants)).Msg("Restaurantes del sistema inicializados correctamente")
	return nil
}

// InitializeTables inicializa las mesas del sistema
func (r *systemRepository) InitializeTables(tables []models.Table) error {
	r.logger.Info().Msg("Inicializando mesas del sistema...")

	for _, table := range tables {
		exists, err := r.tableRepo.ExistsByRestaurantAndNumber(table.RestaurantID, table.Number)
		if err != nil {
			return err
		}

		if !exists {
			if err := r.tableRepo.Create(&table); err != nil {
				return err
			}
		}
	}

	r.logger.Info().Int("tables_count", len(tables)).Msg("Mesas del sistema inicializadas correctamente")
	return nil
}

// InitializeReservationStatuses inicializa los estados de reserva del sistema
func (r *systemRepository) InitializeReservationStatuses(statuses []models.ReservationStatus) error {
	r.logger.Info().Msg("Inicializando estados de reserva del sistema...")

	// Verificar si todos los estados ya existen
	allExist := true
	for _, status := range statuses {
		exists, err := r.statusRepo.ExistsByCode(status.Code)
		if err != nil {
			return err
		}
		if !exists {
			allExist = false
			break
		}
	}

	if allExist {
		r.logger.Info().Int("statuses_count", len(statuses)).Msg("✅ Estados de reserva ya existen, saltando migración de estados de reserva")
		return nil
	}

	for _, status := range statuses {
		exists, err := r.statusRepo.ExistsByCode(status.Code)
		if err != nil {
			return err
		}

		if !exists {
			if err := r.statusRepo.Create(&status); err != nil {
				return err
			}
		}
	}

	r.logger.Info().Int("statuses_count", len(statuses)).Msg("✅ Estados de reserva del sistema inicializados correctamente")
	return nil
}

// AssignPermissionsToRole asigna permisos a un rol por códigos
func (r *systemRepository) AssignPermissionsToRole(roleCode string, permissionCodes []string) error {
	r.logger.Info().Str("role_code", roleCode).Msg("Asignando permisos al rol...")

	// Obtener el rol
	role, err := r.roleRepo.GetByCode(roleCode)
	if err != nil {
		return err
	}
	if role == nil {
		r.logger.Error().Str("role_code", roleCode).Msg("Rol no encontrado")
		return gorm.ErrRecordNotFound
	}

	// Obtener los IDs de los permisos
	var permissionIDs []uint
	for _, permissionCode := range permissionCodes {
		permission, err := r.permissionRepo.GetByCode(permissionCode)
		if err != nil {
			return err
		}
		if permission == nil {
			r.logger.Warn().Str("permission_code", permissionCode).Msg("Permiso no encontrado, saltando...")
			continue
		}
		permissionIDs = append(permissionIDs, permission.ID)
	}

	// Asignar permisos al rol
	if err := r.roleRepo.AssignPermissions(role.ID, permissionIDs); err != nil {
		return err
	}

	r.logger.Info().Str("role_code", roleCode).Int("permissions_count", len(permissionIDs)).Msg("Permisos asignados al rol exitosamente")
	return nil
}

// AssignRolesToUser asigna roles a un usuario por códigos
func (r *systemRepository) AssignRolesToUser(userEmail string, roleCodes []string) error {
	r.logger.Info().Str("user_email", userEmail).Msg("Asignando roles al usuario...")

	// Obtener el usuario
	user, err := r.userRepo.GetByEmail(userEmail)
	if err != nil {
		return err
	}
	if user == nil {
		r.logger.Error().Str("user_email", userEmail).Msg("Usuario no encontrado")
		return gorm.ErrRecordNotFound
	}

	// Obtener los IDs de los roles
	var roleIDs []uint
	for _, roleCode := range roleCodes {
		role, err := r.roleRepo.GetByCode(roleCode)
		if err != nil {
			return err
		}
		if role == nil {
			r.logger.Warn().Str("role_code", roleCode).Msg("Rol no encontrado, saltando...")
			continue
		}
		roleIDs = append(roleIDs, role.ID)
	}

	// Asignar roles al usuario
	if err := r.userRepo.AssignRoles(user.ID, roleIDs); err != nil {
		return err
	}

	r.logger.Info().Str("user_email", userEmail).Int("roles_count", len(roleIDs)).Msg("Roles asignados al usuario exitosamente")
	return nil
}

// ExistsByEmail verifica si existe un usuario por su email
func (r *systemRepository) ExistsByEmail(email string) (bool, error) {
	return r.userRepo.ExistsByEmail(email)
}

// Create crea un nuevo usuario
func (r *systemRepository) Create(user *models.User) error {
	return r.userRepo.Create(user)
}

// UserExists verifica si existe algún usuario en la tabla
func (r *systemRepository) UserExists() (bool, error) {
	return r.userRepo.UserExists()
}

// ExistsByCode verifica si existe un restaurante por su código
func (r *systemRepository) ExistsByCode(code string) (bool, error) {
	return r.restaurantRepo.ExistsByCode(code)
}

// CreateRestaurant crea un nuevo restaurante
func (r *systemRepository) CreateRestaurant(restaurant *models.Restaurant) error {
	return r.restaurantRepo.Create(restaurant)
}

// ExistsByRestaurantAndNumber verifica si existe una mesa por restaurante y número
func (r *systemRepository) ExistsByRestaurantAndNumber(restaurantID uint, number int) (bool, error) {
	return r.tableRepo.ExistsByRestaurantAndNumber(restaurantID, number)
}

// GetRoleByCode obtiene un rol por su código
func (r *systemRepository) GetRoleByCode(code string) (*models.Role, error) {
	return r.roleRepo.GetByCode(code)
}

// GetPermissionByCode obtiene un permiso por su código
func (r *systemRepository) GetPermissionByCode(code string) (*models.Permission, error) {
	return r.permissionRepo.GetByCode(code)
}

// GetAllReservationStatuses obtiene todos los estados de reserva
func (r *systemRepository) GetAllReservationStatuses() ([]models.ReservationStatus, error) {
	return r.statusRepo.GetAll()
}
