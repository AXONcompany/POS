-- name: CreateVenue :one
INSERT INTO venues (
  owner_id, name, address, phone
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetVenueByID :one
SELECT * FROM venues
WHERE id = $1 LIMIT 1;

-- name: ListVenuesByOwner :many
SELECT * FROM venues
WHERE owner_id = $1 AND is_active = true
ORDER BY name;

-- name: UpdateVenue :one
UPDATE venues
SET name = $2,
    address = $3,
    phone = $4,
    is_active = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
