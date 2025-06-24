package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "shofy/db/sqlc"
	role_model "shofy/modules/role/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleService interface {
	ListRole(ctx context.Context) (*role_model.ListRoleResponse, error)
	RolesByID(ctx context.Context, idRole int32) (*role_model.RoleResponse, error)
	DeleteRolesByID(ctx context.Context, idRole int32) error
	UpdateRolesById(ctx context.Context, req *role_model.UpdateRolesRequest) error
	CreateRoles(ctx context.Context, req *role_model.CreateRolesRequest) (*role_model.RoleResponse, error)
}

func NewRoleService(dbPool *pgxpool.Pool) RoleService {
	return &roleService{
		queries: db.New(dbPool),
	}
}

type roleService struct {
	queries *db.Queries
}

func (s *roleService) ListRole(ctx context.Context) (*role_model.ListRoleResponse, error) {
	rows, err := s.queries.ListRole(ctx)
	if err != nil {
		return nil, err
	}

	var roles []role_model.RoleResponse
	for _, row := range rows {
		roles = append(roles, role_model.RoleResponse{
			Id:   row.ID,
			Name: row.Name,
		})
	}

	return &role_model.ListRoleResponse{
		Roles: roles,
	}, nil

}

func (s *roleService) RolesByID(ctx context.Context, idRole int32) (*role_model.RoleResponse, error) {
	rows, err := s.queries.GetRoleByID(ctx, idRole)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Role not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &role_model.RoleResponse{
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

func (s *roleService) UpdateRolesById(ctx context.Context, req *role_model.UpdateRolesRequest) error {
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

func (s *roleService) CreateRoles(ctx context.Context, req *role_model.CreateRolesRequest) (*role_model.RoleResponse, error) {
	// Cek apakah role dengan nama tersebut sudah ada
	_, err := s.queries.GetRoleByName(ctx, req.Name)
	if err == nil {
		// Role ditemukan => return error
		return nil, fmt.Errorf("role with name %s already exists", req.Name)
	} else if !errors.Is(err, sql.ErrNoRows) {
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

	return &role_model.RoleResponse{
		Id:       newRole.ID,
		Name:     newRole.Name,
		IsActive: newRole.IsActive,
	}, nil
}
