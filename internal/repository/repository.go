package repository

import (
	"L0/internal"
	"context"
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	db    *pgxpool.Pool
	redis *redis.Pool
}

func NewRepository(db *pgxpool.Pool, redis *redis.Pool) internal.Repository {
	return &Repository{db: db, redis: redis}
}

const (
	SaveOrderSQL = "INSERT INTO orders(order_uid, track_number, entry, delivery_name, delivery_phone, delivery_zip, " +
		"delivery_city, delivery_address, delivery_region, delivery_email, payment_transaction, payment_request_id, " +
		"payment_currency, payment_provider, payment_amount, payment_dt, payment_bank, payment_delivery_cost, " +
		"payment_goods_total, payment_custom_fee, locale, internal_signature, customer_id, delivery_service, shardkey," +
		" sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, " +
		"$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28);"

	SaveOrderItemsSQL = "INSERT INTO order_item(chrt_id, order_uid) VALUES($1, $2);"

	GetOrderSQL = "SELECT order_uid, track_number, entry, delivery_name, delivery_phone, delivery_zip, " +
		"delivery_city, delivery_address, delivery_region, delivery_email, payment_transaction, payment_request_id, " +
		"payment_currency, payment_provider, payment_amount, payment_dt, payment_bank, payment_delivery_cost, " +
		"payment_goods_total, payment_custom_fee, locale, internal_signature, customer_id, delivery_service, shardkey," +
		" sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1;"

	GetOrdersSQL = "SELECT order_uid, track_number, entry, delivery_name, delivery_phone, delivery_zip, " +
		"delivery_city, delivery_address, delivery_region, delivery_email, payment_transaction, payment_request_id, " +
		"payment_currency, payment_provider, payment_amount, payment_dt, payment_bank, payment_delivery_cost, " +
		"payment_goods_total, payment_custom_fee, locale, internal_signature, customer_id, delivery_service, shardkey," +
		" sm_id, date_created, oof_shard FROM orders;"

	GetItemsSQL = "SELECT items.chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items" +
		" join order_item oi on items.chrt_id = oi.chrt_id WHERE order_uid = $1;"
)

func (r *Repository) getItems(ctx context.Context, orderUID string) ([]internal.Item, error) {
	items := make([]internal.Item, 0)

	rows, err := r.db.Query(ctx, GetItemsSQL, &orderUID)

	for rows.Next() {
		item := internal.Item{}
		err = rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r Repository) GetOrdersPG(ctx context.Context) ([]internal.Order, error) {
	orders := make([]internal.Order, 0)

	rows, err := r.db.Query(ctx, GetOrdersSQL)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		order := internal.Order{}

		err = rows.Scan(&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Delivery.Name,
			&order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
			&order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank,
			&order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee, &order.Locale,
			&order.InternalSignature, &order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated,
			&order.OofShard)
		if err != nil {
			return nil, err
		}
		order.Items, err = r.getItems(ctx, order.OrderUid)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r Repository) GetOrderPG(ctx context.Context, orderUID string) (*internal.Order, error) {
	data := internal.Order{}

	err := r.db.QueryRow(ctx, GetOrderSQL, &orderUID).Scan(&data.OrderUid, &data.TrackNumber, &data.Entry, &data.Delivery.Name,
		&data.Delivery.Phone, &data.Delivery.Zip, &data.Delivery.City, &data.Delivery.Address, &data.Delivery.Region,
		&data.Delivery.Email, &data.Payment.Transaction, &data.Payment.RequestId, &data.Payment.Currency,
		&data.Payment.Provider, &data.Payment.Amount, &data.Payment.PaymentDt, &data.Payment.Bank,
		&data.Payment.DeliveryCost, &data.Payment.GoodsTotal, &data.Payment.CustomFee, &data.Locale,
		&data.InternalSignature, &data.CustomerId, &data.DeliveryService, &data.Shardkey, &data.SmId, &data.DateCreated,
		&data.OofShard)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	data.Items, err = r.getItems(ctx, data.OrderUid)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r Repository) SaveOrderPG(ctx context.Context, data *internal.Order) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, SaveOrderSQL, &data.OrderUid, &data.TrackNumber, &data.Entry, &data.Delivery.Name,
		&data.Delivery.Phone, &data.Delivery.Zip, &data.Delivery.City, &data.Delivery.Address, &data.Delivery.Region,
		&data.Delivery.Email, &data.Payment.Transaction, &data.Payment.RequestId, &data.Payment.Currency,
		&data.Payment.Provider, &data.Payment.Amount, &data.Payment.PaymentDt, &data.Payment.Bank,
		&data.Payment.DeliveryCost, &data.Payment.GoodsTotal, &data.Payment.CustomFee, &data.Locale,
		&data.InternalSignature, &data.CustomerId, &data.DeliveryService, &data.Shardkey, &data.SmId, &data.DateCreated,
		&data.OofShard)
	if err != nil {
		return err
	}

	for _, item := range data.Items {
		_, err = tx.Exec(ctx, SaveOrderItemsSQL, &item.ChrtId, &data.OrderUid)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetOrderRedis(ctx context.Context, uid string) (*internal.Order, error) {
	connRedis := r.redis.Get()
	defer connRedis.Close()

	orderJSON, err := redis.String(connRedis.Do("GET", uid))
	if err != nil {
		return nil, err
	}

	order := internal.Order{}
	err = json.Unmarshal([]byte(orderJSON), &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r Repository) SaveOrderRedis(ctx context.Context, data *internal.Order) error {
	connRedis := r.redis.Get()
	defer connRedis.Close()

	orderJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = connRedis.Do("SET", data.OrderUid, orderJSON)
	if err != nil {
		return err
	}
	return nil
}
