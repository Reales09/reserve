package models

import (
	"time"

	"gorm.io/gorm"
)

// ───────────────────────────────────────────
//
//	RESTAURANTS  (multi-tenant) - MARCA BLANCA
//
// ───────────────────────────────────────────
type Restaurant struct {
	gorm.Model
	Name        string `gorm:"size:120;not null"`
	Code        string `gorm:"size:50;not null;unique"` // slug para URL personalizada
	Timezone    string `gorm:"size:40;default:'America/Bogota'"`
	Address     string `gorm:"size:255"`
	Description string `gorm:"size:500"`

	// Configuración de marca blanca
	LogoURL        string `gorm:"size:255"`
	PrimaryColor   string `gorm:"size:7;default:'#1f2937'"` // Hex color
	SecondaryColor string `gorm:"size:7;default:'#3b82f6'"` // Hex color
	CustomDomain   string `gorm:"size:100;unique"`          // dominio personalizado
	IsActive       bool   `gorm:"default:true"`

	// Configuración de funcionalidades
	EnableDelivery     bool `gorm:"default:false"`
	EnablePickup       bool `gorm:"default:false"`
	EnableReservations bool `gorm:"default:true"`

	Tables       []Table
	Reservations []Reservation
	Staff        []RestaurantStaff
	Clients      []Client
	Users        []User `gorm:"many2many:user_restaurants;"` // Usuarios del restaurante (muchos a muchos)
}

// ───────────────────────────────────────────
//
//	ROLES DEL SISTEMA
//
// ───────────────────────────────────────────
type Role struct {
	gorm.Model
	Name        string `gorm:"size:50;not null;unique"`
	Code        string `gorm:"size:30;not null;unique"` // Código interno
	Description string `gorm:"size:255"`
	Level       int    `gorm:"not null;default:1"` // Nivel jerárquico (1=super, 2=admin, 3=manager, 4=staff)
	IsSystem    bool   `gorm:"default:false"`      // Si es rol del sistema (no se puede eliminar)

	// Scope del rol
	Scope string `gorm:"size:20;not null;default:'restaurant'"` // 'platform', 'restaurant', 'both'

	Permissions []Permission `gorm:"many2many:role_permissions;"`
	Users       []User       `gorm:"many2many:user_roles;"`
}

// ───────────────────────────────────────────
//
//	PERMISOS DEL SISTEMA
//
// ───────────────────────────────────────────
type Permission struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique"`
	Code        string `gorm:"size:50;not null;unique"` // Código interno
	Description string `gorm:"size:255"`
	Resource    string `gorm:"size:50;not null"`                      // Recurso: 'restaurants', 'users', 'reservations', etc.
	Action      string `gorm:"size:20;not null"`                      // Acción: 'create', 'read', 'update', 'delete', 'manage'
	Scope       string `gorm:"size:20;not null;default:'restaurant'"` // 'platform', 'restaurant', 'both'

	Roles []Role `gorm:"many2many:role_permissions;"`
}

// ───────────────────────────────────────────
//
//	USUARIOS DEL SISTEMA
//
// ───────────────────────────────────────────
type User struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"`
	Email       string `gorm:"size:255;not null;unique"`
	Password    string `gorm:"size:255;not null"`
	Phone       string `gorm:"size:20"`
	AvatarURL   string `gorm:"size:255"`
	IsActive    bool   `gorm:"default:true"`
	LastLoginAt *time.Time

	// Relación con restaurantes (un usuario puede estar en múltiples restaurantes)
	Restaurants []Restaurant `gorm:"many2many:user_restaurants;"`

	// Roles del usuario
	Roles []Role `gorm:"many2many:user_roles;"`

	// Relaciones existentes
	StaffOf []RestaurantStaff
}

// ───────────────────────────────────────────
//
//	RESTAURANT STAFF  (N:M usuario ↔ restaurante)
//
// ───────────────────────────────────────────
type RestaurantStaff struct {
	gorm.Model
	UserID       uint `gorm:"not null;index;uniqueIndex:idx_user_restaurant,priority:1"`
	RestaurantID uint `gorm:"not null;index;uniqueIndex:idx_user_restaurant,priority:2"`
	RoleID       uint `gorm:"not null;index"` // Ahora referencia a Role en lugar de string

	User       User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Restaurant Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role       Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// ───────────────────────────────────────────
//
//	CLIENTS – personas que hacen la reserva
//
// ───────────────────────────────────────────
type Client struct {
	gorm.Model
	RestaurantID uint    `gorm:"not null;index;uniqueIndex:idx_rest_client_email,priority:1"`
	Name         string  `gorm:"size:255;not null"`
	Email        string  `gorm:"size:255;uniqueIndex:idx_rest_client_email,priority:2"`
	Phone        string  `gorm:"size:20"`
	Dni          *string `gorm:"size:30;uniqueIndex:idx_rest_client_dni,priority:2"`

	Reservations []Reservation
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// ───────────────────────────────────────────
//
//	TABLES – mesas físicas
//
// ───────────────────────────────────────────
type Table struct {
	gorm.Model
	RestaurantID uint `gorm:"not null;index;uniqueIndex:idx_rest_table_number,priority:1"`
	Number       int  `gorm:"not null;uniqueIndex:idx_rest_table_number,priority:2"`
	Capacity     int  `gorm:"not null"`

	Reservations []Reservation
}

// ───────────────────────────────────────────
//
//	RESERVATIONS
//
// ───────────────────────────────────────────
type Reservation struct {
	gorm.Model
	RestaurantID uint  `gorm:"not null;index"`
	TableID      *uint `gorm:"index"`
	ClientID     uint  `gorm:"not null;index"`
	// Opcional: quién registró la reserva (empleado o sistema)
	CreatedByUserID *uint `gorm:"index"`

	StartAt        time.Time `gorm:"not null;index"`
	EndAt          time.Time `gorm:"not null"`
	NumberOfGuests int       `gorm:"not null"`
	StatusID       uint      `gorm:"not null;index"`

	Restaurant Restaurant        `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Table      Table             `gorm:"foreignKey:TableID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Client     Client            `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedBy  User              `gorm:"foreignKey:CreatedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Status     ReservationStatus `gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// ───────────────────────────────────────────
//
//	RESERVATION STATUS
//
// ───────────────────────────────────────────
type ReservationStatus struct {
	gorm.Model
	Code string `gorm:"size:20;unique;not null"` // Ej: "asignado"
	Name string `gorm:"size:50;not null"`        // Ej: "Asignado"
}

// ───────────────────────────────────────────
//
//	RESERVATION STATUS HISTORY
//
// ───────────────────────────────────────────
type ReservationStatusHistory struct {
	gorm.Model
	ReservationID   uint  `gorm:"not null;index"`
	TableID         *uint `gorm:"index"` // Ahora opcional (nullable)
	StatusID        uint  `gorm:"not null;index"`
	ChangedByUserID *uint `gorm:"index"` // Quién hizo el cambio (puede ser null si fue automático)

	Reservation Reservation       `gorm:"foreignKey:ReservationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Status      ReservationStatus `gorm:"foreignKey:StatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ChangedBy   User              `gorm:"foreignKey:ChangedByUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// ───────────────────────────────────────────
//
//	CONSTANTES Y ENUMS
//
// ───────────────────────────────────────────

// Roles predefinidos del sistema
const (
	// Roles de plataforma (super usuarios)
	ROLE_SUPER_ADMIN       = "super_admin"       // Control total del sistema, puede hacer todo
	ROLE_PLATFORM_ADMIN    = "platform_admin"    // Administra la plataforma, gestiona restaurantes y usuarios
	ROLE_PLATFORM_OPERATOR = "platform_operator" // Operaciones básicas, solo lectura y reportes

	// Roles de restaurante
	ROLE_RESTAURANT_OWNER   = "restaurant_owner"   // Dueño del restaurante, control total del negocio
	ROLE_RESTAURANT_MANAGER = "restaurant_manager" // Gerente, gestiona operaciones diarias
	ROLE_RESTAURANT_STAFF   = "restaurant_staff"   // Personal general, acceso básico
	ROLE_WAITER             = "waiter"             // Mesero, ve reservas y mesas
	ROLE_HOST               = "host"               // Anfitrión, gestiona reservas y clientes
)

// Scopes de roles y permisos
const (
	SCOPE_PLATFORM   = "platform"
	SCOPE_RESTAURANT = "restaurant"
	SCOPE_BOTH       = "both"
)

// Recursos del sistema (entidades que se pueden gestionar)
const (
	RESOURCE_RESTAURANTS  = "restaurants"  // Restaurantes de la plataforma
	RESOURCE_USERS        = "users"        // Usuarios del sistema
	RESOURCE_ROLES        = "roles"        // Roles y permisos
	RESOURCE_PERMISSIONS  = "permissions"  // Permisos individuales
	RESOURCE_RESERVATIONS = "reservations" // Reservas de restaurantes
	RESOURCE_TABLES       = "tables"       // Mesas de los restaurantes
	RESOURCE_CLIENTS      = "clients"      // Clientes que hacen reservas
	RESOURCE_DELIVERY     = "delivery"     // Pedidos de delivery (futuro)
	RESOURCE_MENU         = "menu"         // Menús de restaurantes (futuro)
	RESOURCE_REPORTS      = "reports"      // Reportes y estadísticas
)

// Acciones del sistema (operaciones que se pueden realizar)
const (
	ACTION_CREATE = "create" // Crear nuevos registros
	ACTION_READ   = "read"   // Leer/ver información
	ACTION_UPDATE = "update" // Modificar registros existentes
	ACTION_DELETE = "delete" // Eliminar registros
	ACTION_MANAGE = "manage" // Control total (incluye todas las acciones)
)
