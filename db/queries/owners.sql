-- name: CreateOwner :one
INSERT INTO owners (
  name, email, password_hash
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetOwnerByID :one
SELECT * FROM owners
WHERE id = $1 LIMIT 1;

-- name: GetOwnerByEmail :one
SELECT * FROM owners
WHERE email = $1 LIMIT 1;

-- name: UpdateOwner :one
UPDATE owners
SET name = $2,
    email = $3,
    is_active = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListOwners :many
SELECT * FROM owners
WHERE is_active = true
ORDER BY name;
