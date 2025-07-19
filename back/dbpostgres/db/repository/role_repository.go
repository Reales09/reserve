package repository

import (
	"dbpostgres/db/models"
	"dbpostgres/pkg/log"

	"gorm.io/gorm"
)

// roleRepository implementa RoleRepository
type roleRepository struct {
	db     *gorm.DB
	logger log.ILogger
}

// NewRoleRepository crea una nueva instancia del repositorio de roles
func NewRoleRepository(db *gorm.DB, logger log.ILogger) RoleRepository {
	return &roleRepository{
		db:     db,
		logger: logger,
	}
}

// Create crea un nuevo rol
func (r *roleRepository) Create(role *models.Role) error {
	if err := r.db.Create(role).Error; err != nil {
		r.logger.Error().Err(err).Str("role_code", role.Code).Msg("Error al crear rol")
		return err
	}
	r.logger.Debug().Str("role_code", role.Code).Msg("Rol creado exitosamente")
	return nil
}

// GetByCode obtiene un rol por su código
func (r *roleRepository) GetByCode(code string) (*models.Role, error) {
	var role models.Role
	if err := r.db.Where("code = ?", code).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("role_code", code).Msg("Error al obtener rol por código")
		return nil, err
	}
	return &role, nil
}

// GetAll obtiene todos los roles
func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	if err := r.db.Find(&roles).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener todos los roles")
		return nil, err
	}
	return roles, nil
}

// ExistsByCode verifica si existe un rol por su código
func (r *roleRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Role{}).Where("code = ?", code).Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Str("role_code", code).Msg("Error al verificar existencia de rol")
		return false, err
	}
	return count > 0, nil
}

// AssignPermissions asigna permisos a un rol usando sintaxis GORM
func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	var role models.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		r.logger.Error().Err(err).Uint("role_id", roleID).Msg("Error al buscar rol para asignar permisos")
		return err
	}

	var permissions []models.Permission
	if err := r.db.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al buscar permisos para asignar al rol")
		return err
	}

	// Usar GORM para asociar los permisos al rol (many2many)
	if err := r.db.Model(&role).Association("Permissions").Append(&permissions); err != nil {
		r.logger.Error().Err(err).Uint("role_id", roleID).Msg("Error al asignar permisos al rol")
		return err
	}

	r.logger.Debug().Uint("role_id", roleID).Int("permissions_count", len(permissionIDs)).Msg("Permisos asignados al rol exitosamente")
	return nil
}
