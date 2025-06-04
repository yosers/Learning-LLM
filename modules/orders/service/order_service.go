package service

import (
	"context"

	db "shofy/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*db.CreateOrderRow, error)
	GetOrders(ctx context.Context, limit, offset int32, page int) (*PaginatedOrders, error)
	GetOrderById(ctx context.Context, id string) (*db.GetOrderByIdRow, error)
	UpdateOrder(ctx context.Context, req *UpdateOrderRequest) (*db.Order, error)
	DeleteOrder(ctx context.Context, id string) error
}

func NewOrderService(dbPool *pgxpool.Pool) OrderService {
	return &orderService{
		queries: db.New(dbPool),
	}
}

type CreateOrderRequest struct {
	ID string `json:"id"`
}
