package service

type UpdateRolesRequest struct {
	Name      string `json:"name"`
	Is_active bool   `json:"is_active"`
	Id        int32  `json:"id"`
}

type ListRoleResponse struct {
	Roles []RoleResponse `json:"roles"`
}

type RoleResponse struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type CreateRolesRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}
