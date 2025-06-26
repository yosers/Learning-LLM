package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgtype"

	"fmt"
	db "shofy/db/sqlc"
	model "shofy/modules/shops/model"
	"shofy/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShopService interface {
	GetShopsByID(ctx context.Context, id int32) (model.ShopsResponse, error)
	ListShops(ctx context.Context, req *model.ListShopRequest) (*model.ListShopsResponse, error)
	// GetAllShops(ctx context.Context) ([]db.GetAllShopsRow, error)
	DeleteShopsByID(ctx context.Context, id int32) error
	CreateShops(ctx context.Context, req *model.ShopsRequest) (*model.ShopsResponse, error)
	UpdateShops(ctx context.Context, userId int32, req *model.ShopsRequest) (*model.ShopsResponse, error)
}

func NewShopsService(dbPool *pgxpool.Pool) ShopService {
	return &shopService{
		queries: db.New(dbPool),
	}
}

type shopService struct {
	queries *db.Queries
}

func (s *shopService) ListShops(ctx context.Context, req *model.ListShopRequest) (*model.ListShopsResponse, error) {
	// Get total count first
	total, err := s.queries.CountUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Calculate pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize
	totalPages := int32((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	// Get paginated users
	shops, err := s.queries.ListShops(ctx, db.ListShopsParams{
		Limit:  req.PageSize,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Get user profiles for all users
	shopsResponses := make([]model.ShopsResponse, 0, len(shops))
	for _, shop := range shops {
		shopsResponses = append(shopsResponses, model.ShopsResponse{
			ID:            shop.ID,
			Name:          shop.Name,
			Description:   shop.Description,
			LogoUrl:       shop.LogoUrl.String,
			WebsiteUrl:    shop.WebsiteUrl.String,
			Email:         shop.Email.String,
			WhatsappPhone: shop.WhatsappPhone.String,
			Address:       shop.Address,
			City:          shop.City,
			State:         shop.State,
			IsActive:      shop.IsActive,
			Latitude:      float32(shop.Latitude),
			Longitude:     float32(shop.Longitude),
			ZipCode:       shop.ZipCode,
			Country:       shop.Country,
		})
	}

	return &model.ListShopsResponse{
		Shops:      shopsResponses,
		Total:      int32(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *shopService) GetShopsByID(ctx context.Context, id int32) (model.ShopsResponse, error) {
	shop, err := s.queries.GetShopsById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ShopsResponse{}, fmt.Errorf("shops not found")
		}
		return model.ShopsResponse{}, fmt.Errorf("failed to get shops: %w", err)
	}

	return model.ShopsResponse{
		ID:            shop.ID,
		Name:          shop.Name,
		Description:   shop.Description,
		LogoUrl:       shop.LogoUrl.String,
		WebsiteUrl:    shop.WebsiteUrl.String,
		Email:         shop.Email.String,
		WhatsappPhone: shop.WhatsappPhone.String,
		Address:       shop.Address,
		City:          shop.City,
		State:         shop.State,
		IsActive:      shop.IsActive,
		Latitude:      float32(shop.Latitude),
		Longitude:     float32(shop.Longitude),
		ZipCode:       shop.ZipCode,
		Country:       shop.Country,
	}, nil
}

func (s *shopService) UpdateShops(ctx context.Context, shopsId int32, req *model.ShopsRequest) (*model.ShopsResponse, error) {
	// Get existing user
	_, err := s.queries.GetShopsById(ctx, shopsId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("shops not found")
		}
		return nil, fmt.Errorf("failed to get shops: %w", err)
	}

	// Update user profile
	result, err := s.queries.UpdateShops(ctx, db.UpdateShopsParams{
		ID:          shopsId,
		Name:        utils.StringOrEmpty(&req.Name),
		Description: utils.StringOrEmpty(&req.Description),
		LogoUrl: pgtype.Text{
			String: req.LogoUrl,
			Valid:  req.LogoUrl != "",
		},
		WebsiteUrl: pgtype.Text{
			String: req.WebsiteUrl,
			Valid:  req.WebsiteUrl != "",
		},
		Email: pgtype.Text{
			String: req.Email,
			Valid:  req.Email != "",
		},
		WhatsappPhone: pgtype.Text{
			String: req.WhatsappPhone,
			Valid:  req.WhatsappPhone != "",
		},
		Address:   utils.StringOrEmpty(&req.Address),
		City:      utils.StringOrEmpty(&req.City),
		State:     utils.StringOrEmpty(&req.State),
		ZipCode:   utils.StringOrEmpty(&req.ZipCode),
		Country:   utils.StringOrEmpty(&req.Country),
		Latitude:  float64(req.Latitude),
		Longitude: float64(req.Longitude),
		IsActive:  req.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update Shops profile: %w", err)
	}

	return &model.ShopsResponse{
		ID:            result.ID,
		Name:          result.Name,
		Description:   result.Description,
		LogoUrl:       result.LogoUrl.String,
		WebsiteUrl:    result.WebsiteUrl.String,
		Email:         result.Email.String,
		WhatsappPhone: result.WhatsappPhone.String,
		Address:       result.Address,
		City:          result.City,
		State:         result.State,
		IsActive:      result.IsActive,
		Latitude:      float32(result.Latitude),
		Longitude:     float32(result.Longitude),
		ZipCode:       result.ZipCode,
		Country:       result.Country,
	}, nil
}

func (s *shopService) DeleteShopsByID(ctx context.Context, id int32) error {
	// Check if user exists
	_, err := s.queries.GetShopsById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("Shops not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Proceed to delete user
	err = s.queries.DeleteShopsById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete Shops: %w", err)
	}

	return nil
}

func (s *shopService) CreateShops(ctx context.Context, req *model.ShopsRequest) (*model.ShopsResponse, error) {

	result, err := s.queries.GetShopsByNameOrWhatshapp(ctx, db.GetShopsByNameOrWhatshappParams{
		Name: req.Name,
		WhatsappPhone: pgtype.Text{
			String: req.WhatsappPhone,
			Valid:  req.WhatsappPhone != "",
		},
	})

	if result.ID != 0 {
		return nil, fmt.Errorf("Shops already exist")
	}

	// Create user
	shopResult, err := s.queries.CreateShops(ctx, db.CreateShopsParams{
		Name:        req.Name,
		Description: req.Description,
		LogoUrl: pgtype.Text{
			String: req.LogoUrl,
			Valid:  req.LogoUrl != "",
		},
		WebsiteUrl: pgtype.Text{
			String: req.WebsiteUrl,
			Valid:  req.WebsiteUrl != "",
		},
		Email: pgtype.Text{
			String: req.Email,
			Valid:  req.Email != "",
		},
		WhatsappPhone: pgtype.Text{
			String: req.WhatsappPhone,
			Valid:  req.WhatsappPhone != "",
		},
		Address:   req.Address,
		City:      req.City,
		State:     req.State,
		ZipCode:   req.ZipCode,
		Country:   req.Country,
		Latitude:  float64(req.Latitude),
		Longitude: float64(req.Longitude),
		IsActive:  true,
	})

	// Check for errors
	if err != nil {
		log.Println("Error creating user:", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &model.ShopsResponse{
		Name:          shopResult.Name,
		Description:   shopResult.Description,
		LogoUrl:       shopResult.LogoUrl.String,
		WebsiteUrl:    shopResult.WebsiteUrl.String,
		Email:         shopResult.Email.String,
		WhatsappPhone: shopResult.WhatsappPhone.String,
		Address:       shopResult.Address,
		City:          shopResult.City,
		State:         shopResult.State,
		IsActive:      result.IsActive,
		Latitude:      float32(shopResult.Latitude),
		Longitude:     float32(shopResult.Longitude),
		ZipCode:       shopResult.ZipCode,
		Country:       shopResult.Country,
	}, nil
}
