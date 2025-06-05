package service

import (
	"context"
	"fmt"
	"math/big"

	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService interface {
	GetProductByID(ctx context.Context, id string) (ListProductsRowSnake, error)
	ListProducts(ctx context.Context, limit, offset int32, page int) (*PaginatedProducts, error)
	GetAllProducts(ctx context.Context) ([]db.GetAllProductsRow, error)
	DeleteProductByID(ctx context.Context, id string) error
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*db.CreateProductRow, error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*db.Product, error)
}

func NewProductService(dbPool *pgxpool.Pool) ProductService {
	return &productService{
		queries: db.New(dbPool),
	}
}

type CreateProductRequest struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	CategoryID  int32   `json:"category_id"`
	ShopID      int32   `json:"shop_id"`
}

type ListProductsRowSnake struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description pgtype.Text      `json:"description"`
	Price       pgtype.Numeric   `json:"price"`
	Stock       pgtype.Int4      `json:"stock"`
	CategoryID  string           `json:"category_id"`
	ShopID      string           `json:"shop_id"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	DeletedAt   pgtype.Timestamp `json:"deleted_at"`
}

func (s *productService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*db.CreateProductRow, error) {
	// Convert price to big.Int (cents)
	priceInCents := big.NewInt(int64(req.Price * 100))

	params := db.CreateProductParams{
		ID:   req.ID,
		Name: req.Name,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
		Price: pgtype.Numeric{
			InfinityModifier: pgtype.Finite,
			Valid:            true,
			Int:              priceInCents,
			Exp:              -2, // For 2 decimal places
		},
		Stock: pgtype.Int4{
			Int32: req.Stock,
			Valid: true,
		},
		CategoryID: pgtype.Int4{
			Int32: req.CategoryID,
			Valid: true,
		},
		ShopID: req.ShopID,
	}

	result, err := s.queries.CreateProduct(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return &result, nil
}

type UpdateProductRequest struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int32   `json:"stock"`
	CategoryID  int32   `json:"category_id"`
	ShopID      int32   `json:"shop_id"`
}

type productService struct {
	queries *db.Queries
}

type PaginatedProducts struct {
	Items       []ListProductsRowSnake
	TotalItems  int64
	CurrentPage int
	TotalPages  int
	Limit       int
}

func mapToSnakeCase(rows []db.ListProductsRow) []ListProductsRowSnake {
	snakes := make([]ListProductsRowSnake, len(rows))
	for i, r := range rows {
		snakes[i] = ListProductsRowSnake{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
			Stock:       r.Stock,
			CategoryID:  r.CategoryID,
			ShopID:      r.ShopID,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.UpdatedAt,
			DeletedAt:   r.DeletedAt,
		}
	}
	return snakes
}

func mapRowToSnakeCase(r db.GetProductByIDRow) ListProductsRowSnake {
	return ListProductsRowSnake{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Stock:       r.Stock,
		CategoryID:  r.CategoryID,
		ShopID:      r.ShopID,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}

func (s *productService) GetProductByID(ctx context.Context, id string) (ListProductsRowSnake, error) {
	product, err := s.queries.GetProductByID(ctx, id)
	if err != nil {
		return ListProductsRowSnake{}, err
	}

	getproduct := mapRowToSnakeCase(product)

	return getproduct, nil
}

func (s *productService) ListProducts(ctx context.Context, limit, offset int32, page int) (*PaginatedProducts, error) {
	// Fetch raw product rows from DB
	itemsRaw, err := s.queries.ListProducts(ctx, db.ListProductsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	// Map to snake_case struct
	items := mapToSnakeCase(itemsRaw)

	// Get total count for pagination
	total, err := s.queries.GetCountProduct(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit)) // ceil(total / limit)

	// Return paginated result
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

func (s *productService) DeleteProductByID(ctx context.Context, id string) error {
	// Use SQLC's DeleteProduct method with the provided context
	err := s.queries.DeleteProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}
	return nil
}

func (s *productService) UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*db.Product, error) {
	priceInCents := big.NewInt(int64(req.Price * 100))

	params := db.UpdateProductParams{
		ID:   req.ID,
		Name: req.Name,
		Description: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
		Price: pgtype.Numeric{
			InfinityModifier: pgtype.Finite,
			Valid:            true,
			Int:              priceInCents,
			Exp:              -2,
		},
		Stock: pgtype.Int4{
			Int32: req.Stock,
			Valid: true,
		},
		CategoryID: pgtype.Int4{
			Int32: req.CategoryID,
			Valid: true,
		},
		ShopID: req.ShopID,
	}

	product, err := s.queries.UpdateProduct(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return &product, nil
}
