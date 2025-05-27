package service

import (
	"context"

	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService interface {
	GetProductByID(ctx context.Context, id string) (*db.GetProductByIDRow, error)
	ListProducts(ctx context.Context, limit, offset int32) ([]db.ListProductsRow, error)
	GetAllProducts(ctx context.Context) ([]db.GetAllProductsRow, error)
}

type productService struct {
	queries *db.Queries
}

func NewProductService(dbPool *pgxpool.Pool) ProductService {
	return &productService{
		queries: db.New(dbPool),
	}
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*db.GetProductByIDRow, error) {
	product, err := s.queries.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *productService) ListProducts(ctx context.Context, limit, offset int32) ([]db.ListProductsRow, error) {
	return s.queries.ListProducts(ctx, db.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (s *productService) GetAllProducts(ctx context.Context) ([]db.GetAllProductsRow, error) {
	return s.queries.GetAllProducts(ctx)
}
