package internal

import "context"

type Usecase interface {
	SaveOrder(ctx context.Context, data *Order) error
	InitCache(ctx context.Context) error
	GetOrder(ctx context.Context, uid string) (*Order, error)
}
