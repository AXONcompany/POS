-- name: CreateCategory :one
insert into categories (
  category_name
) values ($1)
returning id, created_at, updated_at, deleted_at, category_name;

-- name: GetCategory :one
select id, created_at, updated_at, deleted_at, category_name
from categories
where id = $1 and deleted_at is null;

-- name: ListCategories :many
select id, created_at, updated_at, deleted_at, category_name
from categories
where deleted_at is null
order by id
limit $1 offset $2;

-- name: UpdateCategory :one
update categories
set
  category_name = $2,
  updated_at = now()
where id = $1 and deleted_at is null
returning id, created_at, updated_at, deleted_at, category_name;

-- name: DeleteCategory :exec
update categories
set deleted_at = now()
where id = $1 and deleted_at is null;
