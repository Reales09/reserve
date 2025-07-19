package repository

import (
	"dbpostgres/db/models"
	"dbpostgres/pkg/log"

	"gorm.io/gorm"
)

// reservationStatusRepository implementa ReservationStatusRepository
type reservationStatusRepository struct {
	db     *gorm.DB
	logger log.ILogger
}

// NewReservationStatusRepository crea una nueva instancia del repositorio de estados de reserva
func NewReservationStatusRepository(db *gorm.DB, logger log.ILogger) ReservationStatusRepository {
	return &reservationStatusRepository{
		db:     db,
		logger: logger,
	}
}

// Create crea un nuevo estado de reserva
func (r *reservationStatusRepository) Create(status *models.ReservationStatus) error {
	if err := r.db.Create(status).Error; err != nil {
		r.logger.Error().Err(err).Str("status_code", status.Code).Msg("Error al crear estado de reserva")
		return err
	}
	r.logger.Debug().Str("status_code", status.Code).Msg("Estado de reserva creado exitosamente")
	return nil
}

// GetByCode obtiene un estado de reserva por su código
func (r *reservationStatusRepository) GetByCode(code string) (*models.ReservationStatus, error) {
	var status models.ReservationStatus
	if err := r.db.Where("code = ?", code).First(&status).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("status_code", code).Msg("Error al obtener estado de reserva por código")
		return nil, err
	}
	return &status, nil
}

// GetAll obtiene todos los estados de reserva
func (r *reservationStatusRepository) GetAll() ([]models.ReservationStatus, error) {
	var statuses []models.ReservationStatus
	if err := r.db.Find(&statuses).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener todos los estados de reserva")
		return nil, err
	}
	return statuses, nil
}

// ExistsByCode verifica si existe un estado de reserva por su código
func (r *reservationStatusRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.ReservationStatus{}).Where("code = ?", code).Count(&count).Error; err != nil {
		r.logger.Error().Err(err).Str("status_code", code).Msg("Error al verificar existencia de estado de reserva")
		return false, err
	}
	return count > 0, nil
}
