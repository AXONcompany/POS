-- name: CreateTerminal :one
INSERT INTO pos_terminals (
  venue_id, terminal_name
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetTerminalByID :one
SELECT * FROM pos_terminals
WHERE id = $1 LIMIT 1;

-- name: ListTerminalsByVenue :many
SELECT * FROM pos_terminals
WHERE venue_id = $1 AND is_active = true
ORDER BY terminal_name;

-- name: UpdateTerminal :one
UPDATE pos_terminals
SET terminal_name = $2,
    is_active = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
