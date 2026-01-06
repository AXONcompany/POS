-- Ejemplo mínimo de sqlc

-- name: CreateCategory :one
insert into categories (category_name)
values ($1)
returning id, created_at, updated_at, deleted_at, category_name;
