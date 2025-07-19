package domain

import (
	"time"
)

// User representa un usuario del sistema
type User struct {
	ID          uint
	Name        string
	Email       string
	Password    string
	Phone       string
	AvatarURL   string
	IsActive    bool
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Role representa un rol del sistema
type Role struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Level       int
	IsSystem    bool
	Scope       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Permission representa un permiso del sistema
type Permission struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Resource    string
	Action      string
	Scope       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// UserRole representa la relación many-to-many entre usuarios y roles
type UserRole struct {
	UserID    uint
	RoleID    uint
	CreatedAt time.Time
}

// RolePermission representa la relación many-to-many entre roles y permisos
type RolePermission struct {
	RoleID       uint
	PermissionID uint
	CreatedAt    time.Time
}

// LoginRequest representa la solicitud de login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse representa la respuesta del login
type LoginResponse struct {
	User        UserInfo         `json:"user"`
	Roles       []RoleInfo       `json:"roles"`
	Permissions []PermissionInfo `json:"permissions"`
	Token       string           `json:"token"`
	IsSuper     bool             `json:"is_super"`
}

// UserInfo información básica del usuario para la respuesta
type UserInfo struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	AvatarURL   string     `json:"avatar_url"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
}

// RoleInfo información básica del rol para la respuesta
type RoleInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	Scope       string `json:"scope"`
}

// PermissionInfo información básica del permiso para la respuesta
type PermissionInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Scope       string `json:"scope"`
}

type Reservation struct {
	Id              uint
	RestaurantID    uint
	TableID         *uint
	ClientID        uint
	CreatedByUserID *uint
	StartAt         time.Time
	EndAt           time.Time
	NumberOfGuests  int
	StatusID        uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type Client struct {
	ID           uint
	RestaurantID uint
	Name         string
	Email        string
	Phone        string
	Dni          *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type Table struct {
	ID           uint
	RestaurantID uint
	Number       int
	Capacity     int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type RestaurantStaff struct {
	ID           uint
	UserID       uint
	RestaurantID uint
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type ReservationStatus struct {
	ID        uint
	Code      string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type ReservationStatusHistory struct {
	ID              uint
	ReservationID   uint
	StatusID        uint
	ChangedByUserID *uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
	StatusCode      string
	StatusName      string
	ChangedByUser   *string
}

type ReserveDetailDTO struct {
	// Reserva
	ReservaID          uint
	StartAt            time.Time
	EndAt              time.Time
	NumberOfGuests     int
	ReservaCreada      time.Time
	ReservaActualizada time.Time

	// Estado
	EstadoCodigo string
	EstadoNombre string

	// Cliente
	ClienteID       uint
	ClienteNombre   string
	ClienteEmail    string
	ClienteTelefono string
	ClienteDni      *string

	// Mesa
	MesaID        *uint
	MesaNumero    *int
	MesaCapacidad *int

	// Restaurante
	RestauranteID        uint
	RestauranteNombre    string
	RestauranteCodigo    string
	RestauranteDireccion string

	// Usuario
	UsuarioID     *uint
	UsuarioNombre *string
	UsuarioEmail  *string

	// Historial de Estados
	StatusHistory []ReservationStatusHistory
}
