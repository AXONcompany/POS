-- name: CreateCategory :one
insert into categories (
  venue_id, category_name, is_active
) values ($1, $2, true)
returning id, venue_id, created_at, updated_at, deleted_at, category_name;

-- name: GetCategory :one
select id, venue_id, created_at, updated_at, deleted_at, category_name
from categories
where id = $1 and venue_id = $2 and deleted_at is null;

-- name: ListCategories :many
select id, venue_id, created_at, updated_at, deleted_at, category_name
from categories
where venue_id = $1 and deleted_at is null
order by id
limit $2 offset $3;

-- name: UpdateCategory :one
update categories
set
  category_name = $3,
  updated_at = now()
where id = $1 and venue_id = $2 and deleted_at is null
returning id, venue_id, created_at, updated_at, deleted_at, category_name;

-- name: DeleteCategory :exec
update categories
set deleted_at = now()
where id = $1 and venue_id = $2 and deleted_at is null;

