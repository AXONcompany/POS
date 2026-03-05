-- name: CreateProduct :one
insert into products (
  venue_id,
  product_name,
  sales_price,
  is_active
) values ($1, $2, $3, $4)
returning id, venue_id, created_at, updated_at, deleted_at, product_name, sales_price, is_active;

-- name: GetProduct :one
select id, venue_id, created_at, updated_at, deleted_at, product_name, sales_price, is_active
from products
where id = $1 and venue_id = $2 and deleted_at is null;

-- name: ListProducts :many
select id, venue_id, created_at, updated_at, deleted_at, product_name, sales_price, is_active
from products
where venue_id = $1 and deleted_at is null
order by id
limit $2 offset $3;

-- name: UpdateProduct :one
update products
set
  product_name = $3,
  sales_price = $4,
  is_active = $5,
  updated_at = now()
where id = $1 and venue_id = $2 and deleted_at is null
returning id, venue_id, created_at, updated_at, deleted_at, product_name, sales_price, is_active;

-- name: DeleteProduct :exec
update products
set deleted_at = now()
where id = $1 and venue_id = $2 and deleted_at is null;
