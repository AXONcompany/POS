-- name: CreateTable :one
INSERT INTO tables (
  table_number, status, capacity, arrival_time
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: ListTables :many
SELECT * FROM tables
ORDER BY table_number;

-- name: GetTable :one
SELECT * FROM tables
WHERE id = $1 LIMIT 1;

-- name: UpdateTableStatus :one
UPDATE tables
SET status = $2, arrival_time = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateTable :exec
UPDATE tables
SET 
    -- COALESCE significa: "Si el primer valor es nulo, usa el segundo"
    -- sqlc.narg() permite pasar nulos desde Go
    table_number = COALESCE(sqlc.narg('table_number'), table_number),
    capacity     = COALESCE(sqlc.narg('capacity'), capacity),
    status       = COALESCE(sqlc.narg('status'), status),
    arrival_time = COALESCE(sqlc.narg('arrival_time'), arrival_time),
    updated_at   = now()
WHERE id = sqlc.arg('id');

-- name: DeleteTable :exec
DELETE FROM tables
WHERE id = $1;

-- SECCION: Asignación de Mesas (Table Waitress) --

-- name: AssignWaitressToTable :one
INSERT INTO table_waitress (
  table_id, waitress_id
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetWaitressByTable :one
SELECT * FROM table_waitress
WHERE table_id = $1 LIMIT 1;

-- name: RemoveWaitressFromTable :exec
DELETE FROM table_waitress
WHERE table_id = $1;