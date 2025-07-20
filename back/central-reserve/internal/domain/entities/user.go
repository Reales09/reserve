package entities

import "time"

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

// BusinessStaff representa la relación entre usuarios y negocios
type BusinessStaff struct {
	ID         uint
	UserID     uint
	BusinessID uint
	RoleID     uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
