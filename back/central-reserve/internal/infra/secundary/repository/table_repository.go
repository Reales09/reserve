package repository

import (
	"central_reserve/internal/domain/entities"
	"central_reserve/internal/domain/ports"
	"central_reserve/internal/infra/secundary/repository/db"
	"central_reserve/internal/pkg/log"
	"context"
	"fmt"
)

// TableRepository implementa ports.ITableRepository
type TableRepository struct {
	database db.IDatabase
	logger   log.ILogger
}

// NewTableRepository crea una nueva instancia del repositorio de mesas
func NewTableRepository(database db.IDatabase, logger log.ILogger) ports.ITableRepository {
	return &TableRepository{
		database: database,
		logger:   logger,
	}
}

// CreateTable crea una nueva mesa
func (r *TableRepository) CreateTable(ctx context.Context, table entities.Table) (string, error) {
	if err := r.database.Conn(ctx).Table("table").Create(&table).Error; err != nil {
		r.logger.Error().Err(err).Msg("Error al crear mesa")
		return "", err
	}
	return fmt.Sprintf("Mesa creada con ID: %d", table.ID), nil
}

// GetTables obtiene todas las mesas
func (r *TableRepository) GetTables(ctx context.Context) ([]entities.Table, error) {
	var tables []entities.Table
	if err := r.database.Conn(ctx).Table("table").Find(&tables).Error; err != nil {
		r.logger.Error().Msg("Error al obtener mesas")
		return nil, err
	}
	return tables, nil
}

// GetTableByID obtiene una mesa por su ID
func (r *TableRepository) GetTableByID(ctx context.Context, id uint) (*entities.Table, error) {
	var table entities.Table
	if err := r.database.Conn(ctx).Table("table").Where("id = ?", id).First(&table).Error; err != nil {
		r.logger.Error().Uint("id", id).Msg("Error al obtener mesa por ID")
		return nil, err
	}
	return &table, nil
}

// UpdateTable actualiza una mesa existente
func (r *TableRepository) UpdateTable(ctx context.Context, id uint, table entities.Table) (string, error) {
	if err := r.database.Conn(ctx).Table("table").Where("id = ?", id).Updates(&table).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al actualizar mesa")
		return "", err
	}
	return fmt.Sprintf("Mesa actualizada con ID: %d", id), nil
}

// DeleteTable elimina una mesa
func (r *TableRepository) DeleteTable(ctx context.Context, id uint) (string, error) {
	if err := r.database.Conn(ctx).Table("table").Where("id = ?", id).Delete(&entities.Table{}).Error; err != nil {
		r.logger.Error().Uint("id", id).Err(err).Msg("Error al eliminar mesa")
		return "", err
	}
	return fmt.Sprintf("Mesa eliminada con ID: %d", id), nil
}
