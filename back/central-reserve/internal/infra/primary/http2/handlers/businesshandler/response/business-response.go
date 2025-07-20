package response

import "time"

// BusinessTypeResponse representa la respuesta de un tipo de negocio
type BusinessTypeResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BusinessResponse representa la respuesta de un negocio
type BusinessResponse struct {
	ID           uint                 `json:"id"`
	Name         string               `json:"name"`
	Code         string               `json:"code"`
	BusinessType BusinessTypeResponse `json:"business_type"`
	Timezone     string               `json:"timezone"`
	Address      string               `json:"address"`
	Description  string               `json:"description"`

	// Configuración de marca blanca
	LogoURL        string `json:"logo_url"`
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	CustomDomain   string `json:"custom_domain"`
	IsActive       bool   `json:"is_active"`

	// Configuración de funcionalidades
	EnableDelivery     bool `json:"enable_delivery"`
	EnablePickup       bool `json:"enable_pickup"`
	EnableReservations bool `json:"enable_reservations"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BusinessListResponse representa la respuesta de lista de negocios
type BusinessListResponse struct {
	Businesses []BusinessResponse `json:"businesses"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
}

// BusinessSuccessResponse representa la respuesta exitosa de un negocio
type BusinessSuccessResponse struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    BusinessResponse `json:"data"`
}

// BusinessListSuccessResponse representa la respuesta exitosa de lista de negocios
type BusinessListSuccessResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    BusinessListResponse `json:"data"`
}

// BusinessErrorResponse representa la respuesta de error
type BusinessErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
