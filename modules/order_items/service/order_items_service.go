package service

import (
	"context"
	"fmt"

	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItemsService interface {
	CreateOrderItems(ctx context.Context, req *CreateOrderItemsRequest) (*db.OrderItem, error)
	GetOrderItemsByID(ctx context.Context, id string) (*db.GetOrderItemsByIDRow, error)
}

func NewOrderItemsService(dbPool *pgxpool.Pool) OrderItemsService {
	return &orderItemsService{
		queries: db.New(dbPool),
	}
}

type orderItemsService struct {
	queries *db.Queries
}

type CreateOrderItemsRequest struct {
	ID        int32
	OrderID   int32
	ProductID string
	Quantity  int32
	UnitPrice float64
}

func (s *orderItemsService) CreateOrderItems(ctx context.Context, req *CreateOrderItemsRequest) (*db.OrderItem, error) {
	params := db.CreateOrderItemsParams{
		ID:        req.ID,
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
	}

	result, err := s.queries.CreateOrderItems(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}
	return &result, nil
}

func (s *orderItemsService) GetOrderItemsByID(ctx context.Context, id string) (*db.GetOrderItemsByIDRow, error) {
	params := db.GetOrderItemsByIDParams{
		ID: id,
	}

	result, err := s.queries.GetOrderItemsByID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items by id: %w", err)
	}
	return &result, nil
}
