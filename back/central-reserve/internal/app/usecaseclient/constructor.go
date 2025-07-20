package usecaseclient

import (
	"central_reserve/internal/domain/entities"
	"central_reserve/internal/domain/ports"
	"context"
)

type IUseCaseClient interface {
	GetClients(ctx context.Context) ([]entities.Client, error)
	GetClientByID(ctx context.Context, id uint) (*entities.Client, error)
	CreateClient(ctx context.Context, client entities.Client) (string, error)
	UpdateClient(ctx context.Context, id uint, client entities.Client) (string, error)
	DeleteClient(ctx context.Context, id uint) (string, error)
}

type ClientUseCase struct {
	repository ports.IClientUseCaseRepository
}

func NewClientUseCase(repository ports.IClientUseCaseRepository) *ClientUseCase {
	return &ClientUseCase{
		repository: repository,
	}
}
