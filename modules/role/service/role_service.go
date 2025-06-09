package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleService interface {
	ListRole(ctx context.Context) (*ListRoleResponse, error)
	RolesByID(ctx context.Context, idRole int32) (*RoleResponse, error)
	DeleteRolesByID(ctx context.Context, idRole int32) error
	UpdateRolesById(ctx context.Context, req *UpdateRolesRequest) error
	CreateRoles(ctx context.Context, req *CreateRolesRequest) (*RoleResponse, error)
}

type ListRoleResponse struct {
	Roles []RoleResponse `json:"roles"`
}

type RoleResponse struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type UpdateRolesRequest struct {
	Name      string `json:"name"`
	Is_active bool   `json:"is_active"`
	Id        int32  `json:"id"`
}

type CreateRolesRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func NewRoleService(dbPool *pgxpool.Pool) RoleService {
	return &roleService{
		queries: db.New(dbPool),
	}
}

type roleService struct {
	queries *db.Queries
}

func (s *roleService) ListRole(ctx context.Context) (*ListRoleResponse, error) {
	rows, err := s.queries.ListRole(ctx)
	if err != nil {
		return nil, err
	}

	var roles []RoleResponse
	for _, row := range rows {
		roles = append(roles, RoleResponse{
			Id:   row.ID,
			Name: row.Name,
		})
	}

	return &ListRoleResponse{
		Roles: roles,
	}, nil

}

func (s *roleService) RolesByID(ctx context.Context, idRole int32) (*RoleResponse, error) {
	rows, err := s.queries.GetRoleByID(ctx, idRole)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &RoleResponse{
		Name: rows.Name,
		Id:   rows.ID,
	}, nil
}

func (s *roleService) DeleteRolesByID(ctx context.Context, idRole int32) error {
	_, err := s.queries.GetRoleByID(ctx, idRole)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("role not found")
		}
		return fmt.Errorf("failed to get role: %w", err)
	}

	err = s.queries.DeleteRoleById(ctx, idRole)

	if err != nil {
		return err
	}
	return nil
}

func (s *roleService) UpdateRolesById(ctx context.Context, req *UpdateRolesRequest) error {
	_, err := s.queries.GetRoleByID(ctx, req.Id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("role not found")
		}
		return fmt.Errorf("failed to get role: %w", err)
	}

	err = s.queries.UpdateRoleById(ctx, db.UpdateRoleByIdParams{
		ID:       req.Id,
		Name:     req.Name,
		IsActive: req.Is_active})

	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

func (s *roleService) CreateRoles(ctx context.Context, req *CreateRolesRequest) (*RoleResponse, error) {
	// Cek apakah role dengan nama tersebut sudah ada
	_, err := s.queries.GetRoleByName(ctx, req.Name)
	if err == nil {
		// Role ditemukan => return error
		return nil, fmt.Errorf("role with name %s already exists", req.Name)
	}

	// Jika error bukan karena "data tidak ditemukan", return error
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}

	// Jika sampai sini, berarti role belum ada => buat baru
	newRole, err := s.queries.CreateRole(ctx, db.CreateRoleParams{
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &RoleResponse{
		Id:       newRole.ID,
		Name:     newRole.Name,
		IsActive: newRole.IsActive,
	}, nil
}
