package services

import (
	"errors"

	"dbpostgres/db/models"

	"gorm.io/gorm"
)

// UserService maneja las operaciones relacionadas con usuarios
type UserService struct {
	db *gorm.DB
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser crea un nuevo usuario
func (us *UserService) CreateUser(user *models.User, roleCodes []string, restaurantIDs []uint) error {
	// Verificar que el email sea único
	var existingUser models.User
	if err := us.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("el email ya está registrado")
	}

	// Iniciar transacción
	tx := us.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Crear usuario
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Asignar roles
	if len(roleCodes) > 0 {
		var roles []models.Role
		if err := tx.Where("code IN ?", roleCodes).Find(&roles).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Asignar restaurantes
	if len(restaurantIDs) > 0 {
		var restaurants []models.Restaurant
		if err := tx.Where("id IN ?", restaurantIDs).Find(&restaurants).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(user).Association("Restaurants").Append(restaurants); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetUserByID obtiene un usuario por ID
func (us *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	err := us.db.Preload("Roles.Permissions").Preload("Restaurant").Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail obtiene un usuario por email
func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := us.db.Preload("Roles.Permissions").Preload("Restaurant").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser actualiza un usuario
func (us *UserService) UpdateUser(userID uint, updates map[string]interface{}) error {
	return us.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

// DeleteUser elimina un usuario
func (us *UserService) DeleteUser(userID uint) error {
	return us.db.Delete(&models.User{}, userID).Error
}

// AssignRolesToUser asigna roles a un usuario
func (us *UserService) AssignRolesToUser(userID uint, roleCodes []string) error {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return err
	}

	var roles []models.Role
	if err := us.db.Where("code IN ?", roleCodes).Find(&roles).Error; err != nil {
		return err
	}

	return us.db.Model(&user).Association("Roles").Replace(roles)
}

// AssignRestaurantsToUser asigna restaurantes a un usuario
func (us *UserService) AssignRestaurantsToUser(userID uint, restaurantIDs []uint) error {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return err
	}

	var restaurants []models.Restaurant
	if err := us.db.Where("id IN ?", restaurantIDs).Find(&restaurants).Error; err != nil {
		return err
	}

	return us.db.Model(&user).Association("Restaurants").Replace(restaurants)
}

// AddRestaurantToUser agrega un restaurante a un usuario
func (us *UserService) AddRestaurantToUser(userID uint, restaurantID uint) error {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return err
	}

	var restaurant models.Restaurant
	if err := us.db.First(&restaurant, restaurantID).Error; err != nil {
		return err
	}

	return us.db.Model(&user).Association("Restaurants").Append(&restaurant)
}

// RemoveRestaurantFromUser remueve un restaurante de un usuario
func (us *UserService) RemoveRestaurantFromUser(userID uint, restaurantID uint) error {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return err
	}

	var restaurant models.Restaurant
	if err := us.db.First(&restaurant, restaurantID).Error; err != nil {
		return err
	}

	return us.db.Model(&user).Association("Restaurants").Delete(&restaurant)
}

// GetUserRestaurants obtiene los restaurantes de un usuario
func (us *UserService) GetUserRestaurants(userID uint) ([]models.Restaurant, error) {
	var user models.User
	if err := us.db.Preload("Restaurants").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Restaurants, nil
}

// GetUsersByRestaurant obtiene usuarios de un restaurante específico
func (us *UserService) GetUsersByRestaurant(restaurantID uint) ([]models.User, error) {
	var users []models.User
	err := us.db.Preload("Roles").Preload("Restaurants").
		Joins("JOIN user_restaurants ON users.id = user_restaurants.user_id").
		Where("user_restaurants.restaurant_id = ?", restaurantID).Find(&users).Error
	return users, err
}

// GetUsersByRole obtiene usuarios con un rol específico
func (us *UserService) GetUsersByRole(roleCode string) ([]models.User, error) {
	var users []models.User
	err := us.db.Preload("Roles").Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("roles.code = ?", roleCode).Find(&users).Error
	return users, err
}

// RestaurantService maneja las operaciones relacionadas con restaurantes
type RestaurantService struct {
	db *gorm.DB
}

// NewRestaurantService crea una nueva instancia del servicio de restaurantes
func NewRestaurantService(db *gorm.DB) *RestaurantService {
	return &RestaurantService{db: db}
}

// CreateRestaurant crea un nuevo restaurante
func (us *RestaurantService) CreateRestaurant(restaurant *models.Restaurant) error {
	// Verificar que el código sea único
	var existingRestaurant models.Restaurant
	if err := us.db.Where("code = ?", restaurant.Code).First(&existingRestaurant).Error; err == nil {
		return errors.New("el código del restaurante ya existe")
	}

	return us.db.Create(restaurant).Error
}

// GetRestaurantByID obtiene un restaurante por ID
func (us *RestaurantService) GetRestaurantByID(restaurantID uint) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := us.db.Preload("Users").Preload("Tables").Where("id = ?", restaurantID).First(&restaurant).Error
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

// GetRestaurantByCode obtiene un restaurante por código
func (us *RestaurantService) GetRestaurantByCode(code string) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := us.db.Preload("Users").Preload("Tables").Where("code = ?", code).First(&restaurant).Error
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

// GetRestaurantByDomain obtiene un restaurante por dominio personalizado
func (us *RestaurantService) GetRestaurantByDomain(domain string) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := us.db.Preload("Users").Preload("Tables").Where("custom_domain = ?", domain).First(&restaurant).Error
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

// UpdateRestaurant actualiza un restaurante
func (us *RestaurantService) UpdateRestaurant(restaurantID uint, updates map[string]interface{}) error {
	return us.db.Model(&models.Restaurant{}).Where("id = ?", restaurantID).Updates(updates).Error
}

// DeleteRestaurant elimina un restaurante
func (us *RestaurantService) DeleteRestaurant(restaurantID uint) error {
	return us.db.Delete(&models.Restaurant{}, restaurantID).Error
}

// GetAllRestaurants obtiene todos los restaurantes
func (us *RestaurantService) GetAllRestaurants() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := us.db.Find(&restaurants).Error
	return restaurants, err
}

// GetActiveRestaurants obtiene solo restaurantes activos
func (us *RestaurantService) GetActiveRestaurants() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := us.db.Where("is_active = ?", true).Find(&restaurants).Error
	return restaurants, err
}

// RoleService maneja las operaciones relacionadas con roles
type RoleService struct {
	db *gorm.DB
}

// NewRoleService crea una nueva instancia del servicio de roles
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

// GetAllRoles obtiene todos los roles
func (rs *RoleService) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	err := rs.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

// GetRoleByCode obtiene un rol por código
func (rs *RoleService) GetRoleByCode(code string) (*models.Role, error) {
	var role models.Role
	err := rs.db.Preload("Permissions").Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRolesByScope obtiene roles por scope
func (rs *RoleService) GetRolesByScope(scope string) ([]models.Role, error) {
	var roles []models.Role
	err := rs.db.Preload("Permissions").Where("scope = ?", scope).Find(&roles).Error
	return roles, err
}

// CreateRole crea un nuevo rol
func (rs *RoleService) CreateRole(role *models.Role) error {
	// Verificar que el código sea único
	var existingRole models.Role
	if err := rs.db.Where("code = ?", role.Code).First(&existingRole).Error; err == nil {
		return errors.New("el código del rol ya existe")
	}

	return rs.db.Create(role).Error
}

// UpdateRole actualiza un rol
func (rs *RoleService) UpdateRole(roleID uint, updates map[string]interface{}) error {
	// Verificar que no sea un rol del sistema
	var role models.Role
	if err := rs.db.First(&role, roleID).Error; err != nil {
		return err
	}

	if role.IsSystem {
		return errors.New("no se puede modificar un rol del sistema")
	}

	return rs.db.Model(&models.Role{}).Where("id = ?", roleID).Updates(updates).Error
}

// DeleteRole elimina un rol
func (rs *RoleService) DeleteRole(roleID uint) error {
	// Verificar que no sea un rol del sistema
	var role models.Role
	if err := rs.db.First(&role, roleID).Error; err != nil {
		return err
	}

	if role.IsSystem {
		return errors.New("no se puede eliminar un rol del sistema")
	}

	return rs.db.Delete(&models.Role{}, roleID).Error
}

// AssignPermissionsToRole asigna permisos a un rol
func (rs *RoleService) AssignPermissionsToRole(roleID uint, permissionCodes []string) error {
	var role models.Role
	if err := rs.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var permissions []models.Permission
	if err := rs.db.Where("code IN ?", permissionCodes).Find(&permissions).Error; err != nil {
		return err
	}

	return rs.db.Model(&role).Association("Permissions").Replace(permissions)
}

// PermissionService maneja las operaciones relacionadas con permisos
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService crea una nueva instancia del servicio de permisos
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// GetAllPermissions obtiene todos los permisos
func (ps *PermissionService) GetAllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	err := ps.db.Find(&permissions).Error
	return permissions, err
}

// GetPermissionsByScope obtiene permisos por scope
func (ps *PermissionService) GetPermissionsByScope(scope string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := ps.db.Where("scope = ?", scope).Find(&permissions).Error
	return permissions, err
}

// GetPermissionsByResource obtiene permisos por recurso
func (ps *PermissionService) GetPermissionsByResource(resource string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := ps.db.Where("resource = ?", resource).Find(&permissions).Error
	return permissions, err
}

// CreatePermission crea un nuevo permiso
func (ps *PermissionService) CreatePermission(permission *models.Permission) error {
	// Verificar que el código sea único
	var existingPermission models.Permission
	if err := ps.db.Where("code = ?", permission.Code).First(&existingPermission).Error; err == nil {
		return errors.New("el código del permiso ya existe")
	}

	return ps.db.Create(permission).Error
}

// UpdatePermission actualiza un permiso
func (ps *PermissionService) UpdatePermission(permissionID uint, updates map[string]interface{}) error {
	return ps.db.Model(&models.Permission{}).Where("id = ?", permissionID).Updates(updates).Error
}

// DeletePermission elimina un permiso
func (ps *PermissionService) DeletePermission(permissionID uint) error {
	return ps.db.Delete(&models.Permission{}, permissionID).Error
}
