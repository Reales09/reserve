package repository

import (
	"dbpostgres/db/models"
	"dbpostgres/pkg/log"

	"gorm.io/gorm"
)

// Repository interfaces para cada entidad

// PermissionRepository define las operaciones para permisos
type PermissionRepository interface {
	Create(permission *models.Permission) error
	GetByCode(code string) (*models.Permission, error)
	GetAll() ([]models.Permission, error)
	ExistsByCode(code string) (bool, error)
}

// RoleRepository define las operaciones para roles
type RoleRepository interface {
	Create(role *models.Role) error
	GetByCode(code string) (*models.Role, error)
	GetAll() ([]models.Role, error)
	ExistsByCode(code string) (bool, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
}

// UserRepository define las operaciones de usuarios
type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	ExistsByEmail(email string) (bool, error)
	AssignRoles(userID uint, roleIDs []uint) error
	UserExists() (bool, error)
}

// RestaurantRepository define las operaciones para restaurantes
type RestaurantRepository interface {
	Create(restaurant *models.Restaurant) error
	GetByCode(code string) (*models.Restaurant, error)
	GetByID(id uint) (*models.Restaurant, error)
	ExistsByCode(code string) (bool, error)
}

// TableRepository define las operaciones para mesas
type TableRepository interface {
	Create(table *models.Table) error
	GetByRestaurantAndNumber(restaurantID uint, number int) (*models.Table, error)
	GetByRestaurant(restaurantID uint) ([]models.Table, error)
	ExistsByRestaurantAndNumber(restaurantID uint, number int) (bool, error)
}

// ReservationStatusRepository define las operaciones para estados de reserva
type ReservationStatusRepository interface {
	Create(status *models.ReservationStatus) error
	GetByCode(code string) (*models.ReservationStatus, error)
	GetAll() ([]models.ReservationStatus, error)
	ExistsByCode(code string) (bool, error)
}

// SystemRepository define las operaciones del sistema para migraciones
type SystemRepository interface {
	// Métodos específicos para migraciones del sistema
	InitializePermissions(permissions []models.Permission) error
	InitializeRoles(roles []models.Role) error
	InitializeUsers(users []models.User) error
	InitializeRestaurants(restaurants []models.Restaurant) error
	InitializeTables(tables []models.Table) error
	InitializeReservationStatuses(statuses []models.ReservationStatus) error
	AssignPermissionsToRole(roleCode string, permissionCodes []string) error
	AssignRolesToUser(userEmail string, roleCodes []string) error

	// Métodos de usuario para migraciones
	ExistsByEmail(email string) (bool, error)
	Create(user *models.User) error
	UserExists() (bool, error)

	// Métodos de restaurante para migraciones
	ExistsByCode(code string) (bool, error)
	CreateRestaurant(restaurant *models.Restaurant) error

	// Métodos de mesa para migraciones
	ExistsByRestaurantAndNumber(restaurantID uint, number int) (bool, error)

	// Métodos de consulta para roles y permisos
	GetRoleByCode(code string) (*models.Role, error)
	GetPermissionByCode(code string) (*models.Permission, error)

	// Métodos de estados de reserva
	GetAllReservationStatuses() ([]models.ReservationStatus, error)
}

// RepositoryManager maneja todas las instancias de repositorios
type RepositoryManager struct {
	PermissionRepository
	RoleRepository
	UserRepository
	RestaurantRepository
	TableRepository
	ReservationStatusRepository
}

// NewRepositoryManager crea una nueva instancia del manager de repositorios
func NewRepositoryManager(db *gorm.DB, logger log.ILogger) *RepositoryManager {
	return &RepositoryManager{
		PermissionRepository:        NewPermissionRepository(db, logger),
		RoleRepository:              NewRoleRepository(db, logger),
		UserRepository:              NewUserRepository(db, logger),
		RestaurantRepository:        NewRestaurantRepository(db, logger),
		TableRepository:             NewTableRepository(db, logger),
		ReservationStatusRepository: NewReservationStatusRepository(db, logger),
	}
}
