package usecasebusiness

import (
	"central_reserve/internal/domain/dtos"
	"context"
	"fmt"
)

// GetBusinesses obtiene todos los negocios
func (uc *BusinessUseCase) GetBusinesses(ctx context.Context) ([]dtos.BusinessResponse, error) {
	uc.log.Info().Msg("Obteniendo negocios")

	businesses, err := uc.repository.GetBusinesses(ctx)
	if err != nil {
		uc.log.Error().Err(err).Msg("Error al obtener negocios")
		return nil, fmt.Errorf("error al obtener negocios: %w", err)
	}

	// Convertir entidades a DTOs
	response := make([]dtos.BusinessResponse, len(businesses))
	for i, business := range businesses {
		response[i] = dtos.BusinessResponse{
			ID:   business.ID,
			Name: business.Name,
			Code: business.Code,
			BusinessType: dtos.BusinessTypeResponse{
				ID: business.BusinessTypeID, // Nota: necesitaríamos obtener el BusinessType completo
			},
			Timezone:           business.Timezone,
			Address:            business.Address,
			Description:        business.Description,
			LogoURL:            business.LogoURL,
			PrimaryColor:       business.PrimaryColor,
			SecondaryColor:     business.SecondaryColor,
			CustomDomain:       business.CustomDomain,
			IsActive:           business.IsActive,
			EnableDelivery:     business.EnableDelivery,
			EnablePickup:       business.EnablePickup,
			EnableReservations: business.EnableReservations,
			CreatedAt:          business.CreatedAt,
			UpdatedAt:          business.UpdatedAt,
		}
	}

	uc.log.Info().Int("count", len(response)).Msg("Negocios obtenidos exitosamente")
	return response, nil
}
