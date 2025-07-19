package repository

import (
	"dbpostgres/db/models"
	"dbpostgres/pkg/log"

	"gorm.io/gorm"
)

// tableRepository implementa TableRepository
type tableRepository struct {
	db     *gorm.DB
	logger log.ILogger
}

// NewTableRepository crea una nueva instancia del repositorio de mesas
func NewTableRepository(db *gorm.DB, logger log.ILogger) TableRepository {
	return &tableRepository{
		db:     db,
		logger: logger,
	}
}

// Create crea una nueva mesa
func (r *tableRepository) Create(table *models.Table) error {
	if err := r.db.Create(table).Error; err != nil {
		r.logger.Error().Err(err).Uint("restaurant_id", table.RestaurantID).Int("table_number", table.Number).Msg("Error al crear mesa")
		return err
	}
	r.logger.Debug().Uint("restaurant_id", table.RestaurantID).Int("table_number", table.Number).Msg("Mesa creada exitosamente")
	return nil
}

// GetByRestaurantAndNumber obtiene una mesa por restaurante y número
func (r *tableRepository) GetByRestaurantAndNumber(restaurantID uint, number int) (*models.Table, error) {
	var table models.Table
	if err := r.db.Where("restaurant_id = ? AND number = ?", restaurantID, number).First(&table).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Uint("restaurant_id", restaurantID).Int("table_number", number).Msg("Error al obtener mesa por restaurante y número")
		return nil, err
	}
	return &table, nil
}

// GetByRestaurant obtiene todas las mesas de un restaurante
func (r *tableRepository) GetByRestaurant(restaurantID uint) ([]models.Table, error) {
	var tables []models.Table
	if err := r.db.Where("restaurant_id = ?", restaurantID).Find(&tables).Error; err != nil {
		r.logger.Error().Err(err).Uint("restaurant_id", restaurantID).Msg("Error al obtener mesas del restaurante")
		return nil, err
	}
	return tables, nil
}

// ExistsByRestaurantAndNumber verifica si existe una mesa por restaurante y número
func (r *tableRepository) ExistsByRestaurantAndNumber(restaurantID uint, number int) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Table{}).Where("restaurant_id = ? AND number = ?", restaurantID, number).Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Uint("restaurant_id", restaurantID).Int("table_number", number).Msg("Error al verificar existencia de mesa")
		return false, err
	}
	return count > 0, nil
}
