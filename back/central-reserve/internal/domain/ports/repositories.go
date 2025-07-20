package ports

import (
	"central_reserve/internal/domain/dtos"
	"central_reserve/internal/domain/entities"
	"context"
	"time"
)

// IClientRepository define las operaciones para clientes
type IClientRepository interface {
	GetClients(ctx context.Context) ([]entities.Client, error)
	GetClientByID(ctx context.Context, id uint) (*entities.Client, error)
	GetClientByEmail(ctx context.Context, email string) (*entities.Client, error)
	GetClientByEmailAndBusiness(ctx context.Context, email string, businessID uint) (*entities.Client, error)
	CreateClient(ctx context.Context, client entities.Client) (string, error)
	UpdateClient(ctx context.Context, id uint, client entities.Client) (string, error)
	DeleteClient(ctx context.Context, id uint) (string, error)
}

// ITableRepository define las operaciones para mesas
type ITableRepository interface {
	CreateTable(ctx context.Context, table entities.Table) (string, error)
	GetTables(ctx context.Context) ([]entities.Table, error)
	GetTableByID(ctx context.Context, id uint) (*entities.Table, error)
	UpdateTable(ctx context.Context, id uint, table entities.Table) (string, error)
	DeleteTable(ctx context.Context, id uint) (string, error)
}

// IReservationRepository define las operaciones para reservas
type IReservationRepository interface {
	CreateReserve(ctx context.Context, reserve entities.Reservation) (string, error)
	GetLatestReservationByClient(ctx context.Context, clientID uint) (*entities.Reservation, error)
	GetReserves(ctx context.Context, statusID *uint, clientID *uint, tableID *uint, startDate *time.Time, endDate *time.Time) ([]dtos.ReserveDetailDTO, error)
	GetReserveByID(ctx context.Context, id uint) (*dtos.ReserveDetailDTO, error)
	CancelReservation(ctx context.Context, id uint, reason string) (string, error)
	UpdateReservation(ctx context.Context, id uint, tableID *uint, startAt *time.Time, endAt *time.Time, numberOfGuests *int) (string, error)
	CreateReservationStatusHistory(ctx context.Context, history entities.ReservationStatusHistory) error
}

// IAuthRepository define las operaciones de autenticación
type IAuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserRoles(ctx context.Context, userID uint) ([]entities.Role, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]entities.Permission, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
}

// IAuthService define las operaciones de autenticación (métodos de repositorio)
type IAuthService interface {
	Login(ctx context.Context, email, password string) (*dtos.LoginResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserRoles(ctx context.Context, userID uint) ([]entities.Role, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]entities.Permission, error)
	ValidatePassword(hashedPassword, password string) error
	GenerateToken(userID uint, email string, roles []string) (string, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
}
