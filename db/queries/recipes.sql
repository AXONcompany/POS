-- name: AddRecipeItem :one
insert into recipe (
  product_id,
  ingredient_id,
  quantity_required
) values ($1, $2, $3)
returning id, created_at, updated_at, deleted_at, product_id, ingredient_id, quantity_required;

-- name: GetRecipeByProductID :many
select r.id, r.product_id, r.ingredient_id, r.quantity_required, i.ingredient_name, i.unit_of_measure
from recipe r
join ingredients i on r.ingredient_id = i.id
where r.product_id = $1 and r.deleted_at is null;

-- name: DeleteRecipeItem :exec
update recipe
set deleted_at = now()
where id = $1 and deleted_at is null;
