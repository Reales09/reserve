package ports

import (
	"central_reserve/internal/domain/entities"
	"context"
)

// IEmailService define las operaciones de envío de emails
type IEmailService interface {
	SendReservationConfirmation(ctx context.Context, email, name string, reservation entities.Reservation) error
	SendReservationCancellation(ctx context.Context, email, name string, reservation entities.Reservation) error
}
