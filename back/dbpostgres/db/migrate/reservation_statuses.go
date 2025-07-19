package migrate

import (
	"dbpostgres/db/models"
	"dbpostgres/db/repository"
	"dbpostgres/pkg/log"
)

// InitializeReservationStatuses inicializa los estados de reserva del sistema
func InitializeReservationStatuses(systemRepo repository.SystemRepository, logger log.ILogger) error {
	// Verificar si ya existen estados de reserva en la base de datos
	existingStatuses, err := systemRepo.GetAllReservationStatuses()
	if err != nil {
		return err
	}

	// Si ya existen estados, saltar la migración
	if len(existingStatuses) > 0 {
		logger.Info().Int("statuses_count", len(existingStatuses)).Msg("✅ Estados de reserva ya existen, saltando migración de estados de reserva")
		return nil
	}

	// Definir los estados a crear
	statuses := []models.ReservationStatus{
		{Code: "pending", Name: "Pendiente"},
		{Code: "confirmed", Name: "Confirmada"},
		{Code: "seated", Name: "Sentado"},
		{Code: "completed", Name: "Completada"},
		{Code: "cancelled", Name: "Cancelada"},
		{Code: "no_show", Name: "No se presentó"},
	}

	// Inicializar estados usando el repositorio
	if err := systemRepo.InitializeReservationStatuses(statuses); err != nil {
		return err
	}

	logger.Info().Int("statuses_count", len(statuses)).Msg("✅ Estados de reserva inicializados correctamente")
	return nil
}
