package repository

import (
	"dbpostgres/db/models"
	"dbpostgres/pkg/log"
	"strings"

	"gorm.io/gorm"
)

// userRepository implementa UserRepository
type userRepository struct {
	db     *gorm.DB
	logger log.ILogger
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository(db *gorm.DB, logger log.ILogger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// Create crea un nuevo usuario
func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error().Err(err).Str("user_email", user.Email).Msg("Error al crear usuario")
		return err
	}
	r.logger.Debug().Str("user_email", user.Email).Msg("Usuario creado exitosamente")
	return nil
}

// GetByEmail obtiene un usuario por su email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("user_email", email).Msg("Error al obtener usuario por email")
		return nil, err
	}
	return &user, nil
}

// GetByID obtiene un usuario por su ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("user_id", id).Msg("Error al obtener usuario por ID")
		return nil, err
	}
	return &user, nil
}

// ExistsByEmail verifica si existe un usuario por su email
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Str("user_email", email).Msg("Error al verificar existencia de usuario")
		return false, err
	}
	return count > 0, nil
}

// UserExists verifica si existe algún usuario en la tabla
func (r *userRepository) UserExists() (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al verificar si existen usuarios en la tabla")
		return false, err
	}
	return count > 0, nil
}

// AssignRoles asigna roles a un usuario usando sintaxis GORM
func (r *userRepository) AssignRoles(userID uint, roleIDs []uint) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Error al buscar usuario para asignar roles")
		return err
	}

	// Eliminar roles actuales del usuario
	if err := r.db.Model(&user).Association("Roles").Clear(); err != nil {
		r.logger.Error().Err(err).Uint("user_id", userID).Msg("Error al limpiar roles actuales del usuario")
		return err
	}

	// Obtener los roles a asignar
	var roles []models.Role
	if err := r.db.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al buscar roles para asignar al usuario")
		return err
	}

	// Asignar los nuevos roles usando sintaxis GORM (insertando en tabla intermedia)
	for _, role := range roles {
		userRole := map[string]interface{}{
			"user_id": userID,
			"role_id": role.ID,
		}
		if err := r.db.Table("user_roles").Create(userRole).Error; err != nil {
			// Si ya existe la relación, ignorar el error de duplicado
			if !isUniqueConstraintError(err) {
				r.logger.Error().Err(err).Uint("user_id", userID).Uint("role_id", role.ID).Msg("Error al insertar relación en user_roles")
				return err
			}
		}
	}

	return nil
}

// isUniqueConstraintError verifica si el error es por restricción de unicidad
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE constraint failed")
}
