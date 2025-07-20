package mapper

import (
	"central_reserve/internal/domain/dtos"
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/request"
	"central_reserve/internal/infra/primary/http2/handlers/businesstypehandler/response"
)

// ToBusinessTypeRequest convierte el request HTTP a DTO de dominio
func ToBusinessTypeRequest(req request.BusinessTypeRequest) dtos.BusinessTypeRequest {
	return dtos.BusinessTypeRequest{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Icon:        req.Icon,
		IsActive:    req.IsActive,
	}
}

// ToBusinessTypeResponse convierte el DTO de dominio a response HTTP
func ToBusinessTypeResponse(dto dtos.BusinessTypeResponse) response.BusinessTypeResponse {
	return response.BusinessTypeResponse{
		ID:          dto.ID,
		Name:        dto.Name,
		Code:        dto.Code,
		Description: dto.Description,
		Icon:        dto.Icon,
		IsActive:    dto.IsActive,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

// ToBusinessTypeListResponse convierte la lista de DTOs de dominio a response HTTP
func ToBusinessTypeListResponse(dtos []dtos.BusinessTypeResponse) response.BusinessTypeListResponse {
	businessTypes := make([]response.BusinessTypeResponse, len(dtos))
	for i, dto := range dtos {
		businessTypes[i] = ToBusinessTypeResponse(dto)
	}

	return response.BusinessTypeListResponse{
		BusinessTypes: businessTypes,
		Total:         int64(len(businessTypes)),
		Page:          1,
		Limit:         len(businessTypes),
	}
}
