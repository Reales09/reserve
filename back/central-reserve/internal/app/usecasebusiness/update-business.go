package usecasebusiness

import (
	"central_reserve/internal/domain/dtos"
	"central_reserve/internal/domain/entities"
	"context"
	"fmt"
)

// UpdateBusiness actualiza un negocio existente
func (uc *BusinessUseCase) UpdateBusiness(ctx context.Context, id uint, request dtos.BusinessRequest) (*dtos.BusinessResponse, error) {
	uc.log.Info().Uint("id", id).Str("name", request.Name).Msg("Actualizando negocio")

	// Verificar que existe
	existing, err := uc.repository.GetBusinessByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al obtener negocio para actualizar")
		return nil, fmt.Errorf("error al obtener negocio: %w", err)
	}

	if existing == nil {
		uc.log.Warn().Uint("id", id).Msg("Negocio no encontrado para actualizar")
		return nil, fmt.Errorf("negocio no encontrado")
	}

	// Verificar que el código no exista en otro negocio
	if request.Code != existing.Code {
		codeExists, err := uc.repository.GetBusinessByCode(ctx, request.Code)
		if err != nil && err.Error() != "negocio no encontrado" {
			uc.log.Error().Err(err).Str("code", request.Code).Msg("Error al verificar código existente")
			return nil, fmt.Errorf("error al verificar código existente: %w", err)
		}

		if codeExists != nil {
			uc.log.Warn().Str("code", request.Code).Msg("Código de negocio ya existe")
			return nil, fmt.Errorf("el código '%s' ya existe", request.Code)
		}
	}

	// Verificar que el dominio personalizado no exista en otro negocio
	if request.CustomDomain != "" && request.CustomDomain != existing.CustomDomain {
		domainExists, err := uc.repository.GetBusinessByCustomDomain(ctx, request.CustomDomain)
		if err != nil && err.Error() != "negocio no encontrado" {
			uc.log.Error().Err(err).Str("domain", request.CustomDomain).Msg("Error al verificar dominio existente")
			return nil, fmt.Errorf("error al verificar dominio existente: %w", err)
		}

		if domainExists != nil {
			uc.log.Warn().Str("domain", request.CustomDomain).Msg("Dominio personalizado ya existe")
			return nil, fmt.Errorf("el dominio '%s' ya existe", request.CustomDomain)
		}
	}

	// Actualizar entidad
	business := entities.Business{
		Name:               request.Name,
		Code:               request.Code,
		BusinessTypeID:     request.BusinessTypeID,
		Timezone:           request.Timezone,
		Address:            request.Address,
		Description:        request.Description,
		LogoURL:            request.LogoURL,
		PrimaryColor:       request.PrimaryColor,
		SecondaryColor:     request.SecondaryColor,
		CustomDomain:       request.CustomDomain,
		IsActive:           request.IsActive,
		EnableDelivery:     request.EnableDelivery,
		EnablePickup:       request.EnablePickup,
		EnableReservations: request.EnableReservations,
	}

	// Guardar en repositorio
	_, err = uc.repository.UpdateBusiness(ctx, id, business)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al actualizar negocio")
		return nil, fmt.Errorf("error al actualizar negocio: %w", err)
	}

	// Obtener el negocio actualizado
	updated, err := uc.repository.GetBusinessByID(ctx, id)
	if err != nil {
		uc.log.Error().Err(err).Uint("id", id).Msg("Error al obtener negocio actualizado")
		return nil, fmt.Errorf("error al obtener negocio actualizado: %w", err)
	}

	response := &dtos.BusinessResponse{
		ID:   updated.ID,
		Name: updated.Name,
		Code: updated.Code,
		BusinessType: dtos.BusinessTypeResponse{
			ID: updated.BusinessTypeID, // Nota: necesitaríamos obtener el BusinessType completo
		},
		Timezone:           updated.Timezone,
		Address:            updated.Address,
		Description:        updated.Description,
		LogoURL:            updated.LogoURL,
		PrimaryColor:       updated.PrimaryColor,
		SecondaryColor:     updated.SecondaryColor,
		CustomDomain:       updated.CustomDomain,
		IsActive:           updated.IsActive,
		EnableDelivery:     updated.EnableDelivery,
		EnablePickup:       updated.EnablePickup,
		EnableReservations: updated.EnableReservations,
		CreatedAt:          updated.CreatedAt,
		UpdatedAt:          updated.UpdatedAt,
	}

	uc.log.Info().Uint("id", id).Str("name", request.Name).Msg("Negocio actualizado exitosamente")
	return response, nil
}
