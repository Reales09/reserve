package email

import (
	"central_reserve/internal/domain/entities"
	"central_reserve/internal/domain/ports"
	"context"
	"fmt"
)

// EmailService implementa ports.IEmailService
type EmailService struct {
	// Aquí se pueden agregar dependencias como configuraciones de SMTP, etc.
}

// NewEmailService crea una nueva instancia del servicio de email
func NewEmailService() ports.IEmailService {
	return &EmailService{}
}

// SendReservationConfirmation envía un email de confirmación de reserva
func (e *EmailService) SendReservationConfirmation(ctx context.Context, email, name string, reservation entities.Reservation) error {
	// Implementación del envío de email de confirmación
	fmt.Printf("Enviando email de confirmación a %s (%s) para la reserva %d\n", name, email, reservation.Id)
	return nil
}

// SendReservationCancellation envía un email de cancelación de reserva
func (e *EmailService) SendReservationCancellation(ctx context.Context, email, name string, reservation entities.Reservation) error {
	// Implementación del envío de email de cancelación
	fmt.Printf("Enviando email de cancelación a %s (%s) para la reserva %d\n", name, email, reservation.Id)
	return nil
}
