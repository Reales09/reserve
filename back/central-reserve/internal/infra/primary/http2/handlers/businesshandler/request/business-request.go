package request

// BusinessRequest representa la solicitud para crear/actualizar un negocio
type BusinessRequest struct {
	Name           string `json:"name" binding:"required"`
	Code           string `json:"code" binding:"required"`
	BusinessTypeID uint   `json:"business_type_id" binding:"required"`
	Timezone       string `json:"timezone"`
	Address        string `json:"address"`
	Description    string `json:"description"`

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
}
