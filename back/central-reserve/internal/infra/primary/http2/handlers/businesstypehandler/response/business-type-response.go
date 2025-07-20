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

// BusinessTypeListResponse representa la respuesta de lista de tipos de negocio
type BusinessTypeListResponse struct {
	BusinessTypes []BusinessTypeResponse `json:"business_types"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
}

// BusinessTypeSuccessResponse representa la respuesta exitosa de un tipo de negocio
type BusinessTypeSuccessResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    BusinessTypeResponse `json:"data"`
}

// BusinessTypeListSuccessResponse representa la respuesta exitosa de lista de tipos de negocio
type BusinessTypeListSuccessResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    BusinessTypeListResponse `json:"data"`
}

// BusinessTypeErrorResponse representa la respuesta de error
type BusinessTypeErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
