package usecasereserve

import (
	"central_reserve/internal/domain"
	"context"
)

func (u *ReserveUseCase) CancelReservation(ctx context.Context, id uint, reason string) (string, error) {
	// Primero obtener los detalles de la reserva para el email
	reservationDetail, err := u.repository.GetReserveByID(ctx, id)
	if err != nil {
		return "", err
	}

	response, err := u.repository.CancelReservation(ctx, id, reason)
	if err != nil {
		return "", err
	}

	if response == "" {
		return "", nil // Reserva no encontrada
	}

	// Enviar email de cancelación (en background para no bloquear la respuesta)
	go func() {
		// Convertir el DTO a la entidad Reservation para el email
		reservation := u.convertDTOToReservation(reservationDetail)
		if err := u.emailService.SendReservationCancellation(ctx, reservationDetail.ClienteEmail, reservationDetail.ClienteNombre, reservation); err != nil {
			u.log.Error(ctx).Err(err).Str("email", reservationDetail.ClienteEmail).Msg("Error enviando email de cancelación")
		} else {
			u.log.Info(ctx).Str("email", reservationDetail.ClienteEmail).Msg("Email de cancelación enviado exitosamente")
		}
	}()

	return response, nil
}

// Función auxiliar para convertir DTO a entidad
func (u *ReserveUseCase) convertDTOToReservation(dto *domain.ReserveDetailDTO) domain.Reservation {
	return domain.Reservation{
		Id:             dto.ReservaID,
		RestaurantID:   dto.RestauranteID,
		TableID:        dto.MesaID,
		ClientID:       dto.ClienteID,
		StartAt:        dto.StartAt,
		EndAt:          dto.EndAt,
		NumberOfGuests: dto.NumberOfGuests,
		StatusID:       0, // No es relevante para el email
		CreatedAt:      dto.ReservaCreada,
		UpdatedAt:      dto.ReservaActualizada,
	}
}
