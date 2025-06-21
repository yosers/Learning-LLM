package service

import (
	"context"
	"fmt"
	"math/big"

	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*db.Order, error)
	GetOrdersList(ctx context.Context, limit, offset int32, page int) (*PaginatedOrders, error)
	GetOrderById(ctx context.Context, id int32) (*db.Order, error)
	UpdateOrder(ctx context.Context, req *UpdateOrderRequest) (*db.Order, error)
	DeleteOrder(ctx context.Context, id int32) error
}

type orderService struct {
	queries *db.Queries
}

func NewOrderService(dbPool *pgxpool.Pool) OrderService {
	return &orderService{
		queries: db.New(dbPool),
	}
}

type CreateOrderRequest struct {
	ShopID int32   `json:"shop_id"`
	UserID int32   `json:"user_id"`
	Total  float64 `json:"total"`
	Status string  `json:"status"`
}

type PaginatedOrders struct {
	Items       []ListOrdersRowSnake
	TotalItems  int64
	CurrentPage int
	TotalPages  int
	Limit       int
}

type ListOrdersRowSnake struct {
	ID        int32            `json:"id"`
	ShopName  string           `json:"shop_name"`
	UserID    int32            `json:"user_id"`
	Total     pgtype.Numeric   `json:"total"`
	Status    pgtype.Text      `json:"status"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func mapToSnakeCase(rows []db.GetListOrdersRow) []ListOrdersRowSnake {
	snakes := make([]ListOrdersRowSnake, len(rows))
	for i, r := range rows {
		snakes[i] = ListOrdersRowSnake{
			ID:        r.ID,
			ShopName:  r.ShopName,
			UserID:    r.UserID,
			Total:     r.Total,
			Status:    r.Status,
			CreatedAt: r.CreatedAt,
		}
	}
	return snakes
}

type UpdateOrderRequest struct {
	ID     int32   `json:"id"`
	Total  float64 `json:"total"`
	Status string  `json:"status"`
}

func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*db.Order, error) {
	totalInCents := big.NewInt(int64(req.Total * 100))

	params := db.CreateOrderParams{
		ShopID: req.ShopID,
		UserID: pgtype.Int4{
			Int32: req.UserID,
			Valid: true,
		},
		Total: pgtype.Numeric{
			InfinityModifier: pgtype.Finite,
			Valid:            true,
			Int:              totalInCents,
			Exp:              -2, // For 2 decimal places
		},
		Status: pgtype.Text{
			String: req.Status,
			Valid:  true,
		},
	}

	// TAMBAHIN INSERT KE TABLE ORDER_ITEMS

	result, err := s.queries.CreateOrder(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return &result, nil
}

func (s *orderService) GetOrdersList(ctx context.Context, limit, offset int32, page int) (*PaginatedOrders, error) {
	itemsRaw, err := s.queries.GetListOrders(ctx, db.GetListOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	items := mapToSnakeCase(itemsRaw)

	total, err := s.queries.GetCountOrder(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &PaginatedOrders{
		Items:       items,
		TotalItems:  total,
		CurrentPage: page,
		TotalPages:  totalPages,
		Limit:       int(limit),
	}, nil
}

func (s *orderService) GetOrderById(ctx context.Context, id int32) (*db.Order, error) {
	order, err := s.queries.GetOrderById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id: %w", err)
	}
	return &order, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, req *UpdateOrderRequest) (*db.Order, error) {
	totalInCents := big.NewInt(int64(req.Total * 100))

	params := db.UpdateOrderParams{
		ID: req.ID,
		Total: pgtype.Numeric{
			InfinityModifier: pgtype.Finite,
			Valid:            true,
			Int:              totalInCents,
			Exp:              -2,
		},
		Status: pgtype.Text{
			String: req.Status,
			Valid:  true,
		},
	}

	order, err := s.queries.UpdateOrder(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return &order, nil
}

func (s *orderService) DeleteOrder(ctx context.Context, id int32) error {
	err := s.queries.DeleteOrder(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}
