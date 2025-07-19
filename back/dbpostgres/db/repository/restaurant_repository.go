package repository

import (
	"dbpostgres/db/models"
	"dbpostgres/pkg/log"

	"gorm.io/gorm"
)

// restaurantRepository implementa RestaurantRepository
type restaurantRepository struct {
	db     *gorm.DB
	logger log.ILogger
}

// NewRestaurantRepository crea una nueva instancia del repositorio de restaurantes
func NewRestaurantRepository(db *gorm.DB, logger log.ILogger) RestaurantRepository {
	return &restaurantRepository{
		db:     db,
		logger: logger,
	}
}

// Create crea un nuevo restaurante
func (r *restaurantRepository) Create(restaurant *models.Restaurant) error {
	if err := r.db.Create(restaurant).Error; err != nil {
		r.logger.Error().Err(err).Str("restaurant_code", restaurant.Code).Msg("Error al crear restaurante")
		return err
	}
	r.logger.Debug().Str("restaurant_code", restaurant.Code).Msg("Restaurante creado exitosamente")
	return nil
}

// GetByCode obtiene un restaurante por su código
func (r *restaurantRepository) GetByCode(code string) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.db.Where("code = ?", code).First(&restaurant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("restaurant_code", code).Msg("Error al obtener restaurante por código")
		return nil, err
	}
	return &restaurant, nil
}

// GetByID obtiene un restaurante por su ID
func (r *restaurantRepository) GetByID(id uint) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.db.First(&restaurant, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("restaurant_id", id).Msg("Error al obtener restaurante por ID")
		return nil, err
	}
	return &restaurant, nil
}

// ExistsByCode verifica si existe un restaurante por su código
func (r *restaurantRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Restaurant{}).Where("code = ?", code).Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Str("restaurant_code", code).Msg("Error al verificar existencia de restaurante")
		return false, err
	}
	return count > 0, nil
}
