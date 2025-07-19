package usecaseauth

import (
	"central_reserve/internal/domain"
	"central_reserve/internal/pkg/jwt"
	"central_reserve/internal/pkg/log"
	"context"
)

type IUseCaseAuth interface {
	Login(ctx context.Context, request domain.LoginRequest) (*domain.LoginResponse, error)
}

type AuthUseCase struct {
	repository domain.IHolaMundo
	jwtService *jwt.JWTService
	log        log.ILogger
}

func NewAuthUseCase(repository domain.IHolaMundo, jwtService *jwt.JWTService, log log.ILogger) IUseCaseAuth {
	return &AuthUseCase{
		repository: repository,
		jwtService: jwtService,
		log:        log,
	}
}
