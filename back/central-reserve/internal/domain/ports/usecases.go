package ports

import (
	"central_reserve/internal/domain/dtos"
	"central_reserve/internal/domain/entities"
	"context"
	"time"
)

// IClientUseCaseRepository define los métodos que necesita ClientUseCase
type IClientUseCaseRepository interface {
	GetClients(ctx context.Context) ([]entities.Client, error)
	GetClientByID(ctx context.Context, id uint) (*entities.Client, error)
	CreateClient(ctx context.Context, client entities.Client) (string, error)
	UpdateClient(ctx context.Context, id uint, client entities.Client) (string, error)
	DeleteClient(ctx context.Context, id uint) (string, error)
	GetClientByEmail(ctx context.Context, email string) (*entities.Client, error)
}

// ITableUseCaseRepository define los métodos que necesita TableUseCase
type ITableUseCaseRepository interface {
	CreateTable(ctx context.Context, table entities.Table) (string, error)
	GetTables(ctx context.Context) ([]entities.Table, error)
	GetTableByID(ctx context.Context, id uint) (*entities.Table, error)
	UpdateTable(ctx context.Context, id uint, table entities.Table) (string, error)
	DeleteTable(ctx context.Context, id uint) (string, error)
}

// IAuthUseCaseRepository define los métodos que necesita AuthUseCase
type IAuthUseCaseRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserRoles(ctx context.Context, userID uint) ([]entities.Role, error)
	GetRolePermissions(ctx context.Context, roleID uint) ([]entities.Permission, error)
	UpdateLastLogin(ctx context.Context, userID uint) error
}

// IReserveUseCaseRepository define los métodos que necesita ReserveUseCase
type IReserveUseCaseRepository interface {
	// Métodos de clientes que necesita
	GetClientByEmailAndBusiness(ctx context.Context, email string, businessID uint) (*entities.Client, error)
	CreateClient(ctx context.Context, client entities.Client) (string, error)

	// Métodos de reservas que necesita
	CreateReserve(ctx context.Context, reserve entities.Reservation) (string, error)
	GetLatestReservationByClient(ctx context.Context, clientID uint) (*entities.Reservation, error)
	GetReserves(ctx context.Context, statusID *uint, clientID *uint, tableID *uint, startDate *time.Time, endDate *time.Time) ([]dtos.ReserveDetailDTO, error)
	GetReserveByID(ctx context.Context, id uint) (*dtos.ReserveDetailDTO, error)
	CancelReservation(ctx context.Context, id uint, reason string) (string, error)
	UpdateReservation(ctx context.Context, id uint, tableID *uint, startAt *time.Time, endAt *time.Time, numberOfGuests *int) (string, error)
	CreateReservationStatusHistory(ctx context.Context, history entities.ReservationStatusHistory) error
}

// IBusinessTypeUseCaseRepository define los métodos que necesita BusinessTypeUseCase
type IBusinessTypeUseCaseRepository interface {
	GetBusinessTypes(ctx context.Context) ([]entities.BusinessType, error)
	GetBusinessTypeByID(ctx context.Context, id uint) (*entities.BusinessType, error)
	GetBusinessTypeByCode(ctx context.Context, code string) (*entities.BusinessType, error)
	CreateBusinessType(ctx context.Context, businessType entities.BusinessType) (string, error)
	UpdateBusinessType(ctx context.Context, id uint, businessType entities.BusinessType) (string, error)
	DeleteBusinessType(ctx context.Context, id uint) (string, error)
}

// IBusinessUseCaseRepository define los métodos que necesita BusinessUseCase
type IBusinessUseCaseRepository interface {
	GetBusinesses(ctx context.Context) ([]entities.Business, error)
	GetBusinessByID(ctx context.Context, id uint) (*entities.Business, error)
	GetBusinessByCode(ctx context.Context, code string) (*entities.Business, error)
	GetBusinessByCustomDomain(ctx context.Context, domain string) (*entities.Business, error)
	CreateBusiness(ctx context.Context, business entities.Business) (string, error)
	UpdateBusiness(ctx context.Context, id uint, business entities.Business) (string, error)
	DeleteBusiness(ctx context.Context, id uint) (string, error)
}
