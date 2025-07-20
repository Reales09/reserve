package usecasereserve

import (
	"central_reserve/internal/domain/dtos"
	"context"
	"fmt"
)

// GetReserveByID obtiene una reserva por su ID
func (u *ReserveUseCase) GetReserveByID(ctx context.Context, id uint) (*dtos.ReserveDetailDTO, error) {
	reservation, err := u.repository.GetReserveByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error al obtener reserva: %w", err)
	}
	return reservation, nil
}
