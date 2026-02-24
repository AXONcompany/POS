-- name: CreateProduct :one
insert into products (
  product_name,
  sales_price,
  is_active
) values ($1, $2, $3)
returning id, created_at, updated_at, deleted_at, product_name, sales_price, is_active;

-- name: GetProduct :one
select id, created_at, updated_at, deleted_at, product_name, sales_price, is_active
from products
where id = $1 and deleted_at is null;

-- name: ListProducts :many
select id, created_at, updated_at, deleted_at, product_name, sales_price, is_active
from products
where deleted_at is null
order by id
limit $1 offset $2;

-- name: UpdateProduct :one
update products
set
  product_name = $2,
  sales_price = $3,
  is_active = $4,
  updated_at = now()
where id = $1 and deleted_at is null
returning id, created_at, updated_at, deleted_at, product_name, sales_price, is_active;

-- name: DeleteProduct :exec
update products
set deleted_at = now()
where id = $1 and deleted_at is null;
