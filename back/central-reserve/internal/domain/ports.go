package domain

import (
	"context"
	"time"
)

type IHolaMundo interface {
	HolaMundo() string
	CreateReserve(ctx context.Context, reserve Reservation) (string, error)
	GetClients(ctx context.Context) ([]Client, error)
	GetClientByID(ctx context.Context, id uint) (*Client, error)
	CreateClient(ctx context.Context, client Client) (string, error)
	UpdateClient(ctx context.Context, id uint, client Client) (string, error)
	DeleteClient(ctx context.Context, id uint) (string, error)
	CreateTable(ctx context.Context, table Table) (string, error)
	GetTables(ctx context.Context) ([]Table, error)
	GetTableByID(ctx context.Context, id uint) (*Table, error)
	UpdateTable(ctx context.Context, id uint, table Table) (string, error)
	DeleteTable(ctx context.Context, id uint) (string, error)
	GetClientByEmail(ctx context.Context, email string) (*Client, error)
	GetClientByEmailAndRestaurant(ctx context.Context, email string, restaurantID uint) (*Client, error)
	CreateReservationStatusHistory(ctx context.Context, history ReservationStatusHistory) error
	GetLatestReservationByClient(ctx context.Context, clientID uint) (*Reservation, error)
	GetReserves(ctx context.Context, statusID *uint, clientID *uint, tableID *uint, startDate *time.Time, endDate *time.Time) ([]ReserveDetailDTO, error)
	GetReserveByID(ctx context.Context, id uint) (*ReserveDetailDTO, error)
	CancelReservation(ctx context.Context, id uint, reason string) (string, error)
	UpdateReservation(ctx context.Context, id uint, tableID *uint, startAt *time.Time, endAt *time.Time, numberOfGuests *int) (string, error)
	GetReservationStatuses(ctx context.Context) ([]ReservationStatus, error)
}

type IEmailService interface {
	SendReservationConfirmation(ctx context.Context, email, name string, reservation Reservation) error
	SendReservationCancellation(ctx context.Context, email, name string, reservation Reservation) error
}
