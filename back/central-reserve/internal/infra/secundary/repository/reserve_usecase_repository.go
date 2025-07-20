package repository

import (
	"central_reserve/internal/domain/dtos"
	"central_reserve/internal/domain/entities"
	"central_reserve/internal/domain/ports"
	"central_reserve/internal/infra/secundary/repository/db"
	"central_reserve/internal/pkg/log"
	"context"
	"time"
)

// ReserveUseCaseRepository implementa ports.IReserveUseCaseRepository
type ReserveUseCaseRepository struct {
	database db.IDatabase
	logger   log.ILogger
}

// NewReserveUseCaseRepository crea una nueva instancia del repositorio de casos de uso de reservas
func NewReserveUseCaseRepository(database db.IDatabase, logger log.ILogger) ports.IReserveUseCaseRepository {
	return &ReserveUseCaseRepository{
		database: database,
		logger:   logger,
	}
}

// GetClientByEmailAndBusiness obtiene un cliente por email y business_id
func (r *ReserveUseCaseRepository) GetClientByEmailAndBusiness(ctx context.Context, email string, businessID uint) (*entities.Client, error) {
	var client entities.Client
	if err := r.database.Conn(ctx).Table("client").Where("email = ? AND business_id = ?", email, businessID).First(&client).Error; err != nil {
		r.logger.Error().Str("email", email).Uint("business_id", businessID).Msg("Error al obtener cliente por email y business")
		return nil, err
	}
	return &client, nil
}

// CreateClient crea un nuevo cliente
func (r *ReserveUseCaseRepository) CreateClient(ctx context.Context, client entities.Client) (string, error) {
	if err := r.database.Conn(ctx).Table("client").Create(&client).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al crear cliente")
		return "", err
	}
	return "Cliente creado exitosamente", nil
}

// CreateReserve crea una nueva reserva
func (r *ReserveUseCaseRepository) CreateReserve(ctx context.Context, reserve entities.Reservation) (string, error) {
	if err := r.database.Conn(ctx).Table("reservation").Create(&reserve).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al crear reserva")
		return "", err
	}
	return "Reserva creada exitosamente", nil
}

// GetLatestReservationByClient obtiene la última reserva de un cliente
func (r *ReserveUseCaseRepository) GetLatestReservationByClient(ctx context.Context, clientID uint) (*entities.Reservation, error) {
	var reservation entities.Reservation
	if err := r.database.Conn(ctx).Table("reservation").Where("client_id = ?", clientID).Order("created_at DESC").First(&reservation).Error; err != nil {
		r.logger.Error().Uint("client_id", clientID).Msg("Error al obtener última reserva del cliente")
		return nil, err
	}
	return &reservation, nil
}

// GetReserves obtiene las reservas con filtros opcionales
func (r *ReserveUseCaseRepository) GetReserves(ctx context.Context, statusID *uint, clientID *uint, tableID *uint, startDate *time.Time, endDate *time.Time) ([]dtos.ReserveDetailDTO, error) {
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

	return tempResults, nil
}

// GetReserveByID obtiene una reserva por su ID
func (r *ReserveUseCaseRepository) GetReserveByID(ctx context.Context, id uint) (*dtos.ReserveDetailDTO, error) {
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

	// Inicializar el slice de historial
	tempResult.StatusHistory = []entities.ReservationStatusHistory{}

	return &tempResult, nil
}

// CancelReservation cancela una reserva
func (r *ReserveUseCaseRepository) CancelReservation(ctx context.Context, id uint, reason string) (string, error) {
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

	return "Reserva cancelada exitosamente", nil
}

// UpdateReservation actualiza una reserva
func (r *ReserveUseCaseRepository) UpdateReservation(ctx context.Context, id uint, tableID *uint, startAt *time.Time, endAt *time.Time, numberOfGuests *int) (string, error) {
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

	return "Reserva actualizada exitosamente", nil
}

// CreateReservationStatusHistory crea un registro en el historial de estados
func (r *ReserveUseCaseRepository) CreateReservationStatusHistory(ctx context.Context, history entities.ReservationStatusHistory) error {
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
