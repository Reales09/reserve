package repository

import (
	"central_reserve/internal/domain/entities"
	"central_reserve/internal/domain/ports"
	"central_reserve/internal/infra/secundary/repository/db"
	"central_reserve/internal/pkg/log"
	"context"
	"fmt"
)

// BusinessTypeRepository implementa ports.IBusinessTypeRepository
type BusinessTypeRepository struct {
	database db.IDatabase
	logger   log.ILogger
}

// NewBusinessTypeRepository crea una nueva instancia del repositorio de tipos de negocio
func NewBusinessTypeRepository(database db.IDatabase, logger log.ILogger) ports.IBusinessTypeRepository {
	return &BusinessTypeRepository{
		database: database,
		logger:   logger,
	}
}

// GetBusinessTypes obtiene todos los tipos de negocio
func (r *BusinessTypeRepository) GetBusinessTypes(ctx context.Context) ([]entities.BusinessType, error) {
	var businessTypes []entities.BusinessType
	if err := r.database.Conn(ctx).Table("business_type").Find(&businessTypes).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al obtener tipos de negocio")
		return nil, err
	}
	return businessTypes, nil
}

// GetBusinessTypeByID obtiene un tipo de negocio por su ID
func (r *BusinessTypeRepository) GetBusinessTypeByID(ctx context.Context, id uint) (*entities.BusinessType, error) {
	var businessType entities.BusinessType
	if err := r.database.Conn(ctx).Table("business_type").Where("id = ?", id).First(&businessType).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al obtener tipo de negocio por ID")
		return nil, err
	}
	return &businessType, nil
}

// GetBusinessTypeByCode obtiene un tipo de negocio por su código
func (r *BusinessTypeRepository) GetBusinessTypeByCode(ctx context.Context, code string) (*entities.BusinessType, error) {
	var businessType entities.BusinessType
	if err := r.database.Conn(ctx).Table("business_type").Where("code = ?", code).First(&businessType).Error; err != nil {
		r.logger.Error().Str("code", code).Err(err).Msg("Error al obtener tipo de negocio por código")
		return nil, err
	}
	return &businessType, nil
}

// CreateBusinessType crea un nuevo tipo de negocio
func (r *BusinessTypeRepository) CreateBusinessType(ctx context.Context, businessType entities.BusinessType) (string, error) {
	if err := r.database.Conn(ctx).Table("business_type").Create(&businessType).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al crear tipo de negocio")
		return "", err
	}
	return fmt.Sprintf("Tipo de negocio creado con ID: %d", businessType.ID), nil
}

// UpdateBusinessType actualiza un tipo de negocio existente
func (r *BusinessTypeRepository) UpdateBusinessType(ctx context.Context, id uint, businessType entities.BusinessType) (string, error) {
	if err := r.database.Conn(ctx).Table("business_type").Where("id = ?", id).Updates(&businessType).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al actualizar tipo de negocio")
		return "", err
	}
	return fmt.Sprintf("Tipo de negocio actualizado con ID: %d", id), nil
}

// DeleteBusinessType elimina un tipo de negocio
func (r *BusinessTypeRepository) DeleteBusinessType(ctx context.Context, id uint) (string, error) {
	if err := r.database.Conn(ctx).Table("business_type").Where("id = ?", id).Delete(&entities.BusinessType{}).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al eliminar tipo de negocio")
		return "", err
	}
	return fmt.Sprintf("Tipo de negocio eliminado con ID: %d", id), nil
}
