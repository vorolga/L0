package internal

import "context"

type Repository interface {
	GetOrdersPG(ctx context.Context) ([]Order, error)
	GetOrderPG(ctx context.Context, uid string) (*Order, error)
	SaveOrderPG(ctx context.Context, data *Order) error

	GetOrderRedis(ctx context.Context, uid string) (*Order, error)
	SaveOrderRedis(ctx context.Context, data *Order) error
}
