package service

import (
	"context"
	"fmt"

	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryService interface {
	GetAllCategory(ctx context.Context) ([]db.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*db.Category, error)
	DeleteByID(ctx context.Context, id string) error
	GetCategoriesPaginated(ctx context.Context, limit, offset int32) (*PaginatedCategories, error)
}

type PaginatedCategories struct {
	Items       []db.Category
	TotalItems  int64
	CurrentPage int32
	TotalPages  int32
	Limit       int32
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
	// Use GetAllCategory with default pagination
	result, err := s.queries.GetAllCategory(ctx, db.GetAllCategoryParams{
		Limit:  100, // Default limit
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %v", err)
	}
	return result, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*db.Category, error) {
	// Convert string ID to int32
	var categoryID int32
	_, err := fmt.Sscanf(id, "%d", &categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID format: %v", err)
	}

	category, err := s.queries.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *categoryService) GetCategoriesPaginated(ctx context.Context, limit, offset int32) (*PaginatedCategories, error) {
	params := db.GetCategoriesPaginatedParams{
		Limit:  limit,
		Offset: offset,
	}

	// Get total count using GetAllCategory with a large limit
	allCategories, err := s.queries.GetAllCategory(ctx, db.GetAllCategoryParams{
		Limit:  1000000, // Very large number to get all
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get total categories: %v", err)
	}
	total := int64(len(allCategories))

	// Get paginated categories
	categories, err := s.queries.GetCategoriesPaginated(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get paginated categories: %v", err)
	}

	// Calculate total pages
	totalPages := (total + int64(limit) - 1) / int64(limit)
	currentPage := (offset / limit) + 1

	return &PaginatedCategories{
		Items:       categories,
		TotalItems:  total,
		CurrentPage: currentPage,
		TotalPages:  int32(totalPages),
		Limit:       limit,
	}, nil
}

func (s *categoryService) DeleteByID(ctx context.Context, id string) error {
	// Convert string ID to int32
	var categoryID int32
	_, err := fmt.Sscanf(id, "%d", &categoryID)
	if err != nil {
		return fmt.Errorf("invalid category ID format: %v", err)
	}

	// Use SQLC's DeleteCategory method with the provided context
	err = s.queries.DeleteCategory(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}
	return nil
}
