-- name: CreateRestaurant :one
INSERT INTO restaurants (
  name, address, phone
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetRestaurantByID :one
SELECT * FROM restaurants
WHERE id = $1 LIMIT 1;

-- name: UpdateRestaurant :one
UPDATE restaurants
SET name = $2,
    address = $3,
    phone = $4,
    is_active = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;


-- name: CreateUser :one
INSERT INTO users (
  restaurant_id, role_id, name, email, password_hash
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET name = $2,
    email = $3,
    is_active = $4,
    role_id = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: ListUsersByRestaurant :many
SELECT * FROM users
WHERE restaurant_id = $1
ORDER BY name;


-- name: GetRoleByName :one
SELECT * FROM roles
WHERE name = $1 LIMIT 1;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1 LIMIT 1;


-- name: CreateSession :one
INSERT INTO sessions (
  user_id, refresh_token, expires_at, device_info, ip_address
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetSessionByToken :one
SELECT * FROM sessions
WHERE refresh_token = $1 LIMIT 1;

-- name: RevokeSession :exec
UPDATE sessions
SET is_revoked = true
WHERE refresh_token = $1;
