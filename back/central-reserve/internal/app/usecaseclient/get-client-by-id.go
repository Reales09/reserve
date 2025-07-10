package usecaseclient

import (
	"central_reserve/internal/domain"
	"context"
)

func (u *ClientUseCase) GetClientByID(ctx context.Context, id uint) (*domain.Client, error) {
	response, err := u.repository.GetClientByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (u *ClientUseCase) GetClientByEmail(ctx context.Context, email string) (*domain.Client, error) {
	return u.repository.GetClientByEmail(ctx, email)
}
