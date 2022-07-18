#!/bin/bash
set -e
export PGPASSWORD=$POSTGRES_PASSWORD;
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE USER $APP_DB_USER WITH PASSWORD '$APP_DB_PASS';
  CREATE DATABASE $APP_DB_NAME;
  GRANT ALL PRIVILEGES ON DATABASE $APP_DB_NAME TO $APP_DB_USER;
  \connect $APP_DB_NAME $APP_DB_USER

  create table if not exists items (
      chrt_id int primary key not null,
      track_number text,
      price int,
      rid text,
      name text,
      sale int,
      size text,
      total_price int,
      nm_id int,
      brand text,
      status int
  );

  create table if not exists orders (
      order_uid text primary key not null ,
      track_number text,
      entry text,
      delivery_name text,
      delivery_phone text,
      delivery_zip text,
      delivery_city text,
      delivery_address text,
      delivery_region text,
      delivery_email text,
      payment_transaction text,
      payment_request_id text,
      payment_currency text,
      payment_provider text,
      payment_amount int,
      payment_dt timestamptz,
      payment_bank text,
      payment_delivery_cost int,
      payment_goods_total int,
      payment_custom_fee int,
      locale text,
      internal_signature text,
      customer_id text,
      delivery_service text,
      shardkey text,
      sm_id int,
      date_created timestamptz,
      oof_shard text
  );

  create table order_item (
      chrt_id int not null ,
      order_uid text not null,
      CONSTRAint fk_order
          FOREIGN KEY(order_uid)
              REFERENCES orders(order_uid),
      CONSTRAint fk_item
          FOREIGN KEY(chrt_id)
              REFERENCES items(chrt_id)
  )
EOSQL
