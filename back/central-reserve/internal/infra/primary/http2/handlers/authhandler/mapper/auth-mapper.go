package mapper

import (
	"central_reserve/internal/domain/dtos"
	"central_reserve/internal/infra/primary/http2/handlers/authhandler/response"
)

// ToLoginResponse convierte el dominio LoginResponse a response.LoginResponse (simplificado)
func ToLoginResponse(domainResponse *dtos.LoginResponse) *response.LoginResponse {
	if domainResponse == nil {
		return nil
	}

	return &response.LoginResponse{
		User: response.UserInfo{
			ID:          domainResponse.User.ID,
			Name:        domainResponse.User.Name,
			Email:       domainResponse.User.Email,
			Phone:       domainResponse.User.Phone,
			AvatarURL:   domainResponse.User.AvatarURL,
			IsActive:    domainResponse.User.IsActive,
			LastLoginAt: domainResponse.User.LastLoginAt,
		},
		Token: domainResponse.Token,
	}
}

// ToUserRolesPermissionsResponse convierte el dominio UserRolesPermissionsResponse a response.UserRolesPermissionsResponse
func ToUserRolesPermissionsResponse(domainResponse *dtos.UserRolesPermissionsResponse) response.UserRolesPermissionsResponse {
	if domainResponse == nil {
		return response.UserRolesPermissionsResponse{}
	}

	return response.UserRolesPermissionsResponse{
		IsSuper:     domainResponse.IsSuper,
		Roles:       toRoleInfoSlice(domainResponse.Roles),
		Permissions: toPermissionInfoSlice(domainResponse.Permissions),
	}
}

// toRoleInfoSlice convierte un slice de dtos.RoleInfo a response.RoleInfo
func toRoleInfoSlice(domainRoles []dtos.RoleInfo) []response.RoleInfo {
	if domainRoles == nil {
		return nil
	}

	roles := make([]response.RoleInfo, len(domainRoles))
	for i, role := range domainRoles {
		roles[i] = response.RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Description: role.Description,
			Level:       role.Level,
			Scope:       role.Scope,
		}
	}
	return roles
}

// toPermissionInfoSlice convierte un slice de dtos.PermissionInfo a response.PermissionInfo
func toPermissionInfoSlice(domainPermissions []dtos.PermissionInfo) []response.PermissionInfo {
	if domainPermissions == nil {
		return nil
	}

	permissions := make([]response.PermissionInfo, len(domainPermissions))
	for i, permission := range domainPermissions {
		permissions[i] = response.PermissionInfo{
			ID:          permission.ID,
			Name:        permission.Name,
			Code:        permission.Code,
			Description: permission.Description,
			Resource:    permission.Resource,
			Action:      permission.Action,
			Scope:       permission.Scope,
		}
	}
	return permissions
}
