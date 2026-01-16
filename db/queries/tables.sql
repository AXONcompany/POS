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