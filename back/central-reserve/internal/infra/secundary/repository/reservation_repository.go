package repository

import (
	"central_reserve/internal/domain/dtos"
	"central_reserve/internal/domain/entities"
	"central_reserve/internal/domain/ports"
	"central_reserve/internal/infra/secundary/repository/db"
	"central_reserve/internal/pkg/log"
	"context"
	"fmt"
	"time"
)

// ReservationRepository implementa ports.IReservationRepository
type ReservationRepository struct {
	database db.IDatabase
	logger   log.ILogger
}

// NewReservationRepository crea una nueva instancia del repositorio de reservas
func NewReservationRepository(database db.IDatabase, logger log.ILogger) ports.IReservationRepository {
	return &ReservationRepository{
		database: database,
		logger:   logger,
	}
}

// CreateReserve crea una nueva reserva
func (r *ReservationRepository) CreateReserve(ctx context.Context, reserve entities.Reservation) (string, error) {
	if err := r.database.Conn(ctx).Table("reservation").Create(&reserve).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al crear reserva")
		return "", err
	}
	return fmt.Sprintf("Reserva creada con ID: %d", reserve.Id), nil
}

// GetLatestReservationByClient obtiene la última reserva de un cliente
func (r *ReservationRepository) GetLatestReservationByClient(ctx context.Context, clientID uint) (*entities.Reservation, error) {
	var reservation entities.Reservation
	if err := r.database.Conn(ctx).Table("reservation").Where("client_id = ?", clientID).Order("created_at DESC").First(&reservation).Error; err != nil {
		r.logger.Error().Uint("client_id", clientID).Msg("Error al obtener última reserva del cliente")
		return nil, err
	}
	return &reservation, nil
}

// GetReserves obtiene las reservas con filtros opcionales
func (r *ReservationRepository) GetReserves(ctx context.Context, statusID *uint, clientID *uint, tableID *uint, startDate *time.Time, endDate *time.Time) ([]dtos.ReserveDetailDTO, error) {
	var tempResults []dtos.ReserveDetailDTO

	query := r.database.Conn(ctx).
		Table("reservation").
		Select(`
			reservation.id as reserva_id,
			reservation.start_at,
			reservation.end_at,
			reservation.number_of_guests,
			reservation.created_at as reserva_creada,
			reservation.updated_at as reserva_actualizada,
			reservation_status.code as estado_codigo,
			reservation_status.name as estado_nombre,
			client.id as cliente_id,
			client.name as cliente_nombre,
			client.email as cliente_email,
			client.phone as cliente_telefono,
			client.dni as cliente_dni,
			table.id as mesa_id,
			table.number as mesa_numero,
			table.capacity as mesa_capacidad,
			business.id as negocio_id,
			business.name as negocio_nombre,
			business.code as negocio_codigo,
			business.address as negocio_direccion,
			user.id as usuario_id,
			user.name as usuario_nombre,
			user.email as usuario_email
		`).
		Joins("LEFT JOIN reservation_status ON reservation.status_id = reservation_status.id").
		Joins("LEFT JOIN client ON reservation.client_id = client.id").
		Joins("LEFT JOIN table ON reservation.table_id = table.id").
		Joins("LEFT JOIN business ON reservation.business_id = business.id").
		Joins("LEFT JOIN user ON reservation.created_by_user_id = user.id")

	// Aplicar filtros
	if statusID != nil {
		query = query.Where("reservation.status_id = ?", *statusID)
	}
	if clientID != nil {
		query = query.Where("reservation.client_id = ?", *clientID)
	}
	if tableID != nil {
		query = query.Where("reservation.table_id = ?", *tableID)
	}
	if startDate != nil {
		query = query.Where("reservation.start_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("reservation.end_at <= ?", *endDate)
	}

	if err := query.Find(&tempResults).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener reservas")
		return []dtos.ReserveDetailDTO{}, nil
	}

	// Inicializar el slice de historial para cada resultado
	for i := range tempResults {
		tempResults[i].StatusHistory = []entities.ReservationStatusHistory{}
	}

	// Obtener el historial de estados para cada reserva
	if len(tempResults) > 0 {
		var historyResults []entities.ReservationStatusHistory
		reservationIDs := make([]uint, len(tempResults))
		for i, result := range tempResults {
			reservationIDs[i] = result.ReservaID
		}

		if err := r.database.Conn(ctx).
			Table("reservation_status_history").
			Select(`
				reservation_status_history.*,
				reservation_status.code as status_code,
				reservation_status.name as status_name,
				user.name as changed_by_user
			`).
			Joins("LEFT JOIN reservation_status ON reservation_status_history.status_id = reservation_status.id").
			Joins("LEFT JOIN user ON reservation_status_history.changed_by_user_id = user.id").
			Where("reservation_status_history.reservation_id IN ?", reservationIDs).
			Order("reservation_status_history.created_at ASC").
			Find(&historyResults).Error; err != nil {
			r.logger.Error().Err(err).Msg("Error al obtener historial de estados")
		} else {
			// Agrupar historial por reservation_id
			historyMap := make(map[uint][]entities.ReservationStatusHistory)
			for _, history := range historyResults {
				historyMap[history.ReservationID] = append(historyMap[history.ReservationID], history)
			}

			// Asignar historial a cada reserva
			for i := range tempResults {
				if history, exists := historyMap[tempResults[i].ReservaID]; exists {
					tempResults[i].StatusHistory = history
				} else {
					tempResults[i].StatusHistory = []entities.ReservationStatusHistory{}
				}
			}
		}
	}

	return tempResults, nil
}

// GetReserveByID obtiene una reserva por su ID
func (r *ReservationRepository) GetReserveByID(ctx context.Context, id uint) (*dtos.ReserveDetailDTO, error) {
	var tempResult dtos.ReserveDetailDTO

	err := r.database.Conn(ctx).
		Table("reservation").
		Select(`
			reservation.id as reserva_id,
			reservation.start_at,
			reservation.end_at,
			reservation.number_of_guests,
			reservation.created_at as reserva_creada,
			reservation.updated_at as reserva_actualizada,
			reservation_status.code as estado_codigo,
			reservation_status.name as estado_nombre,
			client.id as cliente_id,
			client.name as cliente_nombre,
			client.email as cliente_email,
			client.phone as cliente_telefono,
			client.dni as cliente_dni,
			table.id as mesa_id,
			table.number as mesa_numero,
			table.capacity as mesa_capacidad,
			business.id as negocio_id,
			business.name as negocio_nombre,
			business.code as negocio_codigo,
			business.address as negocio_direccion,
			user.id as usuario_id,
			user.name as usuario_nombre,
			user.email as usuario_email
		`).
		Joins("LEFT JOIN reservation_status ON reservation.status_id = reservation_status.id").
		Joins("LEFT JOIN client ON reservation.client_id = client.id").
		Joins("LEFT JOIN table ON reservation.table_id = table.id").
		Joins("LEFT JOIN business ON reservation.business_id = business.id").
		Joins("LEFT JOIN user ON reservation.created_by_user_id = user.id").
		Where("reservation.id = ?", id).
		First(&tempResult).Error

	if err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al obtener reserva por ID")
		return nil, err
	}

	// Obtener el historial de estados
	var historyResults []entities.ReservationStatusHistory
	if err := r.database.Conn(ctx).
		Table("reservation_status_history").
		Select(`
			reservation_status_history.*,
			reservation_status.code as status_code,
			reservation_status.name as status_name,
			user.name as changed_by_user
		`).
		Joins("LEFT JOIN reservation_status ON reservation_status_history.status_id = reservation_status.id").
		Joins("LEFT JOIN user ON reservation_status_history.changed_by_user_id = user.id").
		Where("reservation_status_history.reservation_id = ?", id).
		Order("reservation_status_history.created_at ASC").
		Find(&historyResults).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al obtener historial de estados")
		tempResult.StatusHistory = []entities.ReservationStatusHistory{}
	} else {
		tempResult.StatusHistory = historyResults
	}

	return &tempResult, nil
}

// CancelReservation cancela una reserva
func (r *ReservationRepository) CancelReservation(ctx context.Context, id uint, reason string) (string, error) {
	// Verificar que la reserva existe
	var existingReservation entities.Reservation
	if err := r.database.Conn(ctx).Table("reservation").Where("id = ?", id).First(&existingReservation).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al verificar reserva existente")
		return "", err
	}

	// Actualizar el estado de la reserva a cancelado (asumiendo que el ID del estado cancelado es 3)
	if err := r.database.Conn(ctx).Table("reservation").Where("id = ?", id).Update("status_id", 3).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al cancelar reserva")
		return "", err
	}

	// Crear registro en el historial
	history := entities.ReservationStatusHistory{
		ReservationID: id,
		StatusID:      3, // Estado cancelado
	}

	if err := r.database.Conn(ctx).Table("reservation_status_history").Create(&history).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al crear historial de cancelación")
		// No retornamos error aquí porque la reserva ya fue cancelada
	}

	return fmt.Sprintf("Reserva cancelada con ID: %d", id), nil
}

// UpdateReservation actualiza una reserva
func (r *ReservationRepository) UpdateReservation(ctx context.Context, id uint, tableID *uint, startAt *time.Time, endAt *time.Time, numberOfGuests *int) (string, error) {
	updates := make(map[string]interface{})

	if tableID != nil {
		updates["table_id"] = *tableID
	}
	if startAt != nil {
		updates["start_at"] = *startAt
	}
	if endAt != nil {
		updates["end_at"] = *endAt
	}
	if numberOfGuests != nil {
		updates["number_of_guests"] = *numberOfGuests
	}

	if len(updates) == 0 {
		return "No hay campos para actualizar", nil
	}

	if err := r.database.Conn(ctx).Table("reservation").Where("id = ?", id).Updates(updates).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al actualizar reserva")
		return "", err
	}

	return fmt.Sprintf("Reserva actualizada con ID: %d", id), nil
}

// CreateReservationStatusHistory crea un registro en el historial de estados
func (r *ReservationRepository) CreateReservationStatusHistory(ctx context.Context, history entities.ReservationStatusHistory) error {
	dbHistory := dtos.DBReservationStatusHistory{
		ReservationID:   history.ReservationID,
		StatusID:        history.StatusID,
		ChangedByUserID: history.ChangedByUserID,
	}

	if err := r.database.Conn(ctx).Table("reservation_status_history").Create(&dbHistory).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al crear historial de estado")
		return err
	}

	return nil
}
