-- name: CreateWaitress :one
INSERT INTO waitress (
  id_user
) VALUES (
  $1
) RETURNING *;

-- name: GetWaitress :one
SELECT * FROM waitress
WHERE id_user = $1 LIMIT 1;

-- name: ListWaitresses :many
SELECT * FROM waitress
ORDER BY id_user;

-- name: DeleteWaitress :exec
DELETE FROM waitress
WHERE id_user = $1;