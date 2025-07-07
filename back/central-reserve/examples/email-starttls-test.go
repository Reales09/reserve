package main

import (
	"central_reserve/internal/domain"
	"central_reserve/internal/infra/secundary/email"
	"central_reserve/internal/pkg/env"
	"central_reserve/internal/pkg/log"
	"context"
	"time"
)

func main() {
	// Inicializar logger
	logger := log.New()

	// Cargar configuración (asegúrate de tener un .env con las variables SMTP)
	config, err := env.New(logger)
	if err != nil {
		logger.Error(context.Background()).Err(err).Msg("Error cargando configuración")
		return
	}

	// Crear servicio de email
	emailService := email.New(config, logger)

	// Crear una reserva de ejemplo
	reservation := domain.Reservation{
		Id:             1,
		RestaurantID:   1,
		TableID:        nil,
		ClientID:       1,
		StartAt:        time.Now().Add(24 * time.Hour), // Mañana
		EndAt:          time.Now().Add(26 * time.Hour), // 2 horas después
		NumberOfGuests: 4,
		StatusID:       1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	ctx := context.Background()

	// Mostrar configuración actual
	logger.Info(ctx).
		Str("smtp_host", config.Get("SMTP_HOST")).
		Str("smtp_port", config.Get("SMTP_PORT")).
		Str("use_starttls", config.Get("SMTP_USE_STARTTLS")).
		Str("use_tls", config.Get("SMTP_USE_TLS")).
		Msg("Configuración SMTP actual")

	// Enviar email de confirmación con STARTTLS
	logger.Info(ctx).Msg("Enviando email de confirmación...")
	err = emailService.SendReservationConfirmation(ctx, "cliente@ejemplo.com", "Juan Pérez", reservation)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Error enviando email de confirmación")
	} else {
		logger.Info(ctx).Msg("✅ Email de confirmación enviado exitosamente con STARTTLS")
	}

	// Enviar email de cancelación con STARTTLS
	logger.Info(ctx).Msg("Enviando email de cancelación...")
	err = emailService.SendReservationCancellation(ctx, "cliente@ejemplo.com", "Juan Pérez", reservation)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Error enviando email de cancelación")
	} else {
		logger.Info(ctx).Msg("✅ Email de cancelación enviado exitosamente con STARTTLS")
	}

	logger.Info(ctx).Msg("🎉 Prueba de STARTTLS completada")
}
