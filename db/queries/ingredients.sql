-- name: CreateIngredient :one
insert into ingredients (
  venue_id,
  ingredient_name,
  unit_of_measure,
  ingredient_type
)
values ($1, $2, $3, $4)
returning id, venue_id, created_at, updated_at, deleted_at, ingredient_name, unit_of_measure, ingredient_type, stock;

-- name: GetIngredientByID :one
select id, venue_id, created_at, updated_at, deleted_at, ingredient_name, unit_of_measure, ingredient_type, stock
from ingredients
where id = $1 and venue_id = $2 and deleted_at is null;

-- name: ListIngredients :many
select id, venue_id, created_at, updated_at, deleted_at, ingredient_name, unit_of_measure, ingredient_type, stock
from ingredients
where venue_id = $1 and deleted_at is null
order by id
limit $2 offset $3;

-- name: UpdateIngredient :one
update ingredients
set
  ingredient_name = $3,
  unit_of_measure = $4,
  ingredient_type = $5,
  stock = $6,
  updated_at = now()
where id = $1 and venue_id = $2 and deleted_at is null
returning id, venue_id, created_at, updated_at, deleted_at, ingredient_name, unit_of_measure, ingredient_type, stock;

-- name: DeleteIngredient :exec
update ingredients
set deleted_at = now()
where id = $1 and venue_id = $2 and deleted_at is null;

-- name: ListAllIngredients :many
select id, venue_id, created_at, updated_at, deleted_at, ingredient_name, unit_of_measure, ingredient_type, stock
from ingredients
where venue_id = $1 and deleted_at is null
order by ingredient_name;
