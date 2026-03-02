-- name: CreateSale :one
insert into sales (
  total,
  payment_method,
  date,
  order_id
)
values ($1, $2, $3, $4)
returning id, total, payment_method, date, order_id, created_at, updated_at, deleted_at;

-- name: GetSaleByID :one
select id, total, payment_method, date, order_id, created_at, updated_at, deleted_at
from sales
where id = $1 and deleted_at is null;