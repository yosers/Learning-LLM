package service

import (
	"context"

	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService interface {
	GetProductByID(ctx context.Context, id string) (*db.GetProductByIDRow, error)
	ListProducts(ctx context.Context, limit, offset int32, page int) (*PaginatedProducts, error)
	GetAllProducts(ctx context.Context) ([]db.GetAllProductsRow, error)
}

type productService struct {
	queries *db.Queries
}

type PaginatedProducts struct {
	Items       []db.ListProductsRow
	TotalItems  int64
	CurrentPage int
	TotalPages  int
	Limit       int
}

type PaginatedProductsResponse struct {
	Data        []db.ListProductsRow `json:"data"`
	CurrentPage int                  `json:"current_page"`
	TotalPages  int                  `json:"total_pages"`
	TotalItems  int64                `json:"total_items"`
	Limit       int                  `json:"limit"`
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

func (s *productService) ListProducts(ctx context.Context, limit, offset int32, page int) (*PaginatedProducts, error) {
	items, err := s.queries.ListProducts(ctx, db.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, err
	}

	total, err := s.queries.CountProducts(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit)) // ceil(total / limit)

	return &PaginatedProducts{
		Items:       items,
		TotalItems:  total,
		CurrentPage: page,
		TotalPages:  totalPages,
		Limit:       int(limit),
	}, nil

}

func (s *productService) GetAllProducts(ctx context.Context) ([]db.GetAllProductsRow, error) {
	return s.queries.GetAllProducts(ctx)
}
