package usecase

import (
	"L0/internal"
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4"
)

type Usecase struct {
	repository internal.Repository
}

func NewUsecase(r internal.Repository) (internal.Usecase, error) {
	u := Usecase{repository: r}
	err := u.InitCache(context.Background())
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (u Usecase) SaveOrder(ctx context.Context, data *internal.Order) error {
	err := u.repository.SaveOrderPG(ctx, data)
	if err != nil {
		return err
	}

	err = u.repository.SaveOrderRedis(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) InitCache(ctx context.Context) error {
	orders, err := u.repository.GetOrdersPG(ctx)
	if err != nil {
		return err
	}

	for _, order := range orders {
		err = u.repository.SaveOrderRedis(ctx, &order)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u Usecase) GetOrder(ctx context.Context, uid string) (*internal.Order, error) {
	order, err := u.repository.GetOrderRedis(ctx, uid)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}

	if order != nil {
		return order, nil
	}

	order, err = u.repository.GetOrderPG(ctx, uid)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	return order, nil
}
