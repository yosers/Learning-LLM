package service

import (
	"context"

	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryService interface {
	GetAllCategory(ctx context.Context) ([]db.Category, error)
}

type categoryService struct {
	queries *db.Queries
}

func NewCategoryService(dbPool *pgxpool.Pool) CategoryService {
	return &categoryService{
		queries: db.New(dbPool),
	}
}

func (s *categoryService) GetAllCategory(ctx context.Context) ([]db.Category, error) {
	return s.queries.GetAllCategory(ctx)
}
