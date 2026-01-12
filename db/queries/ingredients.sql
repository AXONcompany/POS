-- name: CreateIngredient :one
insert into ingredients (
  ingredient_name,
  unit_of_measure,
  ingredient_type
)
values ($1, $2, $3)
returning id, created_at, updated_at, deleted_at, ingredient_name, unit_of_measure, ingredient_type, stock;

-- name: GetIngredientByID :one
select id, created_at, updated_at, deleted_at, ingredient_name, unit_of_measure, ingredient_type, stock
from ingredients
where id = $1 and deleted_at is null;

-- name: ListIngredients :many
select id, created_at, updated_at, deleted_at, ingredient_name, unit_of_measure, ingredient_type, stock
from ingredients
where deleted_at is null
order by id
limit $1 offset $2;

-- name: UpdateIngredient :one
update ingredients
set
  ingredient_name = $2,
  unit_of_measure = $3,
  ingredient_type = $4,
  updated_at = now()
where id = $1 and deleted_at is null
returning id, created_at, updated_at, deleted_at, ingredient_name, unit_of_measure, ingredient_type, stock;

-- name: DeleteIngredient :exec
update ingredients
set deleted_at = now()
where id = $1 and deleted_at is null;
