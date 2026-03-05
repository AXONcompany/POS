-- name: CreateTable :one
INSERT INTO tables (
  venue_id, table_number, status, capacity, arrival_time
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListTables :many
SELECT * FROM tables
WHERE venue_id = $1
ORDER BY table_number;

-- name: GetTable :one
SELECT * FROM tables
WHERE id_table = $1 AND venue_id = $2 LIMIT 1;

-- name: UpdateTableStatus :one
UPDATE tables
SET status = $3, arrival_time = $4, updated_at = now()
WHERE id_table = $1 AND venue_id = $2
RETURNING *;

-- name: UpdateTable :exec
UPDATE tables
SET 
    table_number = COALESCE(sqlc.narg('table_number'), table_number),
    capacity     = COALESCE(sqlc.narg('capacity'), capacity),
    status       = COALESCE(sqlc.narg('status'), status),
    arrival_time = COALESCE(sqlc.narg('arrival_time'), arrival_time),
    updated_at   = now()
WHERE id_table = sqlc.arg('id_table') AND venue_id = sqlc.arg('venue_id');

-- name: DeleteTable :exec
DELETE FROM tables
WHERE id_table = $1 AND venue_id = $2;